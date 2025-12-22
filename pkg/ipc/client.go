package ipc

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
)

// Client defines the interface for IPC communication with the daemon
type Client interface {
	Status() (*StatusResponse, error)
	GetProtectionStatus() (bool, error)
	SetFirewallEnabled(enabled bool) error
	IsFirewallEnabled() (bool, error)
	Close() error
}

// socketClient is the real IPC client implementation
type socketClient struct {
	socketPath string
	conn       net.Conn
	reader     *bufio.Reader
	mu         sync.Mutex
	reqID      atomic.Uint64
}

// NewClient creates a new IPC client connected to the daemon
func NewClient(socketPath string) (Client, error) {
	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		return nil, fmt.Errorf("connect to daemon: %w", err)
	}
	return &socketClient{
		socketPath: socketPath,
		conn:       conn,
		reader:     bufio.NewReader(conn),
	}, nil
}

func (c *socketClient) call(cmd string, params interface{}) (*Response, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	id := fmt.Sprintf("%d", c.reqID.Add(1))

	req := Request{ID: id, Command: cmd}
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

	// Re-marshal and unmarshal to get typed response
	data, err := json.Marshal(resp.Data)
	if err != nil {
		return nil, err
	}
	var status StatusResponse
	if err := json.Unmarshal(data, &status); err != nil {
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

func (c *socketClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
