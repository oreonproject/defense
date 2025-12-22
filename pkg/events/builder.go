// oreon/defense · watchthelight <wtl>

package events

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

// Builder accumulates event data during an operation.
// Use Start() to create, chain Set() calls, then End() to finalize.
type Builder struct {
	event     Event
	startedAt time.Time
}

// Start creates a new event builder for an operation.
func Start(eventType EventType, component string) *Builder {
	return &Builder{
		event: Event{
			Type:        eventType,
			OperationID: generateID(),
			Component:   component,
			Fields:      make(map[string]interface{}),
		},
		startedAt: time.Now(),
	}
}

// WithOperationID sets a specific operation ID (for correlation with parent).
func (b *Builder) WithOperationID(id string) *Builder {
	b.event.OperationID = id
	return b
}

// Set adds a field to the event.
func (b *Builder) Set(key string, value interface{}) *Builder {
	b.event.Fields[key] = value
	return b
}

// SetError marks the event as failed with an error.
func (b *Builder) SetError(err error) *Builder {
	if err != nil {
		b.event.Success = false
		b.event.Error = err.Error()
	}
	return b
}

// End finalizes the event with timing and returns it.
func (b *Builder) End() Event {
	b.event.StartedAt = b.startedAt
	b.event.Duration = time.Since(b.startedAt)
	b.event.DurationMs = b.event.Duration.Milliseconds()
	if b.event.Error == "" {
		b.event.Success = true
	}
	return b.event
}

// generateID creates a short random ID for operation tracking.
func generateID() string {
	b := make([]byte, 4)
	rand.Read(b)
	return hex.EncodeToString(b)
}
