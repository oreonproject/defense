// Package tray implements the system tray interface for Oreon Defense
package tray

import (
	"sync"

	"github.com/energye/systray"
	"github.com/oreonproject/defense/pkg/ipc"
)

// Tray represents the system tray application
type Tray struct {
	client        ipc.Client
	menu          *menu
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
	return &Tray{
		client: client,
	}
}

// Run starts the system tray application
func (t *Tray) Run() error {
	systray.Run(t.onReady, t.onExit)
	return nil
}

// onReady is called when the system tray is ready
func (t *Tray) onReady() {
	// Load icons
	t.loadIcons()

	// Set initial state
	t.setIcon("protected")

	// Build the menu
	t.menu = newMenu(t)
	t.menu.build()

	// Start status monitoring
	go t.monitorStatus()
}

// onExit is called when the system tray is exiting
func (t *Tray) onExit() {
	// Cleanup resources if needed
}

// setIcon updates the tray icon based on the current state
func (t *Tray) setIcon(state string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.currentState == state {
		return
	}

	t.currentState = state

	switch state {
	case "protected":
		systray.SetIcon(t.iconProtected)
		systray.SetTooltip("Oreon Defense - Protected")
	case "warning":
		systray.SetIcon(t.iconWarning)
		systray.SetTooltip("Oreon Defense - Warning")
	case "alert":
		systray.SetIcon(t.iconAlert)
		systray.SetTooltip("Oreon Defense - Alert!")
	case "scanning":
		systray.SetIcon(t.iconScanning)
		systray.SetTooltip("Oreon Defense - Scanning...")
	case "paused":
		systray.SetIcon(t.iconPaused)
		systray.SetTooltip("Oreon Defense - Paused")
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
	// TODO: Implement actual status monitoring
	// This will check protection status, scan status, etc.
	// and update the tray icon and menu items accordingly
}

// showNotification displays a desktop notification
func (t *Tray) showNotification(title, message string, isError bool) {
	// TODO: Implement desktop notifications
	// This will use the appropriate notification system for the platform
}
