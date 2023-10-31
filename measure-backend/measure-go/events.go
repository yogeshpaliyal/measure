package main

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
)

// maximum character limits for event fields
const (
	maxTypeChars                       = 32
	maxAppExitReasonChars              = 64
	maxAppExitImportanceChars          = 32
	maxSeverityTextChars               = 10
	maxGestureLongClickTargetChars     = 128
	maxGestureLongClickTargetNameChars = 128
	maxGestureLongClickTargetIDChars   = 128
	maxGestureScrollTargetChars        = 128
	maxGestureScrollTargetNameChars    = 128
	maxGestureScrollTargetIDChars      = 128
	maxGestureScrollDirectionChars     = 8
	maxGestureClickTargetChars         = 128
	maxGestureClickTargetNameChars     = 128
	maxGestureClickTargetIDChars       = 128
	maxHTTPRequestMethodChars          = 16
	maxHTTPRequestProtocolVersionChars = 16
	maxHTTPResponseMethodChars         = 16
	maxLifecycleActivityTypeChars      = 32
	maxLifecycleActivityClassNameChars = 128
	maxLifecycleFragmentTypeChars      = 32
	maxLifecycleFragmentClassNameChars = 128
	maxLifecycleAppTypeChars           = 32
	maxAttrCount                       = 10
)

const TypeANR = "anr"
const TypeException = "exception"
const TypeAppExit = "app_exit"
const TypeString = "string"
const TypeGestureLongClick = "gesture_long_click"
const TypeGestureClick = "gesture_click"
const TypeGestureScroll = "gesture_scroll"
const TypeHTTPRequest = "http_request"
const TypeHTTPResponse = "http_response"
const TypeLifecycleActivity = "lifecycle_activity"
const TypeLifecycleFragment = "lifecycle_fragment"
const TypeLifecycleApp = "lifecycle_app"

// timeFormat is the format of datetime in nanoseconds when
// converting datetime values before inserting into database
const timeFormat = "2006-01-02 15:04:05.999999999"

// list of events table columns
var columns = []string{
	"id",
	"type",
	"session_id",
	"timestamp",
	"resource.device_name",
	"resource.device_model",
	"resource.device_manufacturer",
	"resource.device_type",
	"resource.device_is_foldable",
	"resource.device_is_physical",
	"resource.device_density_dpi",
	"resource.device_width_px",
	"resource.device_height_px",
	"resource.device_density",
	"resource.os_name",
	"resource.os_version",
	"resource.platform",
	"resource.app_version",
	"resource.app_build",
	"resource.app_unique_id",
	"resource.measure_sdk_version",
	"anr.thread_name",
	"anr.handled",
	"anr_exceptions",
	"anr_threads",
	"exception.thread_name",
	"exception.handled",
	"exception_exceptions",
	"exception_threads",
	"app_exit.reason",
	"app_exit.importance",
	"app_exit.trace",
	"app_exit.process_name",
	"app_exit.pid",
	"app_exit.timestamp",
	"string.severity_text",
	"string.string",
	"gesture_long_click.target",
	"gesture_long_click.target_id",
	"gesture_long_click.touch_down_time",
	"gesture_long_click.touch_up_time",
	"gesture_long_click.width",
	"gesture_long_click.height",
	"gesture_long_click.x",
	"gesture_long_click.y",
	"gesture_click.target",
	"gesture_click.target_id",
	"gesture_click.touch_down_time",
	"gesture_click.touch_up_time",
	"gesture_click.width",
	"gesture_click.height",
	"gesture_click.x",
	"gesture_click.y",
	"gesture_scroll.target",
	"gesture_scroll.target_id",
	"gesture_scroll.touch_down_time",
	"gesture_scroll.touch_up_time",
	"gesture_scroll.x",
	"gesture_scroll.y",
	"gesture_scroll.end_x",
	"gesture_scroll.end_y",
	"gesture_scroll.direction",
	"http_request.request_id",
	"http_request.request_url",
	"http_request.method",
	"http_request.http_protocol_version",
	"http_request.request_body_size",
	"http_request.request_body",
	"http_request.request_headers",
	"http_response.request_id",
	"http_response.request_url",
	"http_response.method",
	"http_response.latency_ms",
	"http_response.status_code",
	"http_response.response_body",
	"http_response.response_headers",
	"lifecycle_activity.type",
	"lifecycle_activity.class_name",
	"lifecycle_activity.intent",
	"lifecycle_activity.saved_instance_state",
	"lifecycle_fragment.type",
	"lifecycle_fragment.class_name",
	"lifecycle_fragment.parent_activity",
	"lifecycle_fragment.tag",
	"lifecycle_app.type",
	"attributes",
}

type Frame struct {
	LineNum    int    `json:"line_num"`
	ColNum     int    `json:"col_num"`
	ModuleName string `json:"module_name"`
	FileName   string `json:"file_name"`
	ClassName  string `json:"class_name"`
	MethodName string `json:"method_name"`
}

func (f *Frame) encode() string {
	moduleName := f.ModuleName
	fileName := f.FileName
	className := f.ClassName
	methodName := f.MethodName
	if f.ModuleName == "" {
		moduleName = "__blank__"
	}
	if f.FileName == "" {
		fileName = "__blank__"
	}
	if f.ClassName == "" {
		className = "__blank__"
	}
	if f.MethodName == "" {
		methodName = "__blank__"
	}
	return fmt.Sprintf("(%d, %d, '%s', '%s', '%s', '%s')", f.LineNum, f.ColNum, moduleName, fileName, className, methodName)
}

type Frames []Frame

func (frames Frames) encode() string {
	var collection []string
	for _, frame := range frames {
		collection = append(collection, frame.encode())
	}
	return fmt.Sprintf("[%s]", strings.Join(collection, ", "))
}

type ExceptionUnit struct {
	Type    string `json:"type" binding:"required"`
	Message string `json:"message"`
	Frames  Frames `json:"frames" binding:"required"`
}

func (eu *ExceptionUnit) encode() string {
	return fmt.Sprintf("('%s', '%s', %s)", eu.Type, eu.Message, eu.Frames.encode())
}

type ExceptionUnits []ExceptionUnit

func (exceptionUnits ExceptionUnits) encode() string {
	var units []string
	for _, exceptionUnit := range exceptionUnits {
		units = append(units, exceptionUnit.encode())
	}

	return fmt.Sprintf("[%s]", strings.Join(units, ", "))
}

type Thread struct {
	Name   string `json:"name" binding:"required"`
	Frames Frames `json:"frames" binding:"required"`
}

func (thread *Thread) encode() string {
	return fmt.Sprintf("('%s', %s)", thread.Name, thread.Frames.encode())
}

type Threads []Thread

func (threads Threads) encode() string {
	var collection []string
	for _, thread := range threads {
		collection = append(collection, thread.encode())
	}

	return fmt.Sprintf("[%s]", strings.Join(collection, ", "))
}

type ANR struct {
	ThreadName string         `json:"thread_name" binding:"required"`
	Handled    bool           `json:"handled" binding:"required"`
	Exceptions ExceptionUnits `json:"exceptions" binding:"required"`
	Threads    Threads        `json:"threads" binding:"required"`
}

type Exception struct {
	ThreadName string         `json:"thread_name" binding:"required"`
	Handled    bool           `json:"handled" binding:"required"`
	Exceptions ExceptionUnits `json:"exceptions" binding:"required"`
	Threads    Threads        `json:"threads" binding:"required"`
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
	Target        string    `json:"target"`
	TargetID      string    `json:"target_id"`
	TouchDownTime uint32    `json:"touch_down_time"`
	TouchUpTime   uint32    `json:"touch_up_time"`
	Width         uint16    `json:"width"`
	Height        uint16    `json:"height"`
	X             float32   `json:"x"`
	Y             float32   `json:"y"`
}

type GestureScroll struct {
	Target        string    `json:"target"`
	TargetID      string    `json:"target_id"`
	TouchDownTime uint32    `json:"touch_down_time"`
	TouchUpTime   uint32    `json:"touch_up_time"`
	X             float32   `json:"x"`
	Y             float32   `json:"y"`
	EndX          float32   `json:"end_x"`
	EndY          float32   `json:"end_y"`
	Direction     string    `json:"direction"`
}

type GestureClick struct {
	Target        string    `json:"target"`
	TargetID      string    `json:"target_id"`
	TouchDownTime uint32    `json:"touch_down_time"`
	TouchUpTime   uint32    `json:"touch_up_time"`
	Width         uint16    `json:"width"`
	Height        uint16    `json:"height"`
	X             float32   `json:"x"`
	Y             float32   `json:"y"`
}

type HTTPRequest struct {
	RequestID           string            `json:"request_id"`
	RequestURL          string            `json:"request_url"`
	Method              string            `json:"method"`
	HTTPProtocolVersion string            `json:"http_protocol_version"`
	RequestBodySize     uint32            `json:"request_body_size"`
	RequestBody         string            `json:"request_body"`
	RequestHeaders      map[string]string `json:"request_headers"`
}

type HTTPResponse struct {
	RequestID       string            `json:"request_id"`
	RequestURL      string            `json:"request_url"`
	Method          string            `json:"method"`
	LatencyMS       uint16            `json:"latency_ms"`
	StatusCode      uint16            `json:"status_code"`
	ResponseBody    string            `json:"response_body"`
	ResponseHeaders map[string]string `json:"response_headers"`
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

type EventField struct {
	Timestamp         time.Time         `json:"timestamp" binding:"required"`
	Type              string            `json:"type" binding:"required"`
	ANR               ANR               `json:"anr,omitempty"`
	Exception         Exception         `json:"exception,omitempty"`
	AppExit           AppExit           `json:"app_exit,omitempty"`
	LogString         LogString         `json:"string,omitempty"`
	GestureLongClick  GestureLongClick  `json:"gesture_long_click,omitempty"`
	GestureScroll     GestureScroll     `json:"gesture_scroll,omitempty"`
	GestureClick      GestureClick      `json:"gesture_click,omitempty"`
	HTTPRequest       HTTPRequest       `json:"http_request,omitempty"`
	HTTPResponse      HTTPResponse      `json:"http_response,omitempty"`
	LifecycleActivity LifecycleActivity `json:"lifecycle_activity,omitempty"`
	LifecycleFragment LifecycleFragment `json:"lifecycle_fragment,omitempty"`
	LifecycleApp      LifecycleApp      `json:"lifecycle_app,omitempty"`
	Attributes        map[string]string `json:"attributes"`
}

func (e *EventField) isException() bool {
	return e.Type == TypeException
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

func (e *EventField) isHTTPRequest() bool {
	return e.Type == TypeHTTPRequest
}

func (e *EventField) isHTTPResponse() bool {
	return e.Type == TypeHTTPResponse
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

func (e *EventField) validate() error {
	validTypes := []string{TypeANR, TypeException, TypeAppExit, TypeString, TypeGestureLongClick, TypeGestureScroll, TypeGestureClick, TypeHTTPRequest, TypeHTTPResponse, TypeLifecycleActivity, TypeLifecycleFragment, TypeLifecycleApp}
	if !slices.Contains(validTypes, e.Type) {
		return fmt.Errorf(`"events[].type" is not a valid type`)
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

	if e.isHTTPRequest() {
		if e.HTTPRequest.RequestID == "" || e.HTTPRequest.RequestURL == "" || e.HTTPRequest.Method == "" {
			return fmt.Errorf(`http_request event is invalid`)
		}
	}

	if e.isHTTPResponse() {
		if e.HTTPResponse.RequestID == "" || e.HTTPResponse.RequestURL == "" || e.HTTPResponse.Method == "" {
			return fmt.Errorf(`http_response event is invalid`)
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

	if len(e.Type) > maxTypeChars {
		return fmt.Errorf(`"events[].type" exceeds maximum allowed characters of (%d)`, maxTypeChars)
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
	if len(e.HTTPRequest.Method) > maxHTTPRequestMethodChars {
		return fmt.Errorf(`"events[].http_request.method" exceeds maximum allowed characters of (%d)`, maxHTTPRequestMethodChars)
	}
	if len(e.HTTPRequest.HTTPProtocolVersion) > maxHTTPRequestProtocolVersionChars {
		return fmt.Errorf(`"events[].http_request.http_protocol_version" exceeds maximum allowed characters of (%d)`, maxHTTPRequestProtocolVersionChars)
	}
	if len(e.HTTPResponse.Method) > maxHTTPResponseMethodChars {
		return fmt.Errorf(`"events[].http_response.method" exceeds maximum allowed characters of (%d)`, maxHTTPResponseMethodChars)
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

	if len(e.Attributes) > maxAttrCount {
		return fmt.Errorf(`"events[].attributes" exceeds maximum count of (%d)`, maxAttrCount)
	}

	return nil
}

func makeInsertQuery(table string, columns []string, session *Session) (string, []interface{}) {
	values := []string{}
	valueArgs := []interface{}{}

	placeholder := "(toUUID(?),?,toUUID(?),?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,toUUID(?),?,?,?,?,?,?,toUUID(?),?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"

	for _, event := range session.Events {
		anrExceptions := "[]"
		anrThreads := "[]"
		exceptionExceptions := "[]"
		exceptionThreads := "[]"
		if event.isANR() {
			anrExceptions = event.ANR.Exceptions.encode()
			anrThreads = event.ANR.Threads.encode()
		}
		if event.isException() {
			exceptionExceptions = event.Exception.Exceptions.encode()
			exceptionThreads = event.Exception.Threads.encode()
		}
		values = append(values, placeholder)
		valueArgs = append(valueArgs,
			uuid.New(),
			event.Type,
			session.SessionID,
			event.Timestamp.Format(timeFormat),
			session.Resource.DeviceName,
			session.Resource.DeviceModel,
			session.Resource.DeviceManufacturer,
			session.Resource.DeviceType,
			session.Resource.DeviceIsFoldable,
			session.Resource.DeviceIsPhysical,
			session.Resource.DeviceDensityDPI,
			session.Resource.DeviceWidthPX,
			session.Resource.DeviceHeightPX,
			session.Resource.DeviceDensity,
			session.Resource.OSName,
			session.Resource.OSVersion,
			session.Resource.Platform,
			session.Resource.AppVersion,
			session.Resource.AppBuild,
			session.Resource.AppUniqueID,
			session.Resource.MeasureSDKVersion,
			event.ANR.ThreadName,
			event.ANR.Handled,
			anrExceptions,
			anrThreads,
			event.Exception.ThreadName,
			event.Exception.Handled,
			exceptionExceptions,
			exceptionThreads,
			event.AppExit.Reason,
			event.AppExit.Importance,
			event.AppExit.Trace,
			event.AppExit.ProcessName,
			event.AppExit.PID,
			event.AppExit.Timestamp,
			event.LogString.SeverityText,
			event.LogString.String,
			event.GestureLongClick.Target,
			event.GestureLongClick.TargetID,
			event.GestureLongClick.TouchDownTime,
			event.GestureLongClick.TouchUpTime,
			event.GestureLongClick.Width,
			event.GestureLongClick.Height,
			event.GestureLongClick.X,
			event.GestureLongClick.Y,
			event.GestureClick.Target,
			event.GestureClick.TargetID,
			event.GestureClick.TouchDownTime,
			event.GestureClick.TouchUpTime,
			event.GestureClick.Width,
			event.GestureClick.Height,
			event.GestureClick.X,
			event.GestureClick.Y,
			event.GestureScroll.Target,
			event.GestureScroll.TargetID,
			event.GestureScroll.TouchDownTime,
			event.GestureScroll.TouchUpTime,
			event.GestureScroll.X,
			event.GestureScroll.Y,
			event.GestureScroll.EndX,
			event.GestureScroll.EndY,
			event.GestureScroll.Direction,
			event.HTTPRequest.RequestID,
			event.HTTPRequest.RequestURL,
			event.HTTPRequest.Method,
			event.HTTPRequest.HTTPProtocolVersion,
			event.HTTPRequest.RequestBodySize,
			event.HTTPRequest.RequestBody,
			mapToString(event.HTTPRequest.RequestHeaders),
			event.HTTPResponse.RequestID,
			event.HTTPResponse.RequestURL,
			event.HTTPResponse.Method,
			event.HTTPResponse.LatencyMS,
			event.HTTPResponse.StatusCode,
			event.HTTPResponse.ResponseBody,
			mapToString(event.HTTPResponse.ResponseHeaders),
			event.LifecycleActivity.Type,
			event.LifecycleActivity.ClassName,
			event.LifecycleActivity.Intent,
			event.LifecycleActivity.SavedInstanceState,
			event.LifecycleFragment.Type,
			event.LifecycleFragment.ClassName,
			event.LifecycleFragment.ParentActivity,
			event.LifecycleFragment.Tag,
			event.LifecycleApp.Type,
			mapToString(event.Attributes),
		)
	}

	query := fmt.Sprintf("insert into %s (%s) values %s;", table, strings.Join(columns, ","), strings.Join(values, ", "))

	return query, valueArgs
}
