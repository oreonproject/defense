// oreon/defense · watchthelight <wtl>

package logging

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseLevel(t *testing.T) {
	tests := []struct {
		input string
		want  slog.Level
	}{
		{"debug", slog.LevelDebug},
		{"info", slog.LevelInfo},
		{"warn", slog.LevelWarn},
		{"error", slog.LevelError},
		{"unknown", slog.LevelInfo},
		{"", slog.LevelInfo},
		{"DEBUG", slog.LevelInfo}, // case sensitive, defaults to info
	}

	for _, tt := range tests {
		got := ParseLevel(tt.input)
		if got != tt.want {
			t.Errorf("ParseLevel(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestFileHandler(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.log")

	h, err := NewFileHandler(path, slog.LevelInfo)
	if err != nil {
		t.Fatalf("NewFileHandler: %v", err)
	}
	defer h.Close()

	logger := slog.New(h)
	logger.Info("test message", "key", "value")

	h.Close() // close before reading

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading log file: %v", err)
	}

	// should be JSON
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("log output not valid JSON: %v\noutput: %s", err, data)
	}

	if m["msg"] != "test message" {
		t.Errorf("msg = %v, want 'test message'", m["msg"])
	}
	if m["key"] != "value" {
		t.Errorf("key = %v, want 'value'", m["key"])
	}
}

func TestFileHandlerLevelFiltering(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.log")

	h, err := NewFileHandler(path, slog.LevelWarn)
	if err != nil {
		t.Fatalf("NewFileHandler: %v", err)
	}

	logger := slog.New(h)
	logger.Info("should not appear")
	logger.Warn("should appear")

	h.Close()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading log file: %v", err)
	}

	if strings.Contains(string(data), "should not appear") {
		t.Error("info message should have been filtered")
	}
	if !strings.Contains(string(data), "should appear") {
		t.Error("warn message should have appeared")
	}
}

func TestFileHandlerCreatesDir(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nested", "dir", "test.log")

	h, err := NewFileHandler(path, slog.LevelInfo)
	if err != nil {
		t.Fatalf("NewFileHandler: %v", err)
	}
	defer h.Close()

	if _, err := os.Stat(filepath.Dir(path)); err != nil {
		t.Errorf("directory not created: %v", err)
	}
}

func TestMultiHandler(t *testing.T) {
	var buf1, buf2 bytes.Buffer

	h1 := slog.NewJSONHandler(&buf1, &slog.HandlerOptions{Level: slog.LevelInfo})
	h2 := slog.NewJSONHandler(&buf2, &slog.HandlerOptions{Level: slog.LevelInfo})

	multi := NewMultiHandler(h1, h2)
	logger := slog.New(multi)
	logger.Info("test message")

	if !strings.Contains(buf1.String(), "test message") {
		t.Error("handler 1 did not receive message")
	}
	if !strings.Contains(buf2.String(), "test message") {
		t.Error("handler 2 did not receive message")
	}
}

func TestMultiHandlerFiltersNil(t *testing.T) {
	var buf bytes.Buffer
	h := slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo})

	multi := NewMultiHandler(nil, h, nil)
	logger := slog.New(multi)
	logger.Info("test message")

	if !strings.Contains(buf.String(), "test message") {
		t.Error("handler did not receive message")
	}
}

func TestMultiHandlerWithAttrs(t *testing.T) {
	var buf bytes.Buffer
	h := slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo})

	multi := NewMultiHandler(h)
	withAttrs := multi.WithAttrs([]slog.Attr{slog.String("added", "attr")})

	logger := slog.New(withAttrs)
	logger.Info("test")

	if !strings.Contains(buf.String(), "added") {
		t.Error("WithAttrs did not propagate")
	}
}

func TestMultiHandlerWithGroup(t *testing.T) {
	var buf bytes.Buffer
	h := slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo})

	multi := NewMultiHandler(h)
	withGroup := multi.WithGroup("mygroup")

	logger := slog.New(withGroup)
	logger.Info("test", "key", "value")

	if !strings.Contains(buf.String(), "mygroup") {
		t.Error("WithGroup did not propagate")
	}
}

func TestMultiHandlerEnabled(t *testing.T) {
	h := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelWarn})
	multi := NewMultiHandler(h)

	if multi.Enabled(context.Background(), slog.LevelInfo) {
		t.Error("should not be enabled for info when handler is warn level")
	}
	if !multi.Enabled(context.Background(), slog.LevelError) {
		t.Error("should be enabled for error")
	}
}

func TestWithComponent(t *testing.T) {
	var buf bytes.Buffer
	h := slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo})
	logger := slog.New(h)

	componentLogger := WithComponent(logger, "scanner")
	componentLogger.Info("test message")

	var m map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &m); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if m["component"] != "scanner" {
		t.Errorf("component = %v, want 'scanner'", m["component"])
	}
}

func TestNew(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.log")

	logger, cleanup, err := New(Config{
		Level:    "debug",
		FilePath: path,
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	logger.Info("test message")

	cleanup()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading log: %v", err)
	}

	if !strings.Contains(string(data), "test message") {
		t.Error("message not written to file")
	}
}

func TestNewNoHandlers(t *testing.T) {
	logger, cleanup, err := New(Config{
		Level: "info",
		// no file, no journald
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer cleanup()

	// should return default logger, not nil
	if logger == nil {
		t.Error("expected non-nil logger")
	}
}

func TestNewInvalidPath(t *testing.T) {
	// try to write to a path we can't create
	// use a null device path that should fail on any OS
	_, _, err := New(Config{
		Level:    "info",
		FilePath: string([]byte{0}) + "/test.log", // null byte in path
	})
	if err == nil {
		t.Error("expected error for invalid path")
	}
}
