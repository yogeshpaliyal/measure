package replay

import (
	"measure-backend/measure-go/event"
	"time"
)

// NominalColdLaunchThreshold defines the upper bound
// of a nominal cold launch duration.
var NominalColdLaunchThreshold = 30 * time.Second

// ColdLaunch represents cold launch events
// suitable for session replay.
type ColdLaunch struct {
	EventType  string            `json:"event_type"`
	Duration   time.Duration     `json:"duration"`
	ThreadName string            `json:"-"`
	Timestamp  time.Time         `json:"timestamp"`
	Attributes map[string]string `json:"attributes"`
}

// GetThreadName provides the name of the thread
// where cold launch took place.
func (cl ColdLaunch) GetThreadName() string {
	return cl.ThreadName
}

// GetTimestamp provides the timestamp of
// the cold launch event.
func (cl ColdLaunch) GetTimestamp() time.Time {
	return cl.Timestamp
}

// WarmLaunch represents warm launch events
// suitable for session replay.
type WarmLaunch struct {
	EventType        string            `json:"event_type"`
	Duration         time.Duration     `json:"duration"`
	LaunchedActivity string            `json:"launched_activity"`
	HasSavedState    bool              `json:"has_saved_state"`
	IntentData       string            `json:"intent_data"`
	ThreadName       string            `json:"-"`
	Timestamp        time.Time         `json:"timestamp"`
	Attributes       map[string]string `json:"attributes"`
}

// GetThreadName provides the name of the thread
// where warm launch took place.
func (wl WarmLaunch) GetThreadName() string {
	return wl.ThreadName
}

// GetTimestamp provides the timestamp of
// the warm launch event.
func (wl WarmLaunch) GetTimestamp() time.Time {
	return wl.Timestamp
}

// HotLaunch represents hot launch events
// suitable for session replay.
type HotLaunch struct {
	EventType        string            `json:"event_type"`
	Duration         time.Duration     `json:"duration"`
	LaunchedActivity string            `json:"launched_activity"`
	HasSavedState    bool              `json:"has_saved_state"`
	IntentData       string            `json:"intent_data"`
	ThreadName       string            `json:"-"`
	Timestamp        time.Time         `json:"timestamp"`
	Attributes       map[string]string `json:"attributes"`
}

// GetThreadName provides the name of the thread
// where hot launch took place.
func (hl HotLaunch) GetThreadName() string {
	return hl.ThreadName
}

// GetTimestamp provides the timestamp of
// the hot launch event.
func (hl HotLaunch) GetTimestamp() time.Time {
	return hl.Timestamp
}

// ComputeColdLaunches computes cold launch events
// for session replay.
func ComputeColdLaunches(events []event.EventField) (result []ThreadGrouper) {
	for _, event := range events {
		event.ColdLaunch.Trim()

		coldLaunches := ColdLaunch{
			event.Type,
			event.ColdLaunch.Duration,
			event.ThreadName,
			event.Timestamp,
			event.Attributes,
		}
		result = append(result, coldLaunches)
	}

	return
}

// ComputeWarmLaunches computes warm launch events
// for session replay.
func ComputeWarmLaunches(events []event.EventField) (result []ThreadGrouper) {
	for _, event := range events {
		warmLaunches := WarmLaunch{
			event.Type,
			event.WarmLaunch.Duration,
			event.WarmLaunch.LaunchedActivity,
			event.WarmLaunch.HasSavedState,
			event.WarmLaunch.IntentData,
			event.ThreadName,
			event.Timestamp,
			event.Attributes,
		}
		result = append(result, warmLaunches)
	}

	return
}

// ComputeHotLaunches computes hot launch events
// for session replay.
func ComputeHotLaunches(events []event.EventField) (result []ThreadGrouper) {
	for _, event := range events {
		hotLaunches := HotLaunch{
			event.Type,
			event.HotLaunch.Duration,
			event.HotLaunch.LaunchedActivity,
			event.HotLaunch.HasSavedState,
			event.HotLaunch.IntentData,
			event.ThreadName,
			event.Timestamp,
			event.Attributes,
		}
		result = append(result, hotLaunches)
	}

	return
}
