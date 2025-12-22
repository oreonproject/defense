// oreon/defense · cavaire3d <C3D>
// Icons by serval_serval_serval_serval

package tray

import (
	_ "embed"
)


//go:embed icons/secure-16.png
var ProtectedIcon []byte

//go:embed icons/warning-16.png
var WarningIcon []byte

//go:embed icons/critical-16.png
var AlertIcon []byte

//go:embed icons/pending-16.png
var ScanningIcon []byte

//go:embed icons/paused-16.png
var PausedIcon []byte

// loadIcon returns the appropriate icon for the given state
func loadIcon(state string) []byte {
	switch state {
	case "protected":
		return ProtectedIcon
	case "warning":
		return WarningIcon
	case "alert":
		return AlertIcon
	case "scanning":
		return ScanningIcon
	case "paused":
		return PausedIcon
	default:
		return ProtectedIcon
	}
}
