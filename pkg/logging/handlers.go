// oreon/defense · watchthelight <wtl>

package logging

import (
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

// NewFileHandler creates a handler that writes JSON logs to a file.
// Creates parent directories if needed. Appends to existing file.
func NewFileHandler(path string, level slog.Level) (slog.Handler, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, err
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	return slog.NewJSONHandler(file, &slog.HandlerOptions{
		Level: level,
	}), nil
}
