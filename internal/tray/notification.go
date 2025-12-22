// notification.go
package tray

import (
	"log"
	"os/exec"
	"time"

	"github.com/esiqveland/notify"
	"github.com/godbus/dbus/v5"
)

type NotificationType string

const (
	None                         NotificationType = "none"
	NotificationFirewallDisabled NotificationType = "firewall_disabled"
	NotificationRulesOutdated    NotificationType = "rules_outdated"
	NotificationScanComplete     NotificationType = "scan_complete"
	NotificationThreatBlocked    NotificationType = "threat_blocked"
	NotificationStateChange      NotificationType = "state_change"
)

// Tray embeds the system tray functionality

// showNotification shows a system notification with optional actions
func (t *Tray) showNotification(notificationType NotificationType, title, message string) {
    n := notify.Notification{
        AppName:       "Oreon Defense",
        Summary:       title,
        Body:          message,
        ExpireTimeout: 10 * time.Second,
        Hints:         make(map[string]dbus.Variant),
    }

    // Set all notifications as critical
    n.Hints["urgency"] = dbus.MakeVariant(byte(2)) // Critical

    // Add actions based on notification type
    switch notificationType {
    case NotificationFirewallDisabled:
        n.Actions = []notify.Action{
            {Key: "enable", Label: "Enable Now"},
            {Key: "remind", Label: "Remind Later"},
        }
    case NotificationThreatBlocked:
        n.Actions = []notify.Action{
            {Key: "details", Label: "View Details"},
        }
    }

    // Send the notification
    // if _, err := notify.SendNotification(t.conn,n); err != nil {
    //     log.Printf("Failed to show notification: %v", err)
    // }
}
// executeAction runs a command with the given arguments
func (t *Tray) executeAction(command string, args []string) error {
	cmd := exec.Command(command, args...)
	log.Printf("Executing: %s %v", command, args)
	return cmd.Start()
}

// Close cleans up resources
// func (t *Tray) Close() {
// 	if t.conn != nil {
// 		t.conn.Close()
// 	}
// }