package logging

import (
	"log/slog"
	"os"
	"path/filepath"
)

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
