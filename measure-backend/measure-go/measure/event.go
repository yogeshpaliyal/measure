package measure

import (
	"encoding/json"
	"fmt"
	"measure-backend/measure-go/chrono"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/go-dedup/simhash"
	"github.com/google/uuid"
)

// maximum character limits for event fields
const (
	maxTypeChars                              = 32
	maxThreadNameChars                        = 64
	maxExceptionDeviceLocaleChars             = 64
	maxAnrDeviceLocaleChars                   = 64
	maxAppExitReasonChars                     = 64
	maxAppExitImportanceChars                 = 32
	maxSeverityTextChars                      = 10
	maxGestureLongClickTargetChars            = 128
	maxGestureLongClickTargetNameChars        = 128
	maxGestureLongClickTargetIDChars          = 128
	maxGestureScrollTargetChars               = 128
	maxGestureScrollTargetNameChars           = 128
	maxGestureScrollTargetIDChars             = 128
	maxGestureScrollDirectionChars            = 8
	maxGestureClickTargetChars                = 128
	maxGestureClickTargetNameChars            = 128
	maxGestureClickTargetIDChars              = 128
	maxLifecycleActivityTypeChars             = 32
	maxLifecycleActivityClassNameChars        = 128
	maxLifecycleFragmentTypeChars             = 32
	maxLifecycleFragmentClassNameChars        = 128
	maxLifecycleAppTypeChars                  = 32
	maxColdLaunchLaunchedActivityChars        = 128
	maxWarmLaunchLaunchedActivityChars        = 128
	maxHotLaunchLaunchedActivityChars         = 128
	maxNetworkChangeNetworkTypeChars          = 16
	maxNetworkChangePreviousNetworkTypeChars  = 16
	maxNetworkChangeNetworkGeneration         = 8
	maxNetworkChangePreviousNetworkGeneration = 8
	maxNetworkChangeNetworkProvider           = 64
	maxHttpMethodChars                        = 16
	maxHttpClientChars                        = 32
	maxTrimMemoryLevelChars                   = 64
	maxAttrCount                              = 10
	maxRouteChars                             = 128
)

const TypeANR = "anr"
const TypeException = "exception"
const TypeAppExit = "app_exit"
const TypeString = "string"
const TypeGestureLongClick = "gesture_long_click"
const TypeGestureClick = "gesture_click"
const TypeGestureScroll = "gesture_scroll"
const TypeLifecycleActivity = "lifecycle_activity"
const TypeLifecycleFragment = "lifecycle_fragment"
const TypeLifecycleApp = "lifecycle_app"
const TypeColdLaunch = "cold_launch"
const TypeWarmLaunch = "warm_launch"
const TypeHotLaunch = "hot_launch"
const TypeNetworkChange = "network_change"
const TypeHttp = "http"
const TypeMemoryUsage = "memory_usage"
const TypeLowMemory = "low_memory"
const TypeTrimMemory = "trim_memory"
const TypeCPUUsage = "cpu_usage"
const TypeNavigation = "navigation"

// timeFormat is the format of datetime in nanoseconds when
// converting datetime values before inserting into database
const timeFormat = "2006-01-02 15:04:05.999999999"

var TrimRight = func(s string) string {
	return strings.TrimRight(s, "\x00")
}

type Frame struct {
	LineNum    int    `json:"line_num"`
	ColNum     int    `json:"col_num"`
	ModuleName string `json:"module_name"`
	FileName   string `json:"file_name"`
	ClassName  string `json:"class_name"`
	MethodName string `json:"method_name"`
}

func (f Frame) String() string {
	className := f.ClassName
	methodName := f.MethodName
	fileName := f.FileName
	var lineNum = ""

	if f.LineNum != 0 {
		lineNum = strconv.Itoa(f.LineNum)
	}

	codeInfo := joinNonEmptyStrings(".", className, methodName)
	fileInfo := joinNonEmptyStrings(":", fileName, lineNum)

	if fileInfo != "" {
		fileInfo = fmt.Sprintf(`(%s)`, fileInfo)
	}

	return fmt.Sprintf(`%s%s`, codeInfo, fileInfo)
}

type Frames []Frame

type ExceptionUnit struct {
	Type    string `json:"type" binding:"required"`
	Message string `json:"message"`
	Frames  Frames `json:"frames" binding:"required"`
}

type ExceptionUnits []ExceptionUnit

type Thread struct {
	Name   string `json:"name" binding:"required"`
	Frames Frames `json:"frames" binding:"required"`
}

type Threads []Thread

type ANR struct {
	ThreadName        string         `json:"thread_name" binding:"required"`
	Handled           bool           `json:"handled" binding:"required"`
	Exceptions        ExceptionUnits `json:"exceptions" binding:"required"`
	Threads           Threads        `json:"threads" binding:"required"`
	NetworkType       string         `json:"network_type"`
	NetworkGeneration string         `json:"network_generation"`
	NetworkProvider   string         `json:"network_provider"`
	DeviceLocale      string         `json:"device_locale"`
	Fingerprint       string         `json:"fingerprint"`
}

type Exception struct {
	ThreadName        string         `json:"thread_name" binding:"required"`
	Handled           bool           `json:"handled" binding:"required"`
	Exceptions        ExceptionUnits `json:"exceptions" binding:"required"`
	Threads           Threads        `json:"threads" binding:"required"`
	NetworkType       string         `json:"network_type"`
	NetworkGeneration string         `json:"network_generation"`
	NetworkProvider   string         `json:"network_provider"`
	DeviceLocale      string         `json:"device_locale"`
	Fingerprint       string         `json:"fingerprint"`
}

func (e *Exception) Trim() {
	e.ThreadName = TrimRight(e.ThreadName)
	e.NetworkType = TrimRight(e.NetworkType)
	e.NetworkGeneration = TrimRight(e.NetworkGeneration)
	e.NetworkProvider = TrimRight(e.NetworkProvider)
	e.DeviceLocale = TrimRight(e.DeviceLocale)
}

func (e Exception) Stacktrace() string {
	var b strings.Builder

	b.WriteString(e.getType() + "\n")

	for i := range e.Exceptions {
		for j := range e.Exceptions[i].Frames {
			frame := e.Exceptions[i].Frames[j].String()
			b.WriteString(FramePrefix + frame + "\n")
		}
	}

	return b.String()
}

func (a *ANR) Trim() {
	a.ThreadName = TrimRight(a.ThreadName)
	a.NetworkType = TrimRight(a.NetworkType)
	a.NetworkGeneration = TrimRight(a.NetworkGeneration)
	a.NetworkProvider = TrimRight(a.NetworkProvider)
	a.DeviceLocale = TrimRight(a.DeviceLocale)
}

func (e ANR) Stacktrace() string {
	var b strings.Builder

	b.WriteString(e.getType() + "\n")

	for i := range e.Exceptions {
		for j := range e.Exceptions[i].Frames {
			frame := e.Exceptions[i].Frames[j].String()
			b.WriteString(FramePrefix + frame + "\n")
		}
	}

	return b.String()
}

type AppExit struct {
	Reason      string    `json:"reason" binding:"required"`
	Importance  string    `json:"importance" binding:"required"`
	Trace       string    `json:"trace"`
	ProcessName string    `json:"process_name" binding:"required"`
	PID         string    `json:"pid" binding:"required"`
	Timestamp   time.Time `json:"timestamp" binding:"required"`
}

type LogString struct {
	SeverityText string `json:"severity_text" binding:"required"`
	String       string `json:"string" binding:"required"`
}

type GestureLongClick struct {
	Target        string  `json:"target"`
	TargetID      string  `json:"target_id"`
	TouchDownTime uint32  `json:"touch_down_time"`
	TouchUpTime   uint32  `json:"touch_up_time"`
	Width         uint16  `json:"width"`
	Height        uint16  `json:"height"`
	X             float32 `json:"x"`
	Y             float32 `json:"y"`
}

type GestureScroll struct {
	Target        string  `json:"target"`
	TargetID      string  `json:"target_id"`
	TouchDownTime uint32  `json:"touch_down_time"`
	TouchUpTime   uint32  `json:"touch_up_time"`
	X             float32 `json:"x"`
	Y             float32 `json:"y"`
	EndX          float32 `json:"end_x"`
	EndY          float32 `json:"end_y"`
	Direction     string  `json:"direction"`
}

type GestureClick struct {
	Target        string  `json:"target"`
	TargetID      string  `json:"target_id"`
	TouchDownTime uint32  `json:"touch_down_time"`
	TouchUpTime   uint32  `json:"touch_up_time"`
	Width         uint16  `json:"width"`
	Height        uint16  `json:"height"`
	X             float32 `json:"x"`
	Y             float32 `json:"y"`
}

type LifecycleActivity struct {
	Type               string `json:"type" binding:"required"`
	ClassName          string `json:"class_name" binding:"required"`
	Intent             string `json:"intent"`
	SavedInstanceState bool   `json:"saved_instance_state"`
}

type LifecycleFragment struct {
	Type           string `json:"type" binding:"required"`
	ClassName      string `json:"class_name" binding:"required"`
	ParentActivity string `json:"parent_activity"`
	Tag            string `json:"tag"`
}

type LifecycleApp struct {
	Type string `json:"type" binding:"required"`
}

type ColdLaunch struct {
	ProcessStartUptime          uint32 `json:"process_start_uptime"`
	ProcessStartRequestedUptime uint32 `json:"process_start_requested_uptime"`
	ContentProviderAttachUptime uint32 `json:"content_provider_attach_uptime"`
	OnNextDrawUptime            uint32 `json:"on_next_draw_uptime" binding:"required"`
	LaunchedActivity            string `json:"launched_activity" binding:"required"`
	HasSavedState               bool   `json:"has_saved_state" binding:"required"`
	IntentData                  string `json:"intent_data"`
}

type WarmLaunch struct {
	AppVisibleUptime uint32 `json:"app_visible_uptime"`
	OnNextDrawUptime uint32 `json:"on_next_draw_uptime" binding:"required"`
	LaunchedActivity string `json:"launched_activity" binding:"required"`
	HasSavedState    bool   `json:"has_saved_state" binding:"required"`
	IntentData       string `json:"intent_data"`
}

type HotLaunch struct {
	AppVisibleUptime uint32 `json:"app_visible_uptime"`
	OnNextDrawUptime uint32 `json:"on_next_draw_uptime" binding:"required"`
	LaunchedActivity string `json:"launched_activity" binding:"required"`
	HasSavedState    bool   `json:"has_saved_state" binding:"required"`
	IntentData       string `json:"intent_data"`
}

type NetworkChange struct {
	NetworkType               string `json:"network_type" binding:"required"`
	PreviousNetworkType       string `json:"previous_network_type"`
	NetworkGeneration         string `json:"network_generation"`
	PreviousNetworkGeneration string `json:"previous_network_generation"`
	NetworkProvider           string `json:"network_provider"`
}

type Http struct {
	URL                  string            `json:"url"`
	Method               string            `json:"method"`
	StatusCode           int               `json:"status_code"`
	RequestBodySize      int               `json:"request_body_size"`
	ResponseBodySize     int               `json:"response_body_size"`
	RequestTimestamp     time.Time         `json:"request_timestamp"`
	ResponseTimestamp    time.Time         `json:"response_timestamp"`
	StartTime            uint64            `json:"start_time"`
	EndTime              uint64            `json:"end_time"`
	DNSStart             uint64            `json:"dns_start"`
	DNSEnd               uint64            `json:"dns_end"`
	ConnectStart         uint64            `json:"connect_start"`
	ConnectEnd           uint64            `json:"connect_end"`
	RequestStart         uint64            `json:"request_start"`
	RequestEnd           uint64            `json:"request_end"`
	RequestHeadersStart  uint64            `json:"request_headers_start"`
	RequestHeadersEnd    uint64            `json:"request_headers_end"`
	RequestBodyStart     uint64            `json:"request_body_start"`
	RequestBodyEnd       uint64            `json:"request_body_end"`
	ResponseStart        uint64            `json:"response_start"`
	ResponseEnd          uint64            `json:"response_end"`
	ResponseHeadersStart uint64            `json:"response_headers_start"`
	ResponseHeadersEnd   uint64            `json:"response_headers_end"`
	ResponseBodyStart    uint64            `json:"response_body_start"`
	ResponseBodyEnd      uint64            `json:"response_body_end"`
	RequestHeadersSize   int               `json:"request_headers_size"`
	ResponseHeadersSize  int               `json:"response_headers_size"`
	FailureReason        string            `json:"failure_reason"`
	FailureDescription   string            `json:"failure_description"`
	RequestHeaders       map[string]string `json:"request_headers"`
	ResponseHeaders      map[string]string `json:"response_headers"`
	Client               string            `json:"client"`
}

type MemoryUsage struct {
	JavaMaxHeap       uint64 `json:"java_max_heap" binding:"required"`
	JavaTotalHeap     uint64 `json:"java_total_heap" binding:"required"`
	JavaFreeHeap      uint64 `json:"java_free_heap" binding:"required"`
	TotalPSS          uint64 `json:"total_pss" binding:"required"`
	RSS               uint64 `json:"rss"`
	NativeTotalHeap   uint64 `json:"native_total_heap" binding:"required"`
	NativeFreeHeap    uint64 `json:"native_free_heap" binding:"required"`
	IntervalConfig    uint64 `json:"interval_config" binding:"required"`
	IntervalStartTime uint32 `json:"interval_start_time" binding:"required"`
}

type LowMemory struct {
}

type TrimMemory struct {
	Level string `json:"level" binding:"required"`
}

type CPUUsage struct {
	NumCores       uint8  `json:"num_cores" binding:"required"`
	ClockSpeed     uint64 `json:"clock_speed" binding:"required"`
	StartTime      uint64 `json:"start_time" binding:"required"`
	Uptime         uint64 `json:"uptime" binding:"required"`
	UTime          uint64 `json:"utime" binding:"required"`
	CUTime         uint64 `json:"cutime" binding:"required"`
	STime          uint64 `json:"stime" binding:"required"`
	CSTime         uint64 `json:"cstime" binding:"required"`
	IntervalConfig uint32 `json:"interval_config" binding:"required"`
}

type Navigation struct {
	Route string `json:"route" binding:"required"`
}

type EventField struct {
	ID                uuid.UUID         `json:"id"`
	Timestamp         time.Time         `json:"timestamp" binding:"required"`
	Type              string            `json:"type" binding:"required"`
	ThreadName        string            `json:"thread_name" binding:"required"`
	Resource          Resource          `json:"resource"`
	ANR               ANR               `json:"anr,omitempty"`
	Exception         Exception         `json:"exception,omitempty"`
	AppExit           AppExit           `json:"app_exit,omitempty"`
	LogString         LogString         `json:"string,omitempty"`
	GestureLongClick  GestureLongClick  `json:"gesture_long_click,omitempty"`
	GestureScroll     GestureScroll     `json:"gesture_scroll,omitempty"`
	GestureClick      GestureClick      `json:"gesture_click,omitempty"`
	LifecycleActivity LifecycleActivity `json:"lifecycle_activity,omitempty"`
	LifecycleFragment LifecycleFragment `json:"lifecycle_fragment,omitempty"`
	LifecycleApp      LifecycleApp      `json:"lifecycle_app,omitempty"`
	ColdLaunch        ColdLaunch        `json:"cold_launch,omitempty"`
	WarmLaunch        WarmLaunch        `json:"warm_launch,omitempty"`
	HotLaunch         HotLaunch         `json:"hot_launch,omitempty"`
	NetworkChange     NetworkChange     `json:"network_change,omitempty"`
	Http              Http              `json:"http,omitempty"`
	MemoryUsage       MemoryUsage       `json:"memory_usage,omitempty"`
	LowMemory         LowMemory         `json:"low_memory,omitempty"`
	TrimMemory        TrimMemory        `json:"trim_memory,omitempty"`
	CPUUsage          CPUUsage          `json:"cpu_usage,omitempty"`
	Navigation        Navigation        `json:"navigation,omitempty"`
	Attributes        map[string]string `json:"attributes"`
}

func (e *EventField) isException() bool {
	return e.Type == TypeException
}

func (e *EventField) isUnhandledException() bool {
	return e.Type == TypeException && !e.Exception.Handled
}

func (e *EventField) isANR() bool {
	return e.Type == TypeANR
}

func (e *EventField) isAppExit() bool {
	return e.Type == TypeAppExit
}

func (e *EventField) isString() bool {
	return e.Type == TypeString
}

func (e *EventField) isGestureLongClick() bool {
	return e.Type == TypeGestureLongClick
}

func (e *EventField) isGestureScroll() bool {
	return e.Type == TypeGestureScroll
}

func (e *EventField) isGestureClick() bool {
	return e.Type == TypeGestureClick
}

func (e *EventField) isLifecycleActivity() bool {
	return e.Type == TypeLifecycleActivity
}

func (e *EventField) isLifecycleFragment() bool {
	return e.Type == TypeLifecycleFragment
}

func (e *EventField) isLifecycleApp() bool {
	return e.Type == TypeLifecycleApp
}

func (e *EventField) isColdLaunch() bool {
	return e.Type == TypeColdLaunch
}

func (e *EventField) isWarmLaunch() bool {
	return e.Type == TypeWarmLaunch
}

func (e *EventField) isHotLaunch() bool {
	return e.Type == TypeHotLaunch
}

func (e *EventField) isNetworkChange() bool {
	return e.Type == TypeNetworkChange
}

func (e *EventField) isHttp() bool {
	return e.Type == TypeHttp
}

func (e *EventField) isMemoryUsage() bool {
	return e.Type == TypeMemoryUsage
}

func (e *EventField) isTrimMemory() bool {
	return e.Type == TypeTrimMemory
}

func (e *EventField) isCPUUsage() bool {
	return e.Type == TypeCPUUsage
}

// check if LowMemory event is present
func (e *EventField) isLowMemory() bool {
	return e.Type == TypeLowMemory
}

func (e *EventField) isNavigation() bool {
	return e.Type == TypeNavigation
}

func (e *EventField) Trim() {
	e.ThreadName = TrimRight(e.ThreadName)
	e.Type = TrimRight(e.Type)
	e.Resource.Trim()
	if e.isException() {
		e.Exception.Trim()
	}
}

type EventException struct {
	ID         uuid.UUID         `json:"id"`
	Timestamp  chrono.ISOTime    `json:"timestamp"`
	Type       string            `json:"type"`
	ThreadName string            `json:"thread_name"`
	Resource   Resource          `json:"resource"`
	Exception  Exception         `json:"-"`
	Exceptions []ExceptionView   `json:"exceptions"`
	Threads    []ThreadView      `json:"threads"`
	Attributes map[string]string `json:"attributes"`
}

type ExceptionView struct {
	Type       string `json:"type"`
	Message    string `json:"message"`
	Location   string `json:"location"`
	Stacktrace string `json:"stacktrace"`
}

type ThreadView struct {
	Name   string   `json:"name"`
	Frames []string `json:"frames"`
}

func (e *EventException) Trim() {
	e.ThreadName = TrimRight(e.ThreadName)
	e.Type = TrimRight(e.Type)
	e.Resource.Trim()
	e.Exception.Trim()
}

func (e *EventException) ComputeView() {
	var ev ExceptionView
	ev.Type = e.Exception.getType()
	ev.Message = e.Exception.getMessage()
	ev.Location = e.Exception.getLocation()
	ev.Stacktrace = e.Exception.Stacktrace()
	e.Exceptions = append(e.Exceptions, ev)

	for i := range e.Exception.Threads {
		var tv ThreadView
		tv.Name = e.Exception.Threads[i].Name
		for j := range e.Exception.Threads[i].Frames {
			tv.Frames = append(tv.Frames, e.Exception.Threads[i].Frames[j].String())
		}
		e.Threads = append(e.Threads, tv)
	}
}

type EventANR struct {
	ID         uuid.UUID         `json:"id"`
	Timestamp  chrono.ISOTime    `json:"timestamp"`
	Type       string            `json:"type"`
	ThreadName string            `json:"thread_name"`
	Resource   Resource          `json:"resource"`
	ANR        ANR               `json:"-"`
	ANRs       []ANRView         `json:"anrs"`
	Threads    []ThreadView      `json:"threads"`
	Attributes map[string]string `json:"attributes"`
}

type ANRView struct {
	Type       string `json:"type"`
	Message    string `json:"message"`
	Location   string `json:"location"`
	Stacktrace string `json:"stacktrace"`
}

func (e *EventANR) Trim() {
	e.ThreadName = TrimRight(e.ThreadName)
	e.Type = TrimRight(e.Type)
	e.Resource.Trim()
	e.ANR.Trim()
}

func (e *EventANR) ComputeView() {
	var av ANRView
	av.Type = e.ANR.getType()
	av.Message = e.ANR.getMessage()
	av.Location = e.ANR.getLocation()
	av.Stacktrace = e.ANR.Stacktrace()
	e.ANRs = append(e.ANRs, av)

	for i := range e.ANR.Threads {
		var tv ThreadView
		tv.Name = e.ANR.Threads[i].Name
		for j := range e.ANR.Threads[i].Frames {
			tv.Frames = append(tv.Frames, e.ANR.Threads[i].Frames[j].String())
		}
		e.Threads = append(e.Threads, tv)
	}
}

func (e *EventField) computeExceptionFingerprint() error {
	if !e.isException() {
		return nil
	}

	if e.Exception.Handled {
		return nil
	}

	marshalledException, err := json.Marshal(e.Exception)
	if err != nil {
		return err
	}

	sh := simhash.NewSimhash()
	e.Exception.Fingerprint = fmt.Sprintf("%x", sh.GetSimhash(sh.NewWordFeatureSet(marshalledException)))

	return nil
}

func (e *EventField) computeANRFingerprint() error {
	if !e.isANR() {
		return nil
	}

	marshalledANR, err := json.Marshal(e.ANR)
	if err != nil {
		return err
	}

	sh := simhash.NewSimhash()
	e.ANR.Fingerprint = fmt.Sprintf("%x", sh.GetSimhash(sh.NewWordFeatureSet(marshalledANR)))
	return nil
}

func (e *EventField) validate() error {
	validTypes := []string{TypeANR, TypeException, TypeAppExit, TypeString, TypeGestureLongClick, TypeGestureScroll, TypeGestureClick, TypeLifecycleActivity, TypeLifecycleFragment, TypeLifecycleApp, TypeColdLaunch, TypeWarmLaunch, TypeHotLaunch, TypeNetworkChange, TypeHttp, TypeMemoryUsage, TypeLowMemory, TypeTrimMemory, TypeCPUUsage, TypeNavigation}
	if !slices.Contains(validTypes, e.Type) {
		return fmt.Errorf(`"events[].type" is not a valid type`)
	}
	if e.Timestamp.IsZero() {
		return fmt.Errorf(`events[].timestamp is invalid. Must be a valid ISO 8601 timestamp`)
	}
	if e.ThreadName == "" {
		return fmt.Errorf(`events[].thread_name is invalid`)
	}
	// validate all required fields of each type
	if e.isANR() {
		if len(e.ANR.Exceptions) < 1 || len(e.ANR.Threads) < 1 || e.ANR.ThreadName == "" {
			return fmt.Errorf(`anr event is invalid`)
		}
	}

	if e.isException() {
		if len(e.Exception.Exceptions) < 1 || len(e.Exception.Threads) < 1 || e.Exception.ThreadName == "" {
			return fmt.Errorf(`exception event is invalid`)
		}
	}

	if e.isAppExit() {
		if len(e.AppExit.Reason) < 1 || len(e.AppExit.Importance) < 1 || len(e.AppExit.ProcessName) < 1 || len(e.AppExit.ProcessName) < 1 || e.AppExit.Timestamp.IsZero() {
			return fmt.Errorf(`app_exit event is invalid`)
		}
	}

	if e.isString() {
		if len(e.LogString.String) < 1 {
			return fmt.Errorf(`string event is invalid`)
		}
	}

	if e.isGestureLongClick() {
		if e.GestureLongClick.X < 0 || e.GestureLongClick.Y < 0 {
			return fmt.Errorf(`gesture_long_click event is invalid`)
		}
	}

	if e.isGestureScroll() {
		if e.GestureScroll.X < 0 || e.GestureScroll.Y < 0 {
			return fmt.Errorf(`gesture_scroll event is invalid`)
		}
	}

	if e.isGestureClick() {
		if e.GestureClick.X < 0 || e.GestureClick.Y < 0 {
			return fmt.Errorf(`gesture_click event is invalid`)
		}
	}

	if e.isLifecycleActivity() {
		if e.LifecycleActivity.Type == "" || e.LifecycleActivity.ClassName == "" {
			return fmt.Errorf(`lifecycle_activity event is invalid`)
		}
	}

	if e.isLifecycleFragment() {
		if e.LifecycleFragment.Type == "" || e.LifecycleFragment.ClassName == "" {
			return fmt.Errorf(`lifecycle_fragment event is invalid`)
		}
	}

	if e.isLifecycleApp() {
		if e.LifecycleApp.Type == "" {
			return fmt.Errorf(`lifecycle_app event is invalid`)
		}
	}

	if e.isColdLaunch() {
		if e.ColdLaunch.ProcessStartUptime <= 0 && e.ColdLaunch.ContentProviderAttachUptime <= 0 && e.ColdLaunch.ProcessStartRequestedUptime <= 0 {
			return fmt.Errorf(`one of cold_launch.process_start_uptime, cold_launch.process_start_requested_uptime, cold_launch.content_provider_attach_uptime must be greater than 0`)
		}
		if e.ColdLaunch.OnNextDrawUptime <= 0 {
			return fmt.Errorf(`cold_launch.on_next_draw_uptime must be greater than 0`)
		}
		if e.ColdLaunch.LaunchedActivity == "" {
			return fmt.Errorf(`cold_launch.launched_activity must not be empty`)
		}
	}

	if e.isWarmLaunch() {
		if e.WarmLaunch.AppVisibleUptime <= 0 {
			return fmt.Errorf(`warm_launch.app_visible_uptime must be greater than 0`)
		}
		if e.WarmLaunch.OnNextDrawUptime <= 0 {
			return fmt.Errorf(`warm_launch.on_next_draw_uptime must be greater than 0`)
		}
		if e.WarmLaunch.LaunchedActivity == "" {
			return fmt.Errorf(`warm_launch.launched_activity must not be empty`)
		}
	}

	if e.isHotLaunch() {
		if e.HotLaunch.AppVisibleUptime <= 0 {
			return fmt.Errorf(`hot_launch.app_visible_uptime must be greater than 0`)
		}
		if e.HotLaunch.OnNextDrawUptime <= 0 {
			return fmt.Errorf(`hot_launch.on_next_draw_uptime must be greater than 0`)
		}
		if e.HotLaunch.LaunchedActivity == "" {
			return fmt.Errorf(`hot_launch.launched_activity must not be empty`)
		}
	}

	if e.isNetworkChange() {
		if e.NetworkChange.NetworkType == "" {
			return fmt.Errorf(`network_change.network_type must not be empty`)
		}
	}

	if e.isHttp() {
		if e.Http.URL == "" {
			return fmt.Errorf(`http.url must not be empty`)
		}
		if e.Http.Method == "" {
			return fmt.Errorf(`http.method must not be empty`)
		}
	}

	if e.isMemoryUsage() {
		if e.MemoryUsage.IntervalConfig <= 0 {
			return fmt.Errorf(`memory_usage.interval_config must be greater than 0`)
		}
	}

	if e.isTrimMemory() {
		if e.TrimMemory.Level == "" {
			return fmt.Errorf(`trim_memory.level must not be empty`)
		}
	}

	if e.isCPUUsage() {
		if e.CPUUsage.NumCores <= 0 {
			return fmt.Errorf(`cpu_usage.num_cores must be greater than 0`)
		}
		if e.CPUUsage.ClockSpeed <= 0 {
			return fmt.Errorf(`cpu_usage.clock_speed must be greater than 0`)
		}
		if e.CPUUsage.IntervalConfig <= 0 {
			return fmt.Errorf(`cpu_usage.interval_config must be greater than 0`)
		}
	}

	if e.isNavigation() {
		if e.Navigation.Route == "" {
			return fmt.Errorf(`navigation.route must not be empty`)
		}
	}

	if len(e.Type) > maxTypeChars {
		return fmt.Errorf(`"events[].type" exceeds maximum allowed characters of (%d)`, maxTypeChars)
	}
	if len(e.ThreadName) > maxThreadNameChars {
		return fmt.Errorf(`"events[].thread_name" exceeds maximum allowed characters of (%d)`, maxThreadNameChars)
	}
	if len(e.ANR.ThreadName) > maxThreadNameChars {
		return fmt.Errorf(`"events[].anr.thread_name" exceeds maximum allowed characters of (%d)`, maxThreadNameChars)
	}
	if len(e.Exception.ThreadName) > maxThreadNameChars {
		return fmt.Errorf(`"events[].exception.thread_name" exceeds maximum allowed characters of (%d)`, maxThreadNameChars)
	}
	if len(e.AppExit.Reason) > maxAppExitReasonChars {
		return fmt.Errorf(`"events[].app_exit.reason" exceeds maximum allowed characters of (%d)`, maxAppExitReasonChars)
	}
	if len(e.AppExit.Importance) > maxAppExitImportanceChars {
		return fmt.Errorf(`"events[].app_exit.importance exceeds maximum allowed characters of (%d)`, maxAppExitImportanceChars)
	}
	if len(e.LogString.SeverityText) > maxSeverityTextChars {
		return fmt.Errorf(`"events[].string.severity_text" exceeds maximum allowed characters of (%d)`, maxSeverityTextChars)
	}
	if len(e.GestureLongClick.Target) > maxGestureLongClickTargetChars {
		return fmt.Errorf(`"events[].gesture_long_click.target" exceeds maximum allowed characters of (%d)`, maxGestureLongClickTargetChars)
	}
	if len(e.GestureLongClick.TargetID) > maxGestureLongClickTargetIDChars {
		return fmt.Errorf(`"events[].gesture_long_click.target_id" exceeds maximum allowed characters of (%d)`, maxGestureLongClickTargetIDChars)
	}
	if len(e.GestureClick.Target) > maxGestureClickTargetChars {
		return fmt.Errorf(`"events[].gesture_click.target" exceeds maximum allowed characters of (%d)`, maxGestureClickTargetChars)
	}
	if len(e.GestureClick.TargetID) > maxGestureClickTargetIDChars {
		return fmt.Errorf(`"events[].gesture_click.target_id" exceeds maximum allowed characters of (%d)`, maxGestureClickTargetIDChars)
	}
	if len(e.GestureScroll.Target) > maxGestureScrollTargetChars {
		return fmt.Errorf(`"events[].gesture_scroll.target" exceeds maximum allowed characters of (%d)`, maxGestureScrollTargetChars)
	}
	if len(e.GestureScroll.TargetID) > maxGestureScrollTargetIDChars {
		return fmt.Errorf(`"events[].gesture_scroll.target_id" exceeds maximum allowed characters of (%d)`, maxGestureScrollTargetIDChars)
	}
	if len(e.GestureScroll.Direction) > maxGestureScrollDirectionChars {
		return fmt.Errorf(`"events[].gesture_scroll.direction" exceeds maximum allowed characters of (%d)`, maxGestureScrollDirectionChars)
	}
	if len(e.LifecycleActivity.Type) > maxLifecycleActivityTypeChars {
		return fmt.Errorf(`"events[].lifecycle_activity.type" exceeds maximum allowed characters of (%d)`, maxLifecycleActivityTypeChars)
	}
	if len(e.LifecycleActivity.ClassName) > maxLifecycleActivityClassNameChars {
		return fmt.Errorf(`"events[].lifecycle_activity.class_name" exceeds maximum allowed characters of (%d)`, maxLifecycleActivityClassNameChars)
	}
	if len(e.LifecycleFragment.Type) > maxLifecycleFragmentTypeChars {
		return fmt.Errorf(`"events[].lifecycle_fragment.type" exceeds maximum allowed characters of (%d)`, maxLifecycleFragmentTypeChars)
	}
	if len(e.LifecycleFragment.ClassName) > maxLifecycleFragmentClassNameChars {
		return fmt.Errorf(`"events[].lifecycle_fragment.class_name" exceeds maximum allowed characters of (%d)`, maxLifecycleFragmentClassNameChars)
	}
	if len(e.LifecycleApp.Type) > maxLifecycleAppTypeChars {
		return fmt.Errorf(`"events[].lifecycle_app.type" exceeds maximum allowed characters of (%d)`, maxLifecycleAppTypeChars)
	}
	if len(e.ColdLaunch.LaunchedActivity) == maxColdLaunchLaunchedActivityChars {
		return fmt.Errorf(`events[].cold_launch.launched_activity exceeds maximum allowed characters of (%d)`, maxColdLaunchLaunchedActivityChars)
	}
	if len(e.WarmLaunch.LaunchedActivity) == maxWarmLaunchLaunchedActivityChars {
		return fmt.Errorf(`events[].warm_launch.launched_activity exceeds maximum allowed characters of (%d)`, maxWarmLaunchLaunchedActivityChars)
	}
	if len(e.HotLaunch.LaunchedActivity) == maxHotLaunchLaunchedActivityChars {
		return fmt.Errorf(`events[].hot_launch.launched_activity exceeds maximum allowed characters of (%d)`, maxHotLaunchLaunchedActivityChars)
	}
	if len(e.NetworkChange.NetworkType) == maxNetworkChangeNetworkTypeChars {
		return fmt.Errorf(`events[].network_change.network_type exceeds maximum allowed characters of (%d)`, maxNetworkChangeNetworkTypeChars)
	}
	if len(e.NetworkChange.PreviousNetworkType) == maxNetworkChangePreviousNetworkTypeChars {
		return fmt.Errorf(`events[].network_change.previous_network_type exceeds maximum allowed characters of (%d)`, maxNetworkChangePreviousNetworkTypeChars)
	}
	if len(e.NetworkChange.NetworkGeneration) == maxNetworkChangeNetworkGeneration {
		return fmt.Errorf(`events[].network_change.network_generation exceeds maximum allowed characters of (%d)`, maxNetworkChangeNetworkGeneration)
	}
	if len(e.NetworkChange.PreviousNetworkGeneration) == maxNetworkChangePreviousNetworkGeneration {
		return fmt.Errorf(`events[].network_change.previous_network_generation exceeds maximum allowed characters of (%d)`, maxNetworkChangePreviousNetworkGeneration)
	}
	if len(e.NetworkChange.NetworkProvider) == maxNetworkChangeNetworkProvider {
		return fmt.Errorf(`events[].network_change.network_provider exceeds maximum allowed characters of (%d)`, maxNetworkChangeNetworkProvider)
	}
	if len(e.ANR.DeviceLocale) > maxAnrDeviceLocaleChars {
		return fmt.Errorf(`"events[].anr.device_locale" exceeds maximum allowed characters of (%d)`, maxAnrDeviceLocaleChars)
	}
	if len(e.Exception.DeviceLocale) > maxExceptionDeviceLocaleChars {
		return fmt.Errorf(`"events[].exception.device_locale" exceeds maximum allowed characters of (%d)`, maxExceptionDeviceLocaleChars)
	}
	if len(e.Http.Method) > maxHttpMethodChars {
		return fmt.Errorf(`"events[].http.method" exceeds maximum allowed characters of (%d)`, maxHttpMethodChars)
	}
	if len(e.Http.Client) > maxHttpClientChars {
		return fmt.Errorf(`"events[].http.client" exceeds maximum allowed characters of (%d)`, maxHttpClientChars)
	}
	if len(e.TrimMemory.Level) > maxTrimMemoryLevelChars {
		return fmt.Errorf(`"events[].trim_memo̦ry.level" exceeds maximum allowed characters of (%d)`, maxTrimMemoryLevelChars)
	}
	if len(e.Attributes) > maxAttrCount {
		return fmt.Errorf(`"events[].attributes" exceeds maximum count of (%d)`, maxAttrCount)
	}
	if len(e.Navigation.Route) > maxRouteChars {
		return fmt.Errorf(`"events[].navigation.route" exceeds maximum allowed characters of (%d)`, maxRouteChars)
	}

	return nil
}

func (e Exception) getType() string {
	return e.Exceptions[len(e.Exceptions)-1].Type
}

func (e Exception) getMessage() string {
	return e.Exceptions[len(e.Exceptions)-1].Message
}

func (e Exception) getLocation() string {
	frame := e.Exceptions[len(e.Exceptions)-1].Frames[0]
	return frame.String()
}

func (a ANR) getType() string {
	return a.Exceptions[len(a.Exceptions)-1].Type
}

func (a ANR) getMessage() string {
	return a.Exceptions[len(a.Exceptions)-1].Message
}

func (a ANR) getLocation() string {
	frame := a.Exceptions[len(a.Exceptions)-1].Frames[0]
	return frame.String()
}
