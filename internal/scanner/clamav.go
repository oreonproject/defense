// oreon/defense · watchthelight <wtl>

package scanner

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

// ClamAV provides an interface to the ClamAV daemon.
type ClamAV struct {
	socketPath string
}

// New creates a new ClamAV scanner instance.
func New(socketPath string) *ClamAV {
	return &ClamAV{
		socketPath: socketPath,
	}
}

// IsAvailable checks if the ClamAV daemon is reachable.
func (c *ClamAV) IsAvailable() bool {
	if _, err := os.Stat(c.socketPath); err != nil {
		return false
	}
	return c.Ping() == nil
}

// Ping sends a PING command to clamd and expects PONG.
func (c *ClamAV) Ping() error {
	conn, err := net.DialTimeout("unix", c.socketPath, 5*time.Second)
	if err != nil {
		return fmt.Errorf("connect to clamd: %w", err)
	}
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(5 * time.Second))

	_, err = conn.Write([]byte("PING\n"))
	if err != nil {
		return fmt.Errorf("send PING: %w", err)
	}

	reader := bufio.NewReader(conn)
	response, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	if strings.TrimSpace(response) != "PONG" {
		return fmt.Errorf("unexpected response: %s", response)
	}

	return nil
}
