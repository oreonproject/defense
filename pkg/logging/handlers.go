// oreon/defense · watchthelight <wtl>

package logging

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
)

// IsUnderSystemd checks if we're running as a systemd service.
func IsUnderSystemd() bool {
	_, ok := os.LookupEnv("JOURNAL_STREAM")
	return ok
}

// NewJournaldHandler creates a handler for systemd journal.
// Returns nil if not running under systemd.
func NewJournaldHandler(level slog.Level) slog.Handler {
	if !IsUnderSystemd() {
		return nil
	}

	return slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: level,
	})
}

// FileHandler wraps a slog.Handler and the underlying file for cleanup.
type FileHandler struct {
	slog.Handler
	file *os.File
}

// Close closes the underlying log file.
func (f *FileHandler) Close() error {
	return f.file.Close()
}

// NewFileHandler creates a handler that writes JSON logs to a file.
// Creates parent directories if needed. Appends to existing file.
// Call Close() when done to release the file handle.
func NewFileHandler(path string, level slog.Level) (*FileHandler, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, err
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	handler := slog.NewJSONHandler(file, &slog.HandlerOptions{
		Level: level,
	})

	return &FileHandler{Handler: handler, file: file}, nil
}

// MultiHandler fans out log records to multiple handlers.
type MultiHandler struct {
	handlers []slog.Handler
}

// NewMultiHandler creates a handler that writes to all given handlers.
// Nil handlers are filtered out.
func NewMultiHandler(handlers ...slog.Handler) *MultiHandler {
	var filtered []slog.Handler
	for _, h := range handlers {
		if h != nil {
			filtered = append(filtered, h)
		}
	}
	return &MultiHandler{handlers: filtered}
}

func (m *MultiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, h := range m.handlers {
		if h.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (m *MultiHandler) Handle(ctx context.Context, r slog.Record) error {
	for _, h := range m.handlers {
		if h.Enabled(ctx, r.Level) {
			if err := h.Handle(ctx, r.Clone()); err != nil {
				return err
			}
		}
	}
	return nil
}

func (m *MultiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	handlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		handlers[i] = h.WithAttrs(attrs)
	}
	return &MultiHandler{handlers: handlers}
}

func (m *MultiHandler) WithGroup(name string) slog.Handler {
	handlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		handlers[i] = h.WithGroup(name)
	}
	return &MultiHandler{handlers: handlers}
}
