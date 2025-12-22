// oreon/defense · watchthelight <wtl>

package events

import (
	"context"
	"crypto/rand"
	"encoding/hex"
)

// contextKey is a private type for context keys to avoid collisions.
type contextKey int

const operationIDKey contextKey = iota

// WithOperationID returns a new context with the given operation ID attached.
func WithOperationID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, operationIDKey, id)
}

// OperationIDFromContext extracts the operation ID from the context.
// Returns empty string if no operation ID is present.
func OperationIDFromContext(ctx context.Context) string {
	if id, ok := ctx.Value(operationIDKey).(string); ok {
		return id
	}
	return ""
}

// NewOperationContext creates a new context with a fresh operation ID.
// Returns the context and the generated ID.
func NewOperationContext(ctx context.Context) (context.Context, string) {
	id := newOperationID()
	return WithOperationID(ctx, id), id
}

// newOperationID generates a new random operation ID.
func newOperationID() string {
	b := make([]byte, 4)
	rand.Read(b)
	return hex.EncodeToString(b)
}
