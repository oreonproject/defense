package logging

import "log/slog"

// Config holds logging configuration.
type Config struct {
	Level       string // "debug", "info", "warn", "error"
	FilePath    string // path to log file, empty = no file logging
	UseJournald bool   // write to journald when running under systemd
}

// ParseLevel converts a string level name to slog.Level.
// Defaults to LevelInfo for unknown values.
func ParseLevel(s string) slog.Level {
	switch s {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
