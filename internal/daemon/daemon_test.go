// oreon/defense · watchthelight <wtl>

package daemon

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/oreonproject/defense/pkg/config"
)

func TestDaemonRun(t *testing.T) {
	cfg := &config.Config{}
	logger := slog.Default()

	d := New(cfg, logger)

	socketPath := t.TempDir() + "/test.sock"

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	err := d.Run(ctx, socketPath)
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	if d.State().State() != StateProtected {
		t.Errorf("state = %v, want %v", d.State().State(), StateProtected)
	}
}
