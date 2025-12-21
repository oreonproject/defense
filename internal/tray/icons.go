package tray

import (
	"encoding/base64"
	"log/slog"
)

// decodeIcon decodes a base64 icon string, logging a warning on failure.
func decodeIcon(name, data string) []byte {
	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		slog.Warn("failed to decode icon", "icon", name, "error", err)
		return nil
	}
	return decoded
}

// Placeholder icon data (16x16 transparent PNG) - will be replaced with real icons
const placeholderIconData = "iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAYAAAAf8/9hAAAABmJLR0QA/wD/AP+gvaeTAAAACXBIWXMAAAsTAAALEwEAmpwYAAAAB3RJTUUH4AkEEjQaTQ+3JQAAAB1pVFh0Q29tbWVudAAAAAAAQ3JlYXRlZCB3aXRoIEdJTVBkLmUHAAAAGUlEQVQ4y2NgGAWjYBSMglEwCkbBKBgFgw4AABAAAN1I3kMAAAAASUVORK5CYII="

var protectedIcon = decodeIcon("protected", placeholderIconData)

// TODO: Replace with actual distinct icons from design team
var warningIcon = decodeIcon("warning", placeholderIconData)
var alertIcon = decodeIcon("alert", placeholderIconData)
var scanningIcon = decodeIcon("scanning", placeholderIconData)
var pausedIcon = decodeIcon("paused", placeholderIconData)

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
