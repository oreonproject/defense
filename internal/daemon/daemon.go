// oreon/defense · watchthelight <wtl>

package daemon

import (
	"context"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/oreonproject/defense/internal/scanner"
	"github.com/oreonproject/defense/pkg/config"
	"github.com/oreonproject/defense/pkg/events"
)

// Daemon is the main defense daemon that coordinates scanning,
// firewall, and protection state.
type Daemon struct {
	cfg     *config.Config
	state   *StateManager
	logger  *slog.Logger
	scanner *scanner.ClamAV
	events  *events.Emitter

	// Runtime state (may differ from config)
	firewallEnabled bool
	lastScan        time.Time
	rulesUpdated    time.Time
}

// New creates a new daemon instance.
func New(cfg *config.Config, logger *slog.Logger) *Daemon {
	return &Daemon{
		cfg:             cfg,
		state:           NewStateManager(),
		logger:          logger,
		scanner:         scanner.New(cfg.ClamAV.SocketPath),
		events:          events.NewEmitter(events.WithLogger(logger)),
		firewallEnabled: cfg.Firewall.Enabled,
		rulesUpdated:    time.Now(), // Assume rules are current at startup
	}
}

// Config returns the daemon configuration.
func (d *Daemon) Config() *config.Config {
	return d.cfg
}

// FirewallEnabled returns whether the firewall is currently enabled.
func (d *Daemon) FirewallEnabled() bool {
	return d.firewallEnabled
}

// SetFirewallEnabled enables or disables the firewall.
func (d *Daemon) SetFirewallEnabled(enabled bool) {
	d.firewallEnabled = enabled
	d.cfg.Firewall.Enabled = enabled
	d.logger.Info("firewall toggled", "enabled", enabled)
}

// LastScan returns the time of the last scan.
func (d *Daemon) LastScan() time.Time {
	return d.lastScan
}

// SetLastScan updates the last scan time.
func (d *Daemon) SetLastScan(t time.Time) {
	d.lastScan = t
}

// RulesUpdated returns the time rules were last updated.
func (d *Daemon) RulesUpdated() time.Time {
	return d.rulesUpdated
}

// State returns the state manager for external access (e.g. IPC).
func (d *Daemon) State() *StateManager {
	return d.state
}

// Scanner returns the ClamAV scanner instance.
func (d *Daemon) Scanner() *scanner.ClamAV {
	return d.scanner
}

// Events returns the event emitter for logging wide events.
func (d *Daemon) Events() *events.Emitter {
	return d.events
}

// Run starts the daemon and blocks until context is cancelled.
func (d *Daemon) Run(ctx context.Context, socketPath string) error {
	d.logger.Info("daemon starting")

	// Start IPC server
	server := NewServer(socketPath, d)
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
func (d *Daemon) healthCheck() {
	d.logger.Debug("running health check")

	// Don't change state if we're scanning or paused
	currentState := d.state.State()
	if currentState == StateScanning || currentState == StatePaused {
		return
	}

	// Check ClamAV availability
	clamAvailable := d.checkClamAV()

	// Determine the appropriate state
	var newState State

	if !clamAvailable {
		newState = StateWarning
		d.logger.Warn("ClamAV not available")
	} else if !d.firewallEnabled && d.cfg.Firewall.Enabled {
		// Firewall should be on but isn't
		newState = StateWarning
		d.logger.Warn("firewall disabled but should be enabled")
	} else {
		newState = StateProtected
	}

	if currentState != newState {
		d.state.SetState(newState)
		d.logger.Info("state changed", "from", currentState, "to", newState)
	}

	if currentState == StateStarting {
		d.logger.Info("daemon ready", "state", d.state.State())
	}
}

// checkClamAV verifies ClamAV daemon is available.
func (d *Daemon) checkClamAV() bool {
	// Use scanner to check (does ping)
	if d.scanner.IsAvailable() {
		d.logger.Debug("ClamAV daemon available via configured socket")
		return true
	}

	// Fallback: check common socket paths
	fallbackPaths := []string{
		"/var/run/clamav/clamd.sock",
		"/var/run/clamav/clamd.ctl",
		"/run/clamav/clamd.sock",
		"/tmp/clamd.socket",
	}

	for _, path := range fallbackPaths {
		if _, err := os.Stat(path); err == nil {
			d.logger.Debug("found ClamAV socket", "path", path)
			return true
		}
	}

	// Also check if clamd process is running via /proc
	if _, err := os.Stat("/proc"); err == nil {
		entries, err := os.ReadDir("/proc")
		if err == nil {
			for _, entry := range entries {
				if !entry.IsDir() {
					continue
				}
				cmdlinePath := "/proc/" + entry.Name() + "/cmdline"
				data, err := os.ReadFile(cmdlinePath)
				if err != nil {
					continue
				}
				if strings.Contains(string(data), "clamd") {
					d.logger.Debug("found ClamAV process")
					return true
				}
			}
		}
	}

	return false
}
