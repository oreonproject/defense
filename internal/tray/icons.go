// oreon/defense · cavaire3d <C3D>

package tray

import "github.com/oreonproject/defense/assets"

// loadIcon returns the appropriate icon for the given state
func loadIcon(state string) []byte {
	switch state {
	case "protected":
		return assets.ProtectedIcon
	case "warning":
		return assets.WarningIcon
	case "alert":
		return assets.AlertIcon
	case "scanning":
		return assets.ScanningIcon
	case "paused":
		return assets.PausedIcon
	default:
		return assets.ProtectedIcon
	}
}
