// oreon/defense · watchthelight <wtl>

package daemon

import (
	"bufio"
	"encoding/json"
	"log/slog"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/oreonproject/defense/pkg/ipc"
)

// Server handles IPC connections from clients (tray, CLI).
type Server struct {
	socketPath string
	listener   net.Listener
	state      *StateManager
	done       chan struct{}
}

// NewServer creates an IPC server that exposes daemon state.
func NewServer(socketPath string, state *StateManager) *Server {
	return &Server{
		socketPath: socketPath,
		state:      state,
		done:       make(chan struct{}),
	}
}

// Listen creates the unix socket and starts accepting connections.
func (s *Server) Listen() error {
	// Ensure parent directory exists
	if err := os.MkdirAll(filepath.Dir(s.socketPath), 0755); err != nil {
		return err
	}

	// Remove stale socket if it exists
	os.Remove(s.socketPath)

	ln, err := net.Listen("unix", s.socketPath)
	if err != nil {
		return err
	}
	s.listener = ln

	// Set socket permissions (rw for owner and group)
	if err := os.Chmod(s.socketPath, 0660); err != nil {
		ln.Close()
		return err
	}

	slog.Info("IPC server listening", "socket", s.socketPath)
	return nil
}

// Serve accepts connections until Close is called.
func (s *Server) Serve() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			select {
			case <-s.done:
				return // shutdown
			default:
				slog.Warn("accept error", "error", err)
				continue
			}
		}
		go s.handleConnection(conn)
	}
}

// Close shuts down the server.
func (s *Server) Close() error {
	close(s.done)
	if s.listener != nil {
		s.listener.Close()
	}
	os.Remove(s.socketPath)
	return nil
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	encoder := json.NewEncoder(conn)

	for {
		// Read one line (one JSON request)
		line, err := reader.ReadBytes('\n')
		if err != nil {
			return // client disconnected
		}

		var req ipc.Request
		if err := json.Unmarshal(line, &req); err != nil {
			encoder.Encode(ipc.Response{
				Success: false,
				Error:   "invalid JSON",
			})
			continue
		}

		resp := s.handleRequest(&req)
		encoder.Encode(resp)
	}
}

func (s *Server) handleRequest(req *ipc.Request) *ipc.Response {
	switch req.Command {
	case ipc.CmdPing:
		return &ipc.Response{
			ID:      req.ID,
			Success: true,
			Data:    "pong",
		}

	case ipc.CmdStatus:
		return &ipc.Response{
			ID:      req.ID,
			Success: true,
			Data: ipc.StatusResponse{
				State:           s.state.State().String(),
				FirewallEnabled: false, // TODO: wire to firewall
				LastScan:        time.Time{},
				RulesUpdated:    time.Time{},
			},
		}

	default:
		return &ipc.Response{
			ID:      req.ID,
			Success: false,
			Error:   "unknown command: " + req.Command,
		}
	}
}
