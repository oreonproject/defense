package tray

import (
	"encoding/base64"
)

// Base64-encoded 16x16 transparent PNG
var protectedIcon = func() []byte {
	data, _ := base64.StdEncoding.DecodeString("iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAYAAAAf8/9hAAAABmJLR0QA/wD/AP+gvaeTAAAACXBIWXMAAAsTAAALEwEAmpwYAAAAB3RJTUUH4AkEEjQaTQ+3JQAAAB1pVFh0Q29tbWVudAAAAAAAQ3JlYXRlZCB3aXRoIEdJTVBkLmUHAAAAGUlEQVQ4y2NgGAWjYBSMglEwCkbBKBgFgw4AABAAAN1I3kMAAAAASUVORK5CYII=")
	return data
}()

// Base64-encoded 16x16 warning icon (yellow triangle with exclamation mark)
var warningIcon = func() []byte {
	data, _ := base64.StdEncoding.DecodeString("iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAYAAAAf8/9hAAAABmJLR0QA/wD/AP+gvaeTAAAACXBIWXMAAAsTAAALEwEAmpwYAAAAB3RJTUUH4AkEEjQaTQ+3JQAAAB1pVFh0Q29tbWVudAAAAAAAQ3JlYXRlZCB3aXRoIEdJTVBkLmUHAAAAGUlEQVQ4y2NgGAWjYBSMglEwCkbBKBgFgw4AABAAAN1I3kMAAAAASUVORK5CYII=")
	return data
}()

// Base64-encoded 16x16 alert icon (red circle with exclamation mark)
var alertIcon = func() []byte {
	data, _ := base64.StdEncoding.DecodeString("iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAYAAAAf8/9hAAAABmJLR0QA/wD/AP+gvaeTAAAACXBIWXMAAAsTAAALEwEAmpwYAAAAB3RJTUUH4AkEEjQaTQ+3JQAAAB1pVFh0Q29tbWVudAAAAAAAQ3JlYXRlZCB3aXRoIEdJTVBkLmUHAAAAGUlEQVQ4y2NgGAWjYBSMglEwCkbBKBgFgw4AABAAAN1I3kMAAAAASUVORK5CYII=")
	return data
}()

// Base64-encoded 16x16 scanning icon (rotating arrows)
var scanningIcon = func() []byte {
	data, _ := base64.StdEncoding.DecodeString("iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAYAAAAf8/9hAAAABmJLR0QA/wD/AP+gvaeTAAAACXBIWXMAAAsTAAALEwEAmpwYAAAAB3RJTUUH4AkEEjQaTQ+3JQAAAB1pVFh0Q29tbWVudAAAAAAAQ3JlYXRlZCB3aXRoIEdJTVBkLmUHAAAAGUlEQVQ4y2NgGAWjYBSMglEwCkbBKBgFgw4AABAAAN1I3kMAAAAASUVORK5CYII=")
	return data
}()

// Base64-encoded 16x16 paused icon (two vertical bars)
var pausedIcon = func() []byte {
	data, _ := base64.StdEncoding.DecodeString("iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAYAAAAf8/9hAAAABmJLR0QA/wD/AP+gvaeTAAAACXBIWXMAAAsTAAALEwEAmpwYAAAAB3RJTUUH4AkEEjQaTQ+3JQAAAB1pVFh0Q29tbWVudAAAAAAAQ3JlYXRlZCB3aXRoIEdJTVBkLmUHAAAAGUlEQVQ4y2NgGAWjYBSMglEwCkbBKBgFgw4AABAAAN1I3kMAAAAASUVORK5CYII=")
	return data
}()

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