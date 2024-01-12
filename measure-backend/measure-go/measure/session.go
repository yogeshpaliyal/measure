package measure

import (
	"context"
	"errors"
	"fmt"
	"io"
	"measure-backend/measure-go/inet"
	"measure-backend/measure-go/server"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ipinfo/go/v2/ipinfo"
	"github.com/jackc/pgx/v5"
	"github.com/leporo/sqlf"
)

type Session struct {
	SessionID   uuid.UUID    `json:"session_id" binding:"required"`
	AppID       uuid.UUID    `json:"app_id"`
	Timestamp   time.Time    `json:"timestamp" binding:"required"`
	IPv4        net.IP       `json:"inet_ipv4"`
	IPv6        net.IP       `json:"inet_ipv6"`
	CountryCode string       `json:"inet_country_code"`
	Resource    Resource     `json:"resource" binding:"required"`
	Events      []EventField `json:"events" binding:"required"`
	Attachments []Attachment `json:"attachments"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

func (s *Session) validate() error {
	if err := s.Resource.validate(); err != nil {
		return err
	}

	for _, event := range s.Events {
		if err := event.validate(); err != nil {
			return err
		}
	}

	if s.hasAttachments() {
		for _, attachment := range s.Attachments {
			if err := attachment.validate(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *Session) hasAttachments() bool {
	return len(s.Attachments) > 0
}

func (s *Session) hasUnhandledExceptions() bool {
	result := false

	for i := range s.Events {
		if s.Events[i].isUnhandledException() {
			result = true
			break
		}
	}

	return result
}

func (s *Session) hasANRs() bool {
	result := false

	for i := range s.Events {
		if s.Events[i].isANR() {
			result = true
			break
		}
	}

	return result
}

func (s *Session) needsSymbolication() bool {
	result := false
	for i := range s.Events {
		if s.Events[i].isException() {
			result = true
			break
		}

		if s.Events[i].isANR() {
			result = true
			break
		}

		if s.Events[i].isAppExit() && len(s.Events[i].AppExit.Trace) > 0 {
			result = true
			break
		}

		if s.Events[i].isLifecycleActivity() && len(s.Events[i].LifecycleActivity.ClassName) > 0 {
			result = true
			break
		}

		if s.Events[i].isColdLaunch() && len(s.Events[i].ColdLaunch.LaunchedActivity) > 0 {
			result = true
			break
		}

		if s.Events[i].isWarmLaunch() && len(s.Events[i].WarmLaunch.LaunchedActivity) > 0 {
			result = true
			break
		}

		if s.Events[i].isHotLaunch() && len(s.Events[i].HotLaunch.LaunchedActivity) > 0 {
			result = true
			break
		}

		if s.Events[i].isLifecycleFragment() {
			hasClassName := len(s.Events[i].LifecycleFragment.ClassName) > 0
			hasParentActivity := len(s.Events[i].LifecycleFragment.ParentActivity) > 0

			if hasClassName || hasParentActivity {
				result = true
				break
			}
		}
	}

	return result
}

func (s *Session) uploadAttachments() error {
	for i, a := range s.Attachments {
		a = a.Prepare()
		result, err := a.upload(s)
		if err != nil {
			return err
		}
		a.Location = result.Location
		s.Attachments[i] = a
	}

	return nil
}

func (s *Session) lookupCountry(rawIP string) error {
	ip := net.ParseIP(rawIP)
	if inet.Isv4(ip) {
		s.IPv4 = ip
	} else {
		s.IPv6 = ip
	}

	country, err := inet.LookupCountry(rawIP)
	if err != nil {
		fmt.Println("failed to lookup country by ip")
		return err
	}

	bogon, err := ipinfo.GetIPBogon(ip)
	if err != nil {
		fmt.Println("failed to lookup bogon ip")
	}

	if bogon {
		s.CountryCode = "bogon"
	} else if *country != "" {
		s.CountryCode = *country
	} else {
		s.CountryCode = "not available"
	}

	return nil
}

func (s *Session) getUnhandledExceptions() []EventField {
	var exceptions []EventField
	for _, event := range s.Events {
		if !event.isException() {
			continue
		}
		if event.Exception.Handled {
			continue
		}
		exceptions = append(exceptions, event)
	}

	return exceptions
}

func (s *Session) getANRs() []EventField {
	var anrs []EventField
	for _, event := range s.Events {
		if !event.isANR() {
			continue
		}
		anrs = append(anrs, event)
	}

	return anrs
}

func (s *Session) bucketUnhandledException() error {
	exceptions := s.getUnhandledExceptions()

	type EventGroup struct {
		eventId     uuid.UUID
		exception   Exception
		fingerprint uint64
	}

	var groups []EventGroup

	for _, event := range exceptions {
		if event.Exception.Fingerprint == "" {
			msg := fmt.Sprintf("fingerprint for event %q is empty, cannot bucket", event.ID)
			fmt.Println(msg)
			continue
		}

		fingerprint, err := strconv.ParseUint(event.Exception.Fingerprint, 16, 64)
		if err != nil {
			msg := fmt.Sprintf("failed to parse fingerprint as integer for event %q", event.ID)
			fmt.Println(msg, err)
			return err
		}

		groups = append(groups, EventGroup{
			eventId:     event.ID,
			exception:   event.Exception,
			fingerprint: fingerprint,
		})
	}

	app := App{
		ID: &s.AppID,
	}

	for _, group := range groups {
		appExceptionGroups, err := app.GetExceptionGroups(nil)
		if err != nil {
			return err
		}

		if len(appExceptionGroups) < 1 {
			// insert new exception group
			return NewExceptionGroup(s.AppID, s.Resource.AppVersion, group.exception.getType(), fmt.Sprintf("%x", group.fingerprint), []uuid.UUID{group.eventId}).Insert()
		}

		index, err := ClosestExceptionGroup(appExceptionGroups, group.fingerprint)
		if err != nil {
			return err
		}
		if index < 0 {
			// when no group matches exists, create new exception group
			NewExceptionGroup(s.AppID, s.Resource.AppVersion, group.exception.getType(), fmt.Sprintf("%x", group.fingerprint), []uuid.UUID{group.eventId}).Insert()
			continue
		}
		matchedGroup := appExceptionGroups[index]

		if matchedGroup.EventExists(group.eventId) {
			continue
		}

		if err := matchedGroup.AppendEventId(group.eventId); err != nil {
			return err
		}
	}

	return nil
}

func (s *Session) bucketANRs() error {
	anrs := s.getANRs()

	type EventGroup struct {
		eventId     uuid.UUID
		anr         ANR
		fingerprint uint64
	}

	var groups []EventGroup

	for _, event := range anrs {
		if event.ANR.Fingerprint == "" {
			msg := fmt.Sprintf("fingerprint for anr event %q is empty, cannot bucket", event.ID)
			fmt.Println(msg)
			continue
		}

		fingerprint, err := strconv.ParseUint(event.ANR.Fingerprint, 16, 64)
		if err != nil {
			msg := fmt.Sprintf("failed to parse fingerprint as integer for anr event %q", event.ID)
			fmt.Println(msg, err)
			return err
		}

		groups = append(groups, EventGroup{
			eventId:     event.ID,
			anr:         event.ANR,
			fingerprint: fingerprint,
		})
	}

	app := App{
		ID: &s.AppID,
	}

	for _, group := range groups {
		appANRGroups, err := app.GetANRGroups(nil)
		if err != nil {
			return err
		}

		if len(appANRGroups) < 1 {
			// insert new anr group
			return NewANRGroup(s.AppID, s.Resource.AppVersion, group.anr.getType(), fmt.Sprintf("%x", group.fingerprint), []uuid.UUID{group.eventId}).Insert()
		}

		index, err := ClosestANRGroup(appANRGroups, group.fingerprint)
		if err != nil {
			return err
		}
		if index < 0 {
			// when no group matches exists, create new anr group
			NewANRGroup(s.AppID, s.Resource.AppVersion, group.anr.getType(), fmt.Sprintf("%x", group.fingerprint), []uuid.UUID{group.eventId}).Insert()
			continue
		}
		matchedGroup := appANRGroups[index]

		if matchedGroup.EventExists(group.eventId) {
			continue
		}

		if err := matchedGroup.AppendEventId(group.eventId); err != nil {
			return err
		}
	}

	return nil
}

func (s *Session) ingest() error {
	stmt := sqlf.InsertInto("default.events")
	defer stmt.Close()

	empty := false

	if len(s.Events) == 0 {
		empty = true
		s.Events = append(s.Events, EventField{})
	}

	var args []any
	for i := range s.Events {
		anrExceptions := "[]"
		anrThreads := "[]"
		exceptionExceptions := "[]"
		exceptionThreads := "[]"
		isLowMemory := false
		if s.Events[i].isANR() {
			anrExceptions = s.Events[i].ANR.Exceptions.encode()
			anrThreads = s.Events[i].ANR.Threads.encode()
			s.Events[i].computeANRFingerprint()
		}
		if s.Events[i].isException() {
			exceptionExceptions = s.Events[i].Exception.Exceptions.encode()
			exceptionThreads = s.Events[i].Exception.Threads.encode()
			s.Events[i].computeExceptionFingerprint()
		}
		if s.Events[i].isLowMemory() {
			isLowMemory = true
		}
		if !empty {
			s.Events[i].ID = uuid.New()
		}
		stmt.NewRow().
			Set("id", nil).
			Set("type", nil).
			Set("session_id", nil).
			Set("app_id", nil).
			Set("inet.ipv4", nil).
			Set("inet.ipv6", nil).
			Set("inet.country_code", nil).
			Set("timestamp", nil).
			Set("thread_name", nil).
			Set("resource.device_name", nil).
			Set("resource.device_model", nil).
			Set("resource.device_manufacturer", nil).
			Set("resource.device_type", nil).
			Set("resource.device_is_foldable", nil).
			Set("resource.device_is_physical", nil).
			Set("resource.device_density_dpi", nil).
			Set("resource.device_width_px", nil).
			Set("resource.device_height_px", nil).
			Set("resource.device_density", nil).
			Set("resource.os_name", nil).
			Set("resource.os_version", nil).
			Set("resource.platform", nil).
			Set("resource.app_version", nil).
			Set("resource.app_build", nil).
			Set("resource.app_unique_id", nil).
			Set("resource.measure_sdk_version", nil).
			Set("anr.thread_name", nil).
			Set("anr.handled", nil).
			Set("anr.fingerprint", nil).
			Set("anr_exceptions", nil).
			Set("anr_threads", nil).
			Set("exception.thread_name", nil).
			Set("exception.handled", nil).
			Set("exception.fingerprint", nil).
			Set("exception_exceptions", nil).
			Set("exception_threads", nil).
			Set("app_exit.reason", nil).
			Set("app_exit.importance", nil).
			Set("app_exit.trace", nil).
			Set("app_exit.process_name", nil).
			Set("app_exit.pid", nil).
			Set("app_exit.timestamp", nil).
			Set("string.severity_text", nil).
			Set("string.string", nil).
			Set("gesture_long_click.target", nil).
			Set("gesture_long_click.target_id", nil).
			Set("gesture_long_click.touch_down_time", nil).
			Set("gesture_long_click.touch_up_time", nil).
			Set("gesture_long_click.width", nil).
			Set("gesture_long_click.height", nil).
			Set("gesture_long_click.x", nil).
			Set("gesture_long_click.y", nil).
			Set("gesture_click.target", nil).
			Set("gesture_click.target_id", nil).
			Set("gesture_click.touch_down_time", nil).
			Set("gesture_click.touch_up_time", nil).
			Set("gesture_click.width", nil).
			Set("gesture_click.height", nil).
			Set("gesture_click.x", nil).
			Set("gesture_click.y", nil).
			Set("gesture_scroll.target", nil).
			Set("gesture_scroll.target_id", nil).
			Set("gesture_scroll.touch_down_time", nil).
			Set("gesture_scroll.touch_up_time", nil).
			Set("gesture_scroll.x", nil).
			Set("gesture_scroll.y", nil).
			Set("gesture_scroll.end_x", nil).
			Set("gesture_scroll.end_y", nil).
			Set("gesture_scroll.direction", nil).
			Set("lifecycle_activity.type", nil).
			Set("lifecycle_activity.class_name", nil).
			Set("lifecycle_activity.intent", nil).
			Set("lifecycle_activity.saved_instance_state", nil).
			Set("lifecycle_fragment.type", nil).
			Set("lifecycle_fragment.class_name", nil).
			Set("lifecycle_fragment.parent_activity", nil).
			Set("lifecycle_fragment.tag", nil).
			Set("lifecycle_app.type", nil).
			Set("cold_launch.process_start_uptime", nil).
			Set("cold_launch.process_start_requested_uptime", nil).
			Set("cold_launch.content_provider_attach_uptime", nil).
			Set("cold_launch.on_next_draw_uptime", nil).
			Set("cold_launch.launched_activity", nil).
			Set("cold_launch.has_saved_state", nil).
			Set("cold_launch.intent_data", nil).
			Set("warm_launch.app_visible_uptime", nil).
			Set("warm_launch.on_next_draw_uptime", nil).
			Set("warm_launch.launched_activity", nil).
			Set("warm_launch.has_saved_state", nil).
			Set("warm_launch.intent_data", nil).
			Set("hot_launch.app_visible_uptime", nil).
			Set("hot_launch.on_next_draw_uptime", nil).
			Set("hot_launch.launched_activity", nil).
			Set("hot_launch.has_saved_state", nil).
			Set("hot_launch.intent_data", nil).
			Set("attributes", nil).
			Set("network_change.network_type", nil).
			Set("network_change.previous_network_type", nil).
			Set("network_change.network_generation", nil).
			Set("network_change.previous_network_generation", nil).
			Set("network_change.network_provider", nil).
			Set("anr.network_type", nil).
			Set("anr.network_generation", nil).
			Set("anr.network_provider", nil).
			Set("exception.network_type", nil).
			Set("exception.network_generation", nil).
			Set("exception.network_provider", nil).
			Set("resource.network_type", nil).
			Set("resource.network_generation", nil).
			Set("resource.network_provider", nil).
			Set("resource.device_locale", nil).
			Set("anr.device_locale", nil).
			Set("exception.device_locale", nil).
			Set("http.url", nil).
			Set("http.method", nil).
			Set("http.status_code", nil).
			Set("http.request_body_size", nil).
			Set("http.response_body_size", nil).
			Set("http.request_timestamp", nil).
			Set("http.response_timestamp", nil).
			Set("http.start_time", nil).
			Set("http.end_time", nil).
			Set("http.dns_start", nil).
			Set("http.dns_end", nil).
			Set("http.connect_start", nil).
			Set("http.connect_end", nil).
			Set("http.request_start", nil).
			Set("http.request_end", nil).
			Set("http.request_headers_start", nil).
			Set("http.request_headers_end", nil).
			Set("http.request_body_start", nil).
			Set("http.request_body_end", nil).
			Set("http.response_start", nil).
			Set("http.response_end", nil).
			Set("http.response_headers_start", nil).
			Set("http.response_headers_end", nil).
			Set("http.response_body_start", nil).
			Set("http.response_body_end", nil).
			Set("http.request_headers_size", nil).
			Set("http.response_headers_size", nil).
			Set("http.failure_reason", nil).
			Set("http.failure_description", nil).
			Set("http_request_headers", nil).
			Set("http_response_headers", nil).
			Set("http.client", nil).
			Set("memory_usage.java_max_heap", nil).
			Set("memory_usage.java_total_heap", nil).
			Set("memory_usage.java_free_heap", nil).
			Set("memory_usage.total_pss", nil).
			Set("memory_usage.rss", nil).
			Set("memory_usage.native_total_heap", nil).
			Set("memory_usage.native_free_heap", nil).
			Set("memory_usage.interval_config", nil).
			Set("low_memory", nil).
			Set("trim_memory.level", nil).
			Set("cpu_usage.num_cores", nil).
			Set("cpu_usage.clock_speed", nil).
			Set("cpu_usage.start_time", nil).
			Set("cpu_usage.uptime", nil).
			Set("cpu_usage.utime", nil).
			Set("cpu_usage.cutime", nil).
			Set("cpu_usage.stime", nil).
			Set("cpu_usage.cstime", nil).
			Set("cpu_usage.interval_config", nil)

		args = append(args,
			s.Events[i].ID,
			s.Events[i].Type,
			s.SessionID,
			s.AppID,
			s.IPv4,
			s.IPv6,
			s.CountryCode,
			s.Events[i].Timestamp.Format(timeFormat),
			s.Events[i].ThreadName,
			s.Resource.DeviceName,
			s.Resource.DeviceModel,
			s.Resource.DeviceManufacturer,
			s.Resource.DeviceType,
			s.Resource.DeviceIsFoldable,
			s.Resource.DeviceIsPhysical,
			s.Resource.DeviceDensityDPI,
			s.Resource.DeviceWidthPX,
			s.Resource.DeviceHeightPX,
			s.Resource.DeviceDensity,
			s.Resource.OSName,
			s.Resource.OSVersion,
			s.Resource.Platform,
			s.Resource.AppVersion,
			s.Resource.AppBuild,
			s.Resource.AppUniqueID,
			s.Resource.MeasureSDKVersion,
			s.Events[i].ANR.ThreadName,
			s.Events[i].ANR.Handled,
			s.Events[i].ANR.Fingerprint,
			anrExceptions,
			anrThreads,
			s.Events[i].Exception.ThreadName,
			s.Events[i].Exception.Handled,
			s.Events[i].Exception.Fingerprint,
			exceptionExceptions,
			exceptionThreads,
			s.Events[i].AppExit.Reason,
			s.Events[i].AppExit.Importance,
			s.Events[i].AppExit.Trace,
			s.Events[i].AppExit.ProcessName,
			s.Events[i].AppExit.PID,
			s.Events[i].AppExit.Timestamp,
			s.Events[i].LogString.SeverityText,
			s.Events[i].LogString.String,
			s.Events[i].GestureLongClick.Target,
			s.Events[i].GestureLongClick.TargetID,
			s.Events[i].GestureLongClick.TouchDownTime,
			s.Events[i].GestureLongClick.TouchUpTime,
			s.Events[i].GestureLongClick.Width,
			s.Events[i].GestureLongClick.Height,
			s.Events[i].GestureLongClick.X,
			s.Events[i].GestureLongClick.Y,
			s.Events[i].GestureClick.Target,
			s.Events[i].GestureClick.TargetID,
			s.Events[i].GestureClick.TouchDownTime,
			s.Events[i].GestureClick.TouchUpTime,
			s.Events[i].GestureClick.Width,
			s.Events[i].GestureClick.Height,
			s.Events[i].GestureClick.X,
			s.Events[i].GestureClick.Y,
			s.Events[i].GestureScroll.Target,
			s.Events[i].GestureScroll.TargetID,
			s.Events[i].GestureScroll.TouchDownTime,
			s.Events[i].GestureScroll.TouchUpTime,
			s.Events[i].GestureScroll.X,
			s.Events[i].GestureScroll.Y,
			s.Events[i].GestureScroll.EndX,
			s.Events[i].GestureScroll.EndY,
			s.Events[i].GestureScroll.Direction,
			s.Events[i].LifecycleActivity.Type,
			s.Events[i].LifecycleActivity.ClassName,
			s.Events[i].LifecycleActivity.Intent,
			s.Events[i].LifecycleActivity.SavedInstanceState,
			s.Events[i].LifecycleFragment.Type,
			s.Events[i].LifecycleFragment.ClassName,
			s.Events[i].LifecycleFragment.ParentActivity,
			s.Events[i].LifecycleFragment.Tag,
			s.Events[i].LifecycleApp.Type,
			s.Events[i].ColdLaunch.ProcessStartUptime,
			s.Events[i].ColdLaunch.ProcessStartRequestedUptime,
			s.Events[i].ColdLaunch.ContentProviderAttachUptime,
			s.Events[i].ColdLaunch.OnNextDrawUptime,
			s.Events[i].ColdLaunch.LaunchedActivity,
			s.Events[i].ColdLaunch.HasSavedState,
			s.Events[i].ColdLaunch.IntentData,
			s.Events[i].WarmLaunch.AppVisibleUptime,
			s.Events[i].WarmLaunch.OnNextDrawUptime,
			s.Events[i].WarmLaunch.LaunchedActivity,
			s.Events[i].WarmLaunch.HasSavedState,
			s.Events[i].WarmLaunch.IntentData,
			s.Events[i].HotLaunch.AppVisibleUptime,
			s.Events[i].HotLaunch.OnNextDrawUptime,
			s.Events[i].HotLaunch.LaunchedActivity,
			s.Events[i].HotLaunch.HasSavedState,
			s.Events[i].HotLaunch.IntentData,
			s.Events[i].Attributes,
			s.Events[i].NetworkChange.NetworkType,
			s.Events[i].NetworkChange.PreviousNetworkType,
			s.Events[i].NetworkChange.NetworkGeneration,
			s.Events[i].NetworkChange.PreviousNetworkGeneration,
			s.Events[i].NetworkChange.NetworkProvider,
			s.Events[i].ANR.NetworkType,
			s.Events[i].ANR.NetworkGeneration,
			s.Events[i].ANR.NetworkProvider,
			s.Events[i].Exception.NetworkType,
			s.Events[i].Exception.NetworkGeneration,
			s.Events[i].Exception.NetworkProvider,
			s.Resource.NetworkType,
			s.Resource.NetworkGeneration,
			s.Resource.NetworkProvider,
			s.Resource.DeviceLocale,
			s.Events[i].ANR.DeviceLocale,
			s.Events[i].Exception.DeviceLocale,
			s.Events[i].Http.URL,
			s.Events[i].Http.Method,
			s.Events[i].Http.StatusCode,
			s.Events[i].Http.RequestBodySize,
			s.Events[i].Http.ResponseBodySize,
			s.Events[i].Http.RequestTimestamp,
			s.Events[i].Http.ResponseTimestamp,
			s.Events[i].Http.StartTime,
			s.Events[i].Http.EndTime,
			s.Events[i].Http.DNSStart,
			s.Events[i].Http.DNSEnd,
			s.Events[i].Http.ConnectStart,
			s.Events[i].Http.ConnectEnd,
			s.Events[i].Http.RequestStart,
			s.Events[i].Http.RequestEnd,
			s.Events[i].Http.RequestHeadersStart,
			s.Events[i].Http.RequestHeadersEnd,
			s.Events[i].Http.RequestBodyStart,
			s.Events[i].Http.RequestBodyEnd,
			s.Events[i].Http.ResponseStart,
			s.Events[i].Http.ResponseEnd,
			s.Events[i].Http.ResponseHeadersStart,
			s.Events[i].Http.ResponseHeadersEnd,
			s.Events[i].Http.ResponseBodyStart,
			s.Events[i].Http.ResponseBodyEnd,
			s.Events[i].Http.RequestHeadersSize,
			s.Events[i].Http.ResponseHeadersSize,
			s.Events[i].Http.FailureReason,
			s.Events[i].Http.FailureDescription,
			s.Events[i].Http.RequestHeaders,
			s.Events[i].Http.ResponseHeaders,
			s.Events[i].Http.Client,
			s.Events[i].MemoryUsage.JavaMaxHeap,
			s.Events[i].MemoryUsage.JavaTotalHeap,
			s.Events[i].MemoryUsage.JavaFreeHeap,
			s.Events[i].MemoryUsage.TotalPSS,
			s.Events[i].MemoryUsage.RSS,
			s.Events[i].MemoryUsage.NativeTotalHeap,
			s.Events[i].MemoryUsage.NativeFreeHeap,
			s.Events[i].MemoryUsage.IntervalConfig,
			isLowMemory,
			s.Events[i].TrimMemory.Level,
			s.Events[i].CPUUsage.NumCores,
			s.Events[i].CPUUsage.ClockSpeed,
			s.Events[i].CPUUsage.StartTime,
			s.Events[i].CPUUsage.Uptime,
			s.Events[i].CPUUsage.UTime,
			s.Events[i].CPUUsage.CUTime,
			s.Events[i].CPUUsage.STime,
			s.Events[i].CPUUsage.CSTime,
			s.Events[i].CPUUsage.IntervalConfig,
		)
	}

	if err := server.Server.ChPool.AsyncInsert(context.Background(), stmt.String(), false, args...); err != nil {
		return err
	}

	// keep the emptiness alive
	if empty {
		s.Events = []EventField{}
	}

	return nil
}

func (s *Session) saveWithContext(c *gin.Context) error {
	bytesIn := c.MustGet("bytesIn")
	appId, err := uuid.Parse(c.GetString("appId"))
	if err != nil {
		msg := "error parsing app's uuid"
		fmt.Println(msg, err)
		return err
	}

	app := &App{
		ID: &appId,
	}
	if app, err = app.get(); err != nil {
		msg := "failed to get app"
		fmt.Println(msg, err)
		return err
	}

	tx, err := server.Server.PgPool.Begin(context.Background())
	if err != nil {
		return err
	}

	defer tx.Rollback(context.Background())
	now := time.Now()

	stmt := sqlf.PostgreSQL.InsertInto("public.sessions").
		Set("id", nil).
		Set("event_count", nil).
		Set("attachment_count", nil).
		Set("bytes_in", nil).
		Set("timestamp", nil).
		Set("app_id", nil).
		Set("created_at", nil).
		Set("updated_at", nil)

	defer stmt.Close()

	// insert the session
	_, err = tx.Exec(context.Background(), stmt.String(), s.SessionID, len(s.Events), len(s.Attachments), bytesIn, s.Timestamp, appId, now, now)
	if err != nil {
		fmt.Println(`failed to write session to db`, err.Error())
		return err
	}

	// insert attachments, if present
	if s.hasAttachments() {
		stmt := sqlf.PostgreSQL.InsertInto("public.sessions_attachments")
		defer stmt.Close()
		var args []any
		for _, a := range s.Attachments {
			stmt.NewRow().
				Set("id", nil).
				Set("session_id", nil).
				Set("name", nil).
				Set("extension", nil).
				Set("type", nil).
				Set("key", nil).
				Set("location", nil).
				Set("timestamp", nil)
			args = append(args, a.ID, s.SessionID, a.Name, a.Extension, a.Type, a.Key, a.Location, a.Timestamp)
		}
		_, err := tx.Exec(context.Background(), stmt.String(), args...)
		if err != nil {
			return err
		}
	}

	if !app.Onboarded {
		uniqueIdentifier := s.Resource.AppUniqueID
		platform := s.Resource.Platform
		firstVersion := s.Resource.AppVersion

		if err := app.Onboard(tx, uniqueIdentifier, platform, firstVersion); err != nil {
			return err
		}
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func (s *Session) known() (bool, error) {
	var known string
	if err := server.Server.PgPool.QueryRow(context.Background(), `select id from sessions where id = $1 and app_id = $2;`, s.SessionID, s.AppID).Scan(&known); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *Session) getMappingKey() (string, error) {
	var key string
	if err := server.Server.PgPool.QueryRow(context.Background(), `select key from mapping_files where app_unique_id = $1 and version_name = $2 and version_code = $3 and mapping_type = 'proguard' limit 1;`, s.Resource.AppUniqueID, s.Resource.AppVersion, s.Resource.AppBuild).Scan(&key); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", nil
		}
		return "", err
	}
	return key, nil
}

func (s *Session) EncodeForSymbolication() (CodecMap, []SymbolicationUnit) {
	var symbolicationUnits []SymbolicationUnit
	codecMap := make(CodecMap)

	for eventIdx, event := range s.Events {
		if event.isException() {
			for exceptionIdx, ex := range event.Exception.Exceptions {
				if len(ex.Frames) > 0 {
					idException := uuid.New()
					unitEx := NewCodecMapVal()
					unitEx.Type = TypeException
					unitEx.Event = eventIdx
					unitEx.Exception = exceptionIdx
					unitEx.Frames = TransformSwap
					codecMap[idException] = *unitEx
					su := new(SymbolicationUnit)
					su.ID = idException
					for _, frame := range ex.Frames {
						su.Values = append(su.Values, MarshalRetraceFrame(frame, FramePrefix))
					}
					symbolicationUnits = append(symbolicationUnits, *su)
				}
				if len(ex.Type) > 0 {
					idExceptionType := uuid.New()
					unitExType := NewCodecMapVal()
					unitExType.Type = TypeException
					unitExType.Event = eventIdx
					unitExType.Exception = exceptionIdx
					unitExType.ExceptionType = TransformSwap
					codecMap[idExceptionType] = *unitExType
					su := new(SymbolicationUnit)
					su.ID = idExceptionType
					su.Values = []string{GenericPrefix + ex.Type}
					symbolicationUnits = append(symbolicationUnits, *su)
				}
			}
			for threadIdx, ex := range event.Exception.Threads {
				if len(ex.Frames) > 0 {
					idThread := uuid.New()
					unitTh := NewCodecMapVal()
					unitTh.Type = TypeException
					unitTh.Event = eventIdx
					unitTh.Thread = threadIdx
					unitTh.Frames = TransformSwap
					codecMap[idThread] = *unitTh
					su := new(SymbolicationUnit)
					su.ID = idThread
					for _, frame := range ex.Frames {
						su.Values = append(su.Values, MarshalRetraceFrame(frame, FramePrefix))
					}
					symbolicationUnits = append(symbolicationUnits, *su)
				}
			}
		}

		if event.isANR() {
			for exceptionIdx, ex := range event.ANR.Exceptions {
				if len(ex.Frames) > 0 {
					idException := uuid.New()
					unitEx := NewCodecMapVal()
					unitEx.Type = TypeANR
					unitEx.Event = eventIdx
					unitEx.Exception = exceptionIdx
					unitEx.Frames = TransformSwap
					codecMap[idException] = *unitEx
					su := new(SymbolicationUnit)
					su.ID = idException
					for _, frame := range ex.Frames {
						su.Values = append(su.Values, MarshalRetraceFrame(frame, FramePrefix))
					}
					symbolicationUnits = append(symbolicationUnits, *su)
				}
				if len(ex.Type) > 0 {
					idExceptionType := uuid.New()
					unitExType := NewCodecMapVal()
					unitExType.Type = TypeANR
					unitExType.Event = eventIdx
					unitExType.Exception = exceptionIdx
					unitExType.ExceptionType = TransformSwap
					codecMap[idExceptionType] = *unitExType
					su := new(SymbolicationUnit)
					su.ID = idExceptionType
					su.Values = []string{GenericPrefix + ex.Type}
					symbolicationUnits = append(symbolicationUnits, *su)
				}
			}
			for threadIdx, ex := range event.ANR.Threads {
				if len(ex.Frames) > 0 {
					idThread := uuid.New()
					unitTh := NewCodecMapVal()
					unitTh.Type = TypeANR
					unitTh.Event = eventIdx
					unitTh.Thread = threadIdx
					unitTh.Frames = TransformSwap
					codecMap[idThread] = *unitTh
					su := new(SymbolicationUnit)
					su.ID = idThread
					for _, frame := range ex.Frames {
						su.Values = append(su.Values, MarshalRetraceFrame(frame, FramePrefix))
					}
					symbolicationUnits = append(symbolicationUnits, *su)
				}
			}
		}

		if event.isAppExit() {
			if len(event.AppExit.Trace) > 0 {
				idAppExit := uuid.New()
				unitAE := NewCodecMapVal()
				unitAE.Type = TypeAppExit
				unitAE.Event = eventIdx
				unitAE.Trace = TransformSwap
				codecMap[idAppExit] = *unitAE
				su := new(SymbolicationUnit)
				su.ID = idAppExit
				su.Values = []string{GenericPrefix + event.AppExit.Trace}
				symbolicationUnits = append(symbolicationUnits, *su)
			}
		}

		if event.isLifecycleActivity() {
			if len(event.LifecycleActivity.ClassName) > 0 {
				idLifecycleActivity := uuid.New()
				unitLA := NewCodecMapVal()
				unitLA.Type = TypeLifecycleActivity
				unitLA.Event = eventIdx
				unitLA.ClassName = TransformSwap
				codecMap[idLifecycleActivity] = *unitLA
				su := new(SymbolicationUnit)
				su.ID = idLifecycleActivity
				su.Values = []string{GenericPrefix + event.LifecycleActivity.ClassName}
				symbolicationUnits = append(symbolicationUnits, *su)
			}
		}

		if event.isLifecycleFragment() {
			if len(event.LifecycleFragment.ClassName) > 0 {
				idLifecycleFragment := uuid.New()
				unitLF := NewCodecMapVal()
				unitLF.Type = TypeLifecycleFragment
				unitLF.Event = eventIdx
				unitLF.ClassName = TransformSwap
				codecMap[idLifecycleFragment] = *unitLF
				su := new(SymbolicationUnit)
				su.ID = idLifecycleFragment
				su.Values = []string{GenericPrefix + event.LifecycleFragment.ClassName}
				symbolicationUnits = append(symbolicationUnits, *su)
			}

			if len(event.LifecycleFragment.ParentActivity) > 0 {
				idLifecycleFragment := uuid.New()
				unitLF := NewCodecMapVal()
				unitLF.Type = TypeLifecycleFragment
				unitLF.Event = eventIdx
				unitLF.ParentActivity = TransformSwap
				codecMap[idLifecycleFragment] = *unitLF
				su := new(SymbolicationUnit)
				su.ID = idLifecycleFragment
				su.Values = []string{GenericPrefix + event.LifecycleFragment.ParentActivity}
				symbolicationUnits = append(symbolicationUnits, *su)
			}
		}

		if event.isColdLaunch() {
			if len(event.ColdLaunch.LaunchedActivity) > 0 {
				idColdLaunch := uuid.New()
				unitCL := NewCodecMapVal()
				unitCL.Type = TypeColdLaunch
				unitCL.Event = eventIdx
				unitCL.LaunchedActivity = TransformSwap
				codecMap[idColdLaunch] = *unitCL
				su := new(SymbolicationUnit)
				su.ID = idColdLaunch
				su.Values = []string{GenericPrefix + event.ColdLaunch.LaunchedActivity}
				symbolicationUnits = append(symbolicationUnits, *su)
			}
		}
		if event.isWarmLaunch() {
			if len(event.WarmLaunch.LaunchedActivity) > 0 {
				idWarmLaunch := uuid.New()
				unitCL := NewCodecMapVal()
				unitCL.Type = TypeWarmLaunch
				unitCL.Event = eventIdx
				unitCL.LaunchedActivity = TransformSwap
				codecMap[idWarmLaunch] = *unitCL
				su := new(SymbolicationUnit)
				su.ID = idWarmLaunch
				su.Values = []string{GenericPrefix + event.WarmLaunch.LaunchedActivity}
				symbolicationUnits = append(symbolicationUnits, *su)
			}
		}
		if event.isHotLaunch() {
			if len(event.HotLaunch.LaunchedActivity) > 0 {
				idHotLaunch := uuid.New()
				unitCL := NewCodecMapVal()
				unitCL.Type = TypeHotLaunch
				unitCL.Event = eventIdx
				unitCL.LaunchedActivity = TransformSwap
				codecMap[idHotLaunch] = *unitCL
				su := new(SymbolicationUnit)
				su.ID = idHotLaunch
				su.Values = []string{GenericPrefix + event.HotLaunch.LaunchedActivity}
				symbolicationUnits = append(symbolicationUnits, *su)
			}
		}
	}

	return codecMap, symbolicationUnits
}

func (s *Session) DecodeFromSymbolication(codecMap CodecMap, symbolicationUnits []SymbolicationUnit) {
	for _, su := range symbolicationUnits {
		codecMapVal := codecMap[su.ID]
		switch codecMapVal.Type {
		case TypeException:
			if codecMapVal.Frames == TransformSwap {
				if codecMapVal.Exception > -1 {
					var frames Frames
					for _, value := range su.Values {
						frame, err := UnmarshalRetraceFrame(value, FramePrefix)
						if err != nil {
							fmt.Println("failed to unmarshal retrace frame", err)
							continue
						}
						frames = append(frames, Frame{
							ClassName:  frame.ClassName,
							LineNum:    frame.LineNum,
							FileName:   frame.FileName,
							MethodName: frame.MethodName,
						})
					}
					s.Events[codecMapVal.Event].Exception.Exceptions[codecMapVal.Exception].Frames = frames
				}

				if codecMapVal.Thread > -1 {
					var frames Frames
					for _, value := range su.Values {
						frame, err := UnmarshalRetraceFrame(value, FramePrefix)
						if err != nil {
							fmt.Println("failed to unmarshal retrace frame", err)
							continue
						}
						frames = append(frames, Frame{
							ClassName:  frame.ClassName,
							LineNum:    frame.LineNum,
							FileName:   frame.FileName,
							MethodName: frame.MethodName,
						})
					}
					s.Events[codecMapVal.Event].Exception.Threads[codecMapVal.Thread].Frames = frames
				}
			}

			if codecMapVal.ExceptionType == TransformSwap {
				exceptionType := strings.TrimPrefix(su.Values[0], GenericPrefix)
				s.Events[codecMapVal.Event].Exception.Exceptions[codecMapVal.Exception].Type = exceptionType
			}
		case TypeANR:
			if codecMapVal.Frames == TransformSwap {
				if codecMapVal.Exception > -1 {
					var frames Frames
					for _, value := range su.Values {
						frame, err := UnmarshalRetraceFrame(value, FramePrefix)
						if err != nil {
							fmt.Println("failed to unmarshal retrace frame", err)
							continue
						}
						frames = append(frames, Frame{
							ClassName:  frame.ClassName,
							LineNum:    frame.LineNum,
							FileName:   frame.FileName,
							MethodName: frame.MethodName,
						})
					}
					s.Events[codecMapVal.Event].ANR.Exceptions[codecMapVal.Exception].Frames = frames
				}

				if codecMapVal.Thread > -1 {
					var frames Frames
					for _, value := range su.Values {
						frame, err := UnmarshalRetraceFrame(value, FramePrefix)
						if err != nil {
							fmt.Println("failed to unmarshal retrace frame", err)
							continue
						}
						frames = append(frames, Frame{
							ClassName:  frame.ClassName,
							LineNum:    frame.LineNum,
							FileName:   frame.FileName,
							MethodName: frame.MethodName,
						})
					}
					s.Events[codecMapVal.Event].ANR.Threads[codecMapVal.Thread].Frames = frames
				}
			}

			if codecMapVal.ExceptionType == TransformSwap {
				exceptionType := strings.TrimPrefix(su.Values[0], GenericPrefix)
				s.Events[codecMapVal.Event].ANR.Exceptions[codecMapVal.Exception].Type = exceptionType
			}
		case TypeAppExit:
			if codecMapVal.Trace == TransformSwap {
				s.Events[codecMapVal.Event].AppExit.Trace = strings.TrimPrefix(su.Values[0], GenericPrefix)
			}
		case TypeLifecycleActivity:
			if codecMapVal.ClassName == TransformSwap {
				s.Events[codecMapVal.Event].LifecycleActivity.ClassName = strings.TrimPrefix(su.Values[0], GenericPrefix)
			}
		case TypeLifecycleFragment:
			if codecMapVal.ClassName == TransformSwap {
				s.Events[codecMapVal.Event].LifecycleFragment.ClassName = strings.TrimPrefix(su.Values[0], GenericPrefix)
			}
			if codecMapVal.ParentActivity == TransformSwap {
				s.Events[codecMapVal.Event].LifecycleFragment.ParentActivity = strings.TrimPrefix(su.Values[0], GenericPrefix)
			}
		case TypeColdLaunch:
			if codecMapVal.LaunchedActivity == TransformSwap {
				s.Events[codecMapVal.Event].ColdLaunch.LaunchedActivity = strings.TrimPrefix(su.Values[0], GenericPrefix)
			}
		case TypeWarmLaunch:
			if codecMapVal.LaunchedActivity == TransformSwap {
				s.Events[codecMapVal.Event].WarmLaunch.LaunchedActivity = strings.TrimPrefix(su.Values[0], GenericPrefix)
			}
		case TypeHotLaunch:
			if codecMapVal.LaunchedActivity == TransformSwap {
				s.Events[codecMapVal.Event].HotLaunch.LaunchedActivity = strings.TrimPrefix(su.Values[0], GenericPrefix)
			}

		default:
			continue
		}
	}
}

func PutSession(c *gin.Context) {
	bc := &ByteCounter{}
	c.Request.Body = io.NopCloser(io.TeeReader(c.Request.Body, bc))
	session := new(Session)
	if err := c.ShouldBindJSON(&session); err != nil {
		fmt.Println("gin binding err:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse session payload"})
		return
	}

	appId, err := uuid.Parse(c.GetString("appId"))
	if err != nil {
		msg := "error parsing app's uuid"
		fmt.Println(msg, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to ingest session"})
		return
	}

	session.AppID = appId

	c.Set("bytesIn", bc.Count)

	if known, err := session.known(); err != nil {
		fmt.Println("failed to check existing session", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to ingest session"})
		return
	} else if known {
		c.JSON(http.StatusAccepted, gin.H{"ok": "accepted, known session"})
		return
	}

	if err := session.validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// look up country from ip
	if err := session.lookupCountry(c.ClientIP()); err != nil {
		msg := fmt.Sprintf("failed to lookup country for IP %q", c.ClientIP())
		fmt.Println(msg, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not upload session, failed to lookup country by IP"})
		return

	}

	if session.needsSymbolication() {
		if err := symbolicate(session); err != nil {
			fmt.Println("symbolication failed with error", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not upload session, failed to symbolicate"})
			return
		}
	}

	if session.hasAttachments() {
		if err := session.uploadAttachments(); err != nil {
			fmt.Println("error uploading attachment", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload attachment(s)"})
			return
		}
	}

	if err := session.ingest(); err != nil {
		fmt.Println("clickhouse insert error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := session.saveWithContext(c); err != nil {
		fmt.Println("failed to save session", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save session"})
		return
	}

	if session.hasUnhandledExceptions() {
		if err := session.bucketUnhandledException(); err != nil {
			msg := "failed to save session, error occurred during exception grouping"
			fmt.Println(msg, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
	}

	if session.hasANRs() {
		if err := session.bucketANRs(); err != nil {
			msg := "failed to save session, error occurred during anr grouping"
			fmt.Println(msg, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

	}

	c.JSON(http.StatusAccepted, gin.H{"ok": "accepted"})
}
