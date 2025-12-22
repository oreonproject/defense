// oreon/defense · watchthelight <wtl>

package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/oreonproject/defense/internal/daemon"
	"github.com/oreonproject/defense/pkg/config"
)

var version = "0.1.0-dev"

func main() {
	configPath := flag.String("config", config.SystemConfigPath, "path to config file")
	debug := flag.Bool("debug", false, "enable debug logging")
	flag.Parse()

	fmt.Printf("Oreon Defense v%s\n", version)

	if *debug {
		fmt.Println("debug mode enabled")
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := run(ctx, *configPath); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, configPath string) error {
	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	slog.Info("config loaded", "path", configPath)

	d := daemon.New(cfg, slog.Default())
	return d.Run(ctx, config.SocketPath)
}
