package ipc

// Client defines the interface for IPC communication with the daemon
type Client interface {
	// Protection
	EnableProtection() error
	DisableProtection(duration string) error
	GetProtectionStatus() (bool, error)

	// Scanning
	StartQuickScan() (string, error)             // returns scan ID
	StartFullScan() (string, error)              // returns scan ID
	GetScanStatus(scanID string) (string, error) // returns status: "running", "completed", "failed"

	// Rules
	UpdateRules() error
	GetRulesVersion() (string, error)

	// Firewall
	SetFirewallEnabled(enabled bool) error
	IsFirewallEnabled() (bool, error)

	// Alerts
	GetRecentAlerts(limit int) ([]Alert, error)
}

// Alert represents a security alert notification
type Alert struct {
	ID        string   `json:"id"`
	Title     string   `json:"title"`
	Message   string   `json:"message"`
	Severity  string   `json:"severity"` // info, warning, critical
	Timestamp int64    `json:"timestamp"`
	Actions   []string `json:"actions,omitempty"`
}

// NewClient creates a new IPC client
func NewClient(socketPath string) (Client, error) {
	// TODO: Implement actual client
	return &mockClient{}, nil
}

// mockClient is a stub implementation for development
type mockClient struct{}

func (m *mockClient) EnableProtection() error                     { return nil }
func (m *mockClient) DisableProtection(duration string) error     { return nil }
func (m *mockClient) GetProtectionStatus() (bool, error)          { return true, nil }
func (m *mockClient) StartQuickScan() (string, error)             { return "scan-123", nil }
func (m *mockClient) StartFullScan() (string, error)              { return "scan-456", nil }
func (m *mockClient) GetScanStatus(scanID string) (string, error) { return "completed", nil }
func (m *mockClient) UpdateRules() error                          { return nil }
func (m *mockClient) GetRulesVersion() (string, error)            { return "1.0.0", nil }
func (m *mockClient) SetFirewallEnabled(enabled bool) error       { return nil }
func (m *mockClient) IsFirewallEnabled() (bool, error)            { return true, nil }
func (m *mockClient) GetRecentAlerts(limit int) ([]Alert, error) {
	return []Alert{{
		ID:        "alert-1",
		Title:     "Example Alert",
		Message:   "This is a sample alert",
		Severity:  "info",
		Timestamp: 1671561600,
	}}, nil
}
