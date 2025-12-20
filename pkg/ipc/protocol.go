package ipc

import (
	"encoding/json"
	"time"
)

// Request is sent from client (tray, CLI) to daemon.
//
// Example requests:
//
//	{"id": "1", "cmd": "status"}
//	{"id": "2", "cmd": "scan", "params": {"type": "quick"}}
//	{"id": "3", "cmd": "pause", "params": {"duration": "15m"}}
type Request struct {
	ID      string          `json:"id"`               // unique request ID for matching responses
	Command string          `json:"cmd"`              // command name
	Params  json.RawMessage `json:"params,omitempty"` // command-specific params
}

// Response is sent from daemon to client.
type Response struct {
	ID      string      `json:"id"`              // matches request ID
	Success bool        `json:"success"`         // true if command succeeded
	Data    interface{} `json:"data,omitempty"`  // command-specific response data
	Error   string      `json:"error,omitempty"` // error message if success=false
}

// Commands - use these constants instead of raw strings.
const (
	CmdStatus   = "status"   // get current daemon state
	CmdPing     = "ping"     // health check
	CmdScan     = "scan"     // start a scan
	CmdPause    = "pause"    // pause protection
	CmdResume   = "resume"   // resume protection
	CmdFirewall = "firewall" // firewall control

	// Firewall commands (pan will implement these)
	CmdFirewallStatus  = "firewall_status"
	CmdFirewallEnable  = "firewall_enable"
	CmdFirewallDisable = "firewall_disable"

	// Scan commands
	CmdScanQuick   = "scan_quick"
	CmdScanFull    = "scan_full"
	CmdScanStatus  = "scan_status"
	CmdScanCancel  = "scan_cancel"
	CmdScanHistory = "scan_history"

	// Rule updates
	CmdRulesStatus = "rules_status"
	CmdRulesUpdate = "rules_update"
)

// StatusResponse is returned by CmdStatus.
//
// Example (josh will use this for tray icon):
//
//	resp, _ := client.Call(CmdStatus, nil)
//	status := resp.Data.(*StatusResponse)
//	updateTrayIcon(status.State)
type StatusResponse struct {
	State           string    `json:"state"`            // "protected", "warning", etc
	FirewallEnabled bool      `json:"firewall_enabled"` // pan's firewall integration
	LastScan        time.Time `json:"last_scan"`
	RulesUpdated    time.Time `json:"rules_updated"`
}

// ScanParams for CmdScan.
type ScanParams struct {
	Type string `json:"type"` // "quick" or "full"
}

// ScanResponse is returned when starting a scan.
type ScanResponse struct {
	JobID string `json:"job_id"`
}

// ScanStatusResponse is returned by CmdScanStatus.
type ScanStatusResponse struct {
	JobID        string    `json:"job_id"`
	Status       string    `json:"status"` // "running", "completed", "cancelled"
	Progress     float64   `json:"progress"`
	FilesScanned int       `json:"files_scanned"`
	ThreatsFound int       `json:"threats_found"`
	StartedAt    time.Time `json:"started_at"`
}

// PauseParams for CmdPause.
type PauseParams struct {
	Duration string `json:"duration"` // "15m", "1h", "reboot"
}

// FirewallStatusResponse is returned by CmdFirewallStatus.
// Pan will implement the firewall package that provides this data.
type FirewallStatusResponse struct {
	Enabled    bool `json:"enabled"`
	TableCount int  `json:"table_count"`
	ChainCount int  `json:"chain_count"`
	RuleCount  int  `json:"rule_count"`
}
