package tray

import (
	_ "embed"
)

//go:embed assets/icons/protected.png
var protectedIcon []byte

//go:embed assets/icons/warning.png
var warningIcon []byte

// go:embed assets/icons/alert.png
var alertIcon []byte

//go:embed assets/icons/scanning.png
var scanningIcon []byte

//go:embed assets/icons/paused.png
var pausedIcon []byte

// loadIcon returns the appropriate icon for the given state
func loadIcon(state string) []byte {
	switch state {
	case "protected":
		return protectedIcon
	case "warning":
		return warningIcon
	case "alert":
		return alertIcon
	case "scanning":
		return scanningIcon
	case "paused":
		return pausedIcon
	default:
		return protectedIcon // Default to protected icon
	}
}
