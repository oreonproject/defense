// oreon/defense · cavaire3d <C3D>

package tray

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
	_ "image/png"
	"log"
	"log/slog"
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
//
//go:embed icons/logo.png
var logo []byte

func (t *Tray) showNotification(notificationType NotificationType, title, message string) {
	n := notify.Notification{
		AppName:       "Oreon Defense",
		Summary:       title,
		Body:          message,
		ExpireTimeout: 10 * time.Second,
		Hints:         make(map[string]dbus.Variant),
	}
	decoded, err := decodeBytesToRGBA(logo)
	if err != nil {
		slog.Error("failed to decode image", "error", err)
	}
	hintIcon := notify.HintImageDataRGBA(decoded)
	n.AddHint(hintIcon)
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
	case NotificationRulesOutdated:
		n.Actions = []notify.Action{
			{Key: "update", Label: "Update Now"},
			{Key: "remind", Label: "Remind Later"},
		}
	case NotificationScanComplete:
		n.Actions = []notify.Action{
			{Key: "view_results", Label: "View Results"},
			{Key: "dismiss", Label: "Dismiss"},
		}
	case NotificationStateChange:
		// No actions for state change notifications
		n.ExpireTimeout = 5 * time.Second
	case None:
		// No actions for generic notifications
		n.Hints["urgency"] = dbus.MakeVariant(byte(1))
	}

	// Send the notification and get the notification ID
	if t.notifier != nil {
		id, err := t.notifier.SendNotification(n)
		if err != nil {
			slog.Error("failed to show notification", "error", err)
			return
		}
		log.Println("Notification sent with ID:", id)
		// Listen for actions

	}
}

// executeOrder66 runs a command with the given arguments
func (t *Tray) executeOrder66(command string, args []string) error {
	cmd := exec.Command(command, args...)
	slog.Debug("executing order", "command", command, "args", args)
	if err := cmd.Start(); err != nil {
		slog.Error("failed to execute command", "command", command, "args", args, "error", err)
		return err
	}
	return nil
}
func decodeBytesToRGBA(imageData []byte) (*image.RGBA, error) {
	// Decode the byte slice into a generic image.Image
	img, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	// Convert the decoded image to *image.RGBA format if it isn't already
	if rgbaImg, ok := img.(*image.RGBA); ok {
		return rgbaImg, nil
	}

	// If it's another format (like image.NRGBA, image.YCbCr, etc.),
	// draw it onto a new image.RGBA canvas.
	bounds := img.Bounds()
	rgbaImg := image.NewRGBA(bounds)
	// Drawing an image onto an RGBA canvas converts the pixels
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			rgbaImg.Set(x, y, img.At(x, y))
		}
	}
	return rgbaImg, nil
}
