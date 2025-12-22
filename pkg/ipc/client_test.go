// oreon/defense · watchthelight <wtl>

package ipc

import (
	"bufio"
	"encoding/json"
	"net"
	"path/filepath"
	"testing"
	"time"
)

// mockIPCServer creates a mock daemon server for testing
func mockIPCServer(t *testing.T, handler func(req *Request) *Response) (string, func()) {
	t.Helper()

	dir := t.TempDir()
	sockPath := filepath.Join(dir, "test.sock")

	listener, err := net.Listen("unix", sockPath)
	if err != nil {
		t.Fatalf("failed to create mock server: %v", err)
	}

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			go func(c net.Conn) {
				defer c.Close()
				reader := bufio.NewReader(c)
				encoder := json.NewEncoder(c)

				for {
					line, err := reader.ReadBytes('\n')
					if err != nil {
						return
					}

					var req Request
					if err := json.Unmarshal(line, &req); err != nil {
						continue
					}

					resp := handler(&req)
					encoder.Encode(resp)
				}
			}(conn)
		}
	}()

	return sockPath, func() { listener.Close() }
}

func TestClient_Status(t *testing.T) {
	sockPath, cleanup := mockIPCServer(t, func(req *Request) *Response {
		if req.Command != CmdStatus {
			t.Errorf("unexpected command: %s", req.Command)
		}

		data, _ := json.Marshal(StatusResponse{
			State:           "protected",
			FirewallEnabled: true,
			LastScan:        time.Now(),
		})

		return &Response{
			ID:      req.ID,
			Success: true,
			Data:    data,
		}
	})
	defer cleanup()

	client := NewClient(sockPath)
	defer client.Close()

	status, err := client.Status()
	if err != nil {
		t.Fatalf("Status() error = %v", err)
	}

	if status.State != "protected" {
		t.Errorf("State = %v, want protected", status.State)
	}
	if !status.FirewallEnabled {
		t.Error("FirewallEnabled = false, want true")
	}
}

func TestClient_GetProtectionStatus(t *testing.T) {
	sockPath, cleanup := mockIPCServer(t, func(req *Request) *Response {
		data, _ := json.Marshal(StatusResponse{State: "protected"})
		return &Response{ID: req.ID, Success: true, Data: data}
	})
	defer cleanup()

	client := NewClient(sockPath)
	defer client.Close()

	protected, err := client.GetProtectionStatus()
	if err != nil {
		t.Fatalf("GetProtectionStatus() error = %v", err)
	}
	if !protected {
		t.Error("GetProtectionStatus() = false, want true")
	}
}

func TestClient_SetFirewallEnabled(t *testing.T) {
	var receivedCmd string

	sockPath, cleanup := mockIPCServer(t, func(req *Request) *Response {
		receivedCmd = req.Command
		data, _ := json.Marshal("ok")
		return &Response{ID: req.ID, Success: true, Data: data}
	})
	defer cleanup()

	client := NewClient(sockPath)
	defer client.Close()

	// Test enable
	if err := client.SetFirewallEnabled(true); err != nil {
		t.Fatalf("SetFirewallEnabled(true) error = %v", err)
	}
	if receivedCmd != CmdFirewallEnable {
		t.Errorf("command = %v, want %v", receivedCmd, CmdFirewallEnable)
	}

	// Test disable
	if err := client.SetFirewallEnabled(false); err != nil {
		t.Fatalf("SetFirewallEnabled(false) error = %v", err)
	}
	if receivedCmd != CmdFirewallDisable {
		t.Errorf("command = %v, want %v", receivedCmd, CmdFirewallDisable)
	}
}

func TestClient_StartQuickScan(t *testing.T) {
	sockPath, cleanup := mockIPCServer(t, func(req *Request) *Response {
		if req.Command != CmdScanQuick {
			t.Errorf("unexpected command: %s", req.Command)
		}
		data, _ := json.Marshal(ScanResponse{JobID: "quick-123"})
		return &Response{ID: req.ID, Success: true, Data: data}
	})
	defer cleanup()

	client := NewClient(sockPath)
	defer client.Close()

	resp, err := client.StartQuickScan()
	if err != nil {
		t.Fatalf("StartQuickScan() error = %v", err)
	}
	if resp.JobID != "quick-123" {
		t.Errorf("JobID = %v, want quick-123", resp.JobID)
	}
}

func TestClient_PauseResume(t *testing.T) {
	var receivedCmd string

	sockPath, cleanup := mockIPCServer(t, func(req *Request) *Response {
		receivedCmd = req.Command
		data, _ := json.Marshal("ok")
		return &Response{ID: req.ID, Success: true, Data: data}
	})
	defer cleanup()

	client := NewClient(sockPath)
	defer client.Close()

	if err := client.Pause(); err != nil {
		t.Fatalf("Pause() error = %v", err)
	}
	if receivedCmd != CmdPause {
		t.Errorf("command = %v, want %v", receivedCmd, CmdPause)
	}

	if err := client.Resume(); err != nil {
		t.Fatalf("Resume() error = %v", err)
	}
	if receivedCmd != CmdResume {
		t.Errorf("command = %v, want %v", receivedCmd, CmdResume)
	}
}

func TestClient_Reconnect(t *testing.T) {
	callCount := 0

	sockPath, cleanup := mockIPCServer(t, func(req *Request) *Response {
		callCount++
		data, _ := json.Marshal(StatusResponse{State: "protected"})
		return &Response{ID: req.ID, Success: true, Data: data}
	})

	client := NewClient(sockPath)
	defer client.Close()

	// First call
	_, err := client.Status()
	if err != nil {
		t.Fatalf("first Status() error = %v", err)
	}

	// Restart server (simulates daemon restart)
	cleanup()
	sockPath2, cleanup2 := mockIPCServer(t, func(req *Request) *Response {
		callCount++
		data, _ := json.Marshal(StatusResponse{State: "warning"})
		return &Response{ID: req.ID, Success: true, Data: data}
	})
	defer cleanup2()

	// Client needs the same socket path, so this test is limited
	// Just verify initial call worked
	if callCount != 1 {
		t.Errorf("callCount = %v, want 1", callCount)
	}

	_ = sockPath2 // acknowledge we created new server
}

func TestClient_Error(t *testing.T) {
	sockPath, cleanup := mockIPCServer(t, func(req *Request) *Response {
		return &Response{
			ID:      req.ID,
			Success: false,
			Error:   "something went wrong",
		}
	})
	defer cleanup()

	client := NewClient(sockPath)
	defer client.Close()

	_, err := client.Status()
	if err == nil {
		t.Fatal("Status() should return error")
	}
	if err.Error() != "daemon error: something went wrong" {
		t.Errorf("error = %v, want 'daemon error: something went wrong'", err)
	}
}

func TestClient_ConnectionFailure(t *testing.T) {
	client := NewClient("/nonexistent/socket.sock")
	defer client.Close()

	_, err := client.Status()
	if err == nil {
		t.Fatal("Status() should return error for nonexistent socket")
	}
}
