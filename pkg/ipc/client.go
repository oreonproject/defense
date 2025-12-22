package ipc

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

// Client defines the interface for IPC communication with the daemon
type Client interface {
	Status() (*StatusResponse, error)
	GetProtectionStatus() (bool, error)
	SetFirewallEnabled(enabled bool) error
	IsFirewallEnabled() (bool, error)
	StartQuickScan() (*ScanResponse, error)
	StartFullScan() (*ScanResponse, error)
	Pause() error
	Resume() error
	Subscribe() (<-chan StateChangeEvent, error)
	Close() error
}

// socketClient is the real IPC client implementation
type socketClient struct {
	socketPath string
	conn       net.Conn
	reader     *bufio.Reader
	mu         sync.Mutex
	reqID      atomic.Uint64
	connected  bool
}

// NewClient creates a new IPC client. Connection is established lazily on first call.
func NewClient(socketPath string) Client {
	return &socketClient{socketPath: socketPath}
}

// connect establishes a connection to the daemon
func (c *socketClient) connect() error {
	conn, err := net.Dial("unix", c.socketPath)
	if err != nil {
		return fmt.Errorf("connect to daemon: %w", err)
	}
	c.conn = conn
	c.reader = bufio.NewReader(conn)
	c.connected = true
	return nil
}

// reconnect closes the existing connection and establishes a new one
func (c *socketClient) reconnect() error {
	if c.conn != nil {
		c.conn.Close()
	}
	c.connected = false
	return c.connect()
}

func (c *socketClient) call(cmd string, params interface{}) (*Response, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Connect on first use
	if !c.connected {
		if err := c.connect(); err != nil {
			return nil, err
		}
	}

	resp, err := c.doCall(cmd, params)
	if err != nil && c.isConnectionError(err) {
		// Try to reconnect once
		if reconnErr := c.reconnect(); reconnErr != nil {
			return nil, fmt.Errorf("reconnect failed: %w", reconnErr)
		}
		// Retry the call
		resp, err = c.doCall(cmd, params)
	}
	return resp, err
}

// isConnectionError checks if the error indicates a broken connection
func (c *socketClient) isConnectionError(err error) bool {
	if err == nil {
		return false
	}
	// EOF means connection closed
	if errors.Is(err, io.EOF) {
		return true
	}
	// Check for network operation errors
	var opErr *net.OpError
	if errors.As(err, &opErr) {
		return true
	}
	return false
}

// doCall performs the actual IPC call without reconnect logic
func (c *socketClient) doCall(cmd string, params interface{}) (*Response, error) {
	id := fmt.Sprintf("%d", c.reqID.Add(1))

	req := Request{Version: ProtocolVersion, ID: id, Command: cmd}
	if params != nil {
		p, err := json.Marshal(params)
		if err != nil {
			return nil, err
		}
		req.Params = p
	}

	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	data = append(data, '\n')

	if _, err := c.conn.Write(data); err != nil {
		return nil, err
	}

	// Set read deadline to avoid blocking forever if daemon hangs
	c.conn.SetReadDeadline(time.Now().Add(30 * time.Second))
	line, err := c.reader.ReadBytes('\n')
	if err != nil {
		return nil, err
	}

	var resp Response
	if err := json.Unmarshal(line, &resp); err != nil {
		return nil, err
	}

	if !resp.Success {
		return nil, fmt.Errorf("daemon error: %s", resp.Error)
	}

	return &resp, nil
}

func (c *socketClient) Status() (*StatusResponse, error) {
	resp, err := c.call(CmdStatus, nil)
	if err != nil {
		return nil, err
	}

	var status StatusResponse
	if err := resp.UnmarshalData(&status); err != nil {
		return nil, err
	}
	return &status, nil
}

func (c *socketClient) GetProtectionStatus() (bool, error) {
	status, err := c.Status()
	if err != nil {
		return false, err
	}
	return status.State == "protected", nil
}

func (c *socketClient) SetFirewallEnabled(enabled bool) error {
	cmd := CmdFirewallDisable
	if enabled {
		cmd = CmdFirewallEnable
	}
	_, err := c.call(cmd, nil)
	return err
}

func (c *socketClient) IsFirewallEnabled() (bool, error) {
	status, err := c.Status()
	if err != nil {
		return false, err
	}
	return status.FirewallEnabled, nil
}

func (c *socketClient) StartQuickScan() (*ScanResponse, error) {
	resp, err := c.call(CmdScanQuick, nil)
	if err != nil {
		return nil, err
	}

	var scanResp ScanResponse
	if err := resp.UnmarshalData(&scanResp); err != nil {
		return nil, err
	}
	return &scanResp, nil
}

func (c *socketClient) StartFullScan() (*ScanResponse, error) {
	resp, err := c.call(CmdScanFull, nil)
	if err != nil {
		return nil, err
	}

	var scanResp ScanResponse
	if err := resp.UnmarshalData(&scanResp); err != nil {
		return nil, err
	}
	return &scanResp, nil
}

func (c *socketClient) Pause() error {
	_, err := c.call(CmdPause, nil)
	return err
}

func (c *socketClient) Resume() error {
	_, err := c.call(CmdResume, nil)
	return err
}

func (c *socketClient) Subscribe() (<-chan StateChangeEvent, error) {
	// Create a dedicated connection for subscription
	conn, err := net.Dial("unix", c.socketPath)
	if err != nil {
		return nil, fmt.Errorf("connect for subscribe: %w", err)
	}

	// Send subscribe request
	req := Request{Version: ProtocolVersion, ID: "sub", Command: CmdSubscribe}
	data, _ := json.Marshal(req)
	data = append(data, '\n')
	if _, err := conn.Write(data); err != nil {
		conn.Close()
		return nil, fmt.Errorf("send subscribe: %w", err)
	}

	// Read subscription confirmation
	reader := bufio.NewReader(conn)
	line, err := reader.ReadBytes('\n')
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("read subscribe response: %w", err)
	}

	var resp Response
	if err := json.Unmarshal(line, &resp); err != nil || !resp.Success {
		conn.Close()
		return nil, fmt.Errorf("subscribe failed: %s", resp.Error)
	}

	// Create channel and start reading events
	events := make(chan StateChangeEvent, 10)
	go func() {
		defer conn.Close()
		defer close(events)

		for {
			line, err := reader.ReadBytes('\n')
			if err != nil {
				return
			}

			var resp Response
			if err := json.Unmarshal(line, &resp); err != nil {
				continue
			}

			var event StateChangeEvent
			if err := resp.UnmarshalData(&event); err != nil {
				continue
			}

			select {
			case events <- event:
			default:
				// drop if channel is full
			}
		}
	}()

	return events, nil
}

func (c *socketClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
