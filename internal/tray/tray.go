// Package tray implements the system tray interface for Oreon Defense
package tray

import (
	"log/slog"
	"sync"
	"time"

	"github.com/energye/systray"
	"github.com/esiqveland/notify"
	"github.com/godbus/dbus/v5"
	"github.com/oreonproject/defense/pkg/ipc"
)

// Tray represents the system tray application
type Tray struct {
	client   ipc.Client
	menu     *menu
	notifier notify.Notifier

	iconProtected []byte
	iconWarning   []byte
	iconAlert     []byte
	iconScanning  []byte
	iconPaused    []byte

	mu           sync.Mutex
	currentState string
}

// New creates a new Tray instance
func New(client ipc.Client) *Tray {
	return &Tray{client: client}
}

// Run starts the system tray application
func (t *Tray) Run() error {
	systray.Run(t.onReady, t.onExit)

	return nil
}

// onReady is called when the system tray is ready
func (t *Tray) onReady() {
	// Initialize D-Bus notifier
	conn, err := dbus.SessionBus()
	if err != nil {
		slog.Error("failed to connect to session bus", "error", err)
	} else {
		t.notifier, err = notify.New(conn)
		if err != nil {
			slog.Error("failed to create notifier", "error", err)
		}
	}

	// Load icons
	t.loadIcons()

	// Set initial state
	t.setIcon("protected")

	// Build the menu
	t.menu = newMenu(t)
	t.menu.build()

	// Start status monitoring
	go t.monitorStatus()

	// Show initial notification
	t.showNotification(NotificationStateChange, "Oreon Defense", "Protection is now active")
}

// onExit is called when the system tray is exiting
func (t *Tray) onExit() {
	// Cleanup resources if needed
}

// setIcon updates the tray icon based on the current state
func (t *Tray) setIcon(state string) {
	t.mu.Lock()
	oldState := t.currentState

	if oldState == state {
		t.mu.Unlock()
		return
	}

	t.currentState = state
	t.mu.Unlock()

	// Update icon and tooltip
	switch state {
	case "protected":
		systray.SetIcon(t.iconProtected)
		systray.SetTooltip("Oreon Defense - Protected")
	case "warning":
		systray.SetIcon(t.iconWarning)
		systray.SetTooltip("Oreon Defense - Warning")
		t.showNotification(NotificationRulesOutdated, "Rules Outdated", "Your security rules are out of date")
	case "alert":
		systray.SetIcon(t.iconAlert)
		systray.SetTooltip("Oreon Defense - Alert!")
		t.showNotification(NotificationThreatBlocked, "Threat Blocked", "A potential threat has been blocked")
	case "scanning":
		systray.SetIcon(t.iconScanning)
		systray.SetTooltip("Oreon Defense - Scanning...")
	case "paused":
		systray.SetIcon(t.iconPaused)
		systray.SetTooltip("Oreon Defense - Paused")
		t.showNotification(NotificationFirewallDisabled, "Firewall Disabled", "Your firewall protection is currently disabled")
	}
}

// loadIcons loads all the required icons
func (t *Tray) loadIcons() {
	// These will be implemented in icons.go
	t.iconProtected = loadIcon("protected")
	t.iconWarning = loadIcon("warning")
	t.iconAlert = loadIcon("alert")
	t.iconScanning = loadIcon("scanning")
	t.iconPaused = loadIcon("paused")
}

// monitorStatus periodically checks the system status and updates the UI
func (t *Tray) monitorStatus() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		status, err := t.client.Status()
		if err != nil {
			t.setIcon("warning")
			continue
		}
		t.setIcon(status.State)
	}
}

// showNotification displays a desktop notification
