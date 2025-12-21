// oreon/defense · watchthelight <wtl>

package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

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

	fmt.Printf("loaded config from %s\n", configPath)
	fmt.Printf("  real-time protection: %v\n", cfg.General.RealTimeProtection)
	fmt.Printf("  firewall integration: %v\n", cfg.Firewall.Enabled)

	// TODO: Initialize logging
	// TODO: Initialize state machine
	// TODO: Start IPC server
	// TODO: Start health check loop

	fmt.Println("daemon ready, waiting for shutdown...")

	<-ctx.Done()

	fmt.Println("shutting down...")
	return nil
}
