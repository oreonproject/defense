// oreon/defense · watchthelight <wtl>

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

// New creates a configured logger.
// Sets up file and/or journald handlers based on config.
// Returns a cleanup function that should be called on shutdown.
func New(cfg Config) (*slog.Logger, func() error, error) {
	level := ParseLevel(cfg.Level)
	var handlers []slog.Handler
	var closers []func() error

	if cfg.FilePath != "" {
		h, err := NewFileHandler(cfg.FilePath, level)
		if err != nil {
			return nil, nil, err
		}
		handlers = append(handlers, h)
		closers = append(closers, h.Close)
	}

	if cfg.UseJournald {
		if h := NewJournaldHandler(level); h != nil {
			handlers = append(handlers, h)
		}
	}

	cleanup := func() error {
		for _, fn := range closers {
			if err := fn(); err != nil {
				return err
			}
		}
		return nil
	}

	if len(handlers) == 0 {
		return slog.Default(), cleanup, nil
	}

	return slog.New(NewMultiHandler(handlers...)), cleanup, nil
}

// WithComponent returns a child logger with "component" field set.
func WithComponent(logger *slog.Logger, component string) *slog.Logger {
	return logger.With("component", component)
}
