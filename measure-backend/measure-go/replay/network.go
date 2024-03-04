package replay

import (
	"measure-backend/measure-go/event"
	"time"
)

// NetworkChange represents network change events
// suitable for session replay.
type NetworkChange struct {
	EventType string `json:"event_type"`
	*event.NetworkChange
	ThreadName string            `json:"-"`
	Timestamp  time.Time         `json:"timestamp"`
	Attributes map[string]string `json:"attributes"`
}

// GetThreadName provides the name of the thread
// where the network change took place.
func (nc NetworkChange) GetThreadName() string {
	return nc.ThreadName
}

// GetTimestamp provides the timestamp of
// the network change event.
func (nc NetworkChange) GetTimestamp() time.Time {
	return nc.Timestamp
}

// ComputeNetworkChange computes network change
// events for session replay.
func ComputeNetworkChange(events []event.EventField) (result []ThreadGrouper) {
	for _, event := range events {
		event.NetworkChange.Trim()
		netChanges := NetworkChange{
			event.Type,
			&event.NetworkChange,
			event.ThreadName,
			event.Timestamp,
			event.Attributes,
		}
		result = append(result, netChanges)
	}

	return
}
