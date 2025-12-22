package tray

import (
	"fmt"

	"github.com/energye/systray"
)

// menu represents the system tray menu structure
type menu struct {
	tray *Tray

	// Main menu items
	quickScanItem    *systray.MenuItem
	fullScanItem     *systray.MenuItem
	updateRulesItem  *systray.MenuItem
	pauseMenu        *systray.MenuItem
	pause15MinItem   *systray.MenuItem
	pause1HourItem   *systray.MenuItem
	pauseUntilReboot *systray.MenuItem
	firewallItem     *systray.MenuItem
	alertsMenu       *systray.MenuItem
	openAppItem      *systray.MenuItem
	settingsItem     *systray.MenuItem
	quitItem         *systray.MenuItem

	// State tracking
	isPaused bool
}

// newMenu creates a new menu instance
func newMenu(t *Tray) *menu {
	return &menu{
		tray: t,
	}
}

// build creates and initializes the menu structure
func (m *menu) build() {
	// Add menu items
	m.quickScanItem = systray.AddMenuItem("Run Quick Scan", "Run a quick system scan")
	m.fullScanItem = systray.AddMenuItem("Run Full Scan", "Run a full system scan")
	systray.AddSeparator()

	m.updateRulesItem = systray.AddMenuItem("Update Security Rules", "Download and update security rules")
	systray.AddSeparator()

	// Pause protection submenu
	m.pauseMenu = systray.AddMenuItem("Pause Protection", "Temporarily pause protection")
	m.pause15MinItem = m.pauseMenu.AddSubMenuItem("15 Minutes", "Pause protection for 15 minutes")
	m.pause1HourItem = m.pauseMenu.AddSubMenuItem("1 Hour", "Pause protection for 1 hour")
	m.pauseUntilReboot = m.pauseMenu.AddSubMenuItem("Until Reboot", "Pause protection until next system reboot")

	systray.AddSeparator()

	// Firewall toggle
	m.firewallItem = systray.AddMenuItemCheckbox("Firewall: Enabled ✓", "Toggle firewall protection", true)

	systray.AddSeparator()

	// Alerts submenu
	m.alertsMenu = systray.AddMenuItem("Recent Alerts (0)", "View recent security alerts")

	systray.AddSeparator()

	// App actions
	m.openAppItem = systray.AddMenuItem("Open Oreon Defense", "Open main application window")
	m.settingsItem = systray.AddMenuItem("Settings", "Configure application settings")
	m.quitItem = systray.AddMenuItem("Quit", "Exit the application")

	// Set up click handlers
	m.setupHandlers()
}

// setupHandlers configures all menu item click handlers
func (m *menu) setupHandlers() {
	// Sync state with daemon
	go m.syncStateWithDaemon()

	// Set up click handlers for each menu item
	m.quickScanItem.Click(m.handleQuickScan)
	m.fullScanItem.Click(m.handleFullScan)
	m.updateRulesItem.Click(m.handleUpdateRules)
	m.pause15MinItem.Click(func() { m.handlePause("15m") })
	m.pause1HourItem.Click(func() { m.handlePause("1h") })
	m.pauseUntilReboot.Click(func() { m.handlePause("reboot") })
	m.firewallItem.Click(m.handleFirewallToggle)
	m.openAppItem.Click(m.handleOpenApp)
	m.settingsItem.Click(m.handleOpenSettings)
	m.quitItem.Click(systray.Quit)
}

func (m *menu) handleQuickScan() {
	m.tray.setIcon("scanning")
	go func() {
		_, err := m.tray.client.StartQuickScan()
		if err != nil {
			m.tray.setIcon("warning")
			m.tray.showNotification(None, "Scan Failed", "Failed to start quick scan: "+err.Error())
			return
		}
		// The daemon will set state back to protected when scan completes
		// We show notification here as scan has started
		m.tray.showNotification(NotificationScanComplete, "Quick Scan Started", "Scanning your system...")
	}()
}

func (m *menu) handleFullScan() {
	m.tray.setIcon("scanning")
	go func() {
		_, err := m.tray.client.StartFullScan()
		if err != nil {
			m.tray.setIcon("warning")
			m.tray.showNotification(None, "Scan Failed", "Failed to start full scan: "+err.Error())
			return
		}
		m.tray.showNotification(NotificationScanComplete, "Full Scan Started", "Scanning your entire system...")
	}()
}
func (m *menu) handleUpdateRules() {
	// TODO: Implement rules update
	go func() {
		// err := m.tray.client.UpdateRules()
		// if err != nil {
		// 	m.tray.showNotification("Update Failed", "Failed to update security rules", true)
		// 	return
		// }
		m.tray.showNotification(None, "Rules Updated", "Security rules have been updated successfully")
	}()
}

func (m *menu) handlePause(duration string) {
	// Check if we're resuming
	if m.isPaused {
		err := m.tray.client.Resume()
		if err != nil {
			m.tray.showNotification(None, "Error", "Failed to resume protection: "+err.Error())
			return
		}
		m.isPaused = false
		m.pauseMenu.SetTitle("Pause Protection")
		m.tray.showNotification(NotificationStateChange, "Protection Resumed", "Protection is now active")
		return
	}

	// Pause protection
	err := m.tray.client.Pause()
	if err != nil {
		m.tray.showNotification(None, "Error", "Failed to pause protection: "+err.Error())
		return
	}

	m.isPaused = true
	m.pauseMenu.SetTitle("Resume Protection")

	switch duration {
	case "15m":
		m.tray.showNotification(NotificationStateChange, "Protection Paused", "Protection will resume in 15 minutes")
	case "1h":
		m.tray.showNotification(NotificationStateChange, "Protection Paused", "Protection will resume in 1 hour")
	case "reboot":
		m.tray.showNotification(NotificationStateChange, "Protection Paused", "Protection will resume after system reboot")
	}
}

// syncStateWithDaemon queries the daemon and updates menu state to match
func (m *menu) syncStateWithDaemon() {
	status, err := m.tray.client.Status()
	if err != nil {
		// Can't reach daemon, leave defaults
		return
	}

	// Sync firewall checkbox
	if status.FirewallEnabled {
		m.firewallItem.Check()
		m.firewallItem.SetTitle("Firewall: Enabled ✓")
	} else {
		m.firewallItem.Uncheck()
		m.firewallItem.SetTitle("Firewall: Disabled")
	}

	// Sync pause state
	if status.State == "paused" {
		m.isPaused = true
		m.pauseMenu.SetTitle("Resume Protection")
	} else {
		m.isPaused = false
		m.pauseMenu.SetTitle("Pause Protection")
	}
}

func (m *menu) handleFirewallToggle() {
	newState := !m.firewallItem.Checked()

	err := m.tray.client.SetFirewallEnabled(newState)
	if err != nil {
		m.tray.showNotification(None, "Error", "Failed to toggle firewall: "+err.Error())
		return
	}

	if newState {
		m.firewallItem.Check()
		m.firewallItem.SetTitle("Firewall: Enabled ✓")
		m.firewallItem.SetTooltip("Firewall protection is enabled")
	} else {
		m.firewallItem.Uncheck()
		m.firewallItem.SetTitle("Firewall: Disabled")
		m.firewallItem.SetTooltip("Firewall protection is disabled")
		m.tray.showNotification(NotificationFirewallDisabled, "Firewall Disabled", "Your firewall protection is now disabled")
	}
}

func (m *menu) handleOpenApp() {
	// TODO: Implement app open
	fmt.Println("Open app requested")
}

func (m *menu) handleOpenSettings() {
	// TODO: Open settings
	fmt.Println("Open settings requested")
}

// updateAlertsMenu updates the alerts menu with the latest alerts
func (m *menu) updateAlertsMenu() {
	// TODO: Implement alerts menu update
	// alerts, err := m.tray.client.GetRecentAlerts(5)
	// if err != nil {
	// 	return
	// }

	// Update menu title
	// m.alertsMenu.SetTitle(fmt.Sprintf("Recent Alerts (%d)", len(alerts)))

	// Clear existing menu items
	// for _, item := range m.alertsMenu.GetSubMenu().Items() {
	// 	item.Hide()
	// }

	// Add new alerts
	// for _, alert := range alerts {
	// 	item := m.alertsMenu.AddSubMenuItem(alert.Title, alert.Message)
	// 	// Add click handler for each alert
	// }
}
