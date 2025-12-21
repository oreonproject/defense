// oreon/defense · watchthelight <wtl>

package daemon

import (
	"context"
	"log/slog"
	"time"

	"github.com/oreonproject/defense/pkg/config"
)

// Daemon is the main defense daemon that coordinates scanning,
// firewall, and protection state.
type Daemon struct {
	cfg    *config.Config
	state  *StateManager
	logger *slog.Logger
}

// New creates a new daemon instance.
func New(cfg *config.Config, logger *slog.Logger) *Daemon {
	return &Daemon{
		cfg:    cfg,
		state:  NewStateManager(),
		logger: logger,
	}
}

// State returns the state manager for external access (e.g. IPC).
func (d *Daemon) State() *StateManager {
	return d.state
}

// Run starts the daemon and blocks until context is cancelled.
func (d *Daemon) Run(ctx context.Context, socketPath string) error {
	d.logger.Info("daemon starting")

	// Start IPC server
	server := NewServer(socketPath, d.state)
	if err := server.Listen(); err != nil {
		return err
	}
	go server.Serve()
	defer server.Close()

	// initial health check
	d.healthCheck()

	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			d.logger.Info("daemon shutting down")
			return nil
		case <-ticker.C:
			d.healthCheck()
		}
	}
}

// healthCheck evaluates system state and updates the state machine.
// TODO: actually check clamav, firewall, rules
func (d *Daemon) healthCheck() {
	d.logger.Debug("running health check")

	// placeholder - just set to protected for now
	// real implementation will check:
	// - clamav available?
	// - firewall enabled (if configured)?
	// - rules up to date?

	if d.state.State() == StateStarting {
		d.state.SetState(StateProtected)
		d.logger.Info("daemon ready", "state", d.state.State())
	}
}
