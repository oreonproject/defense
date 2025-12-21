// oreon/defense · watchthelight <wtl>

package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	sddaemon "github.com/coreos/go-systemd/v22/daemon"
	"github.com/oreonproject/defense/internal/daemon"
	"github.com/oreonproject/defense/pkg/config"
	"github.com/oreonproject/defense/pkg/logging"
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
	// Load config
	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	fmt.Printf("loaded config from %s\n", configPath)
	fmt.Printf("  real-time protection: %v\n", cfg.General.RealTimeProtection)
	fmt.Printf("  firewall integration: %v\n", cfg.Firewall.Enabled)

	// Initialize logger
	logCfg := logging.Config{
		Level:       cfg.General.LogLevel,
		FilePath:    config.LogPath,
		UseJournald: true, // detect systemd and use journald if available
	}
	logger, cleanup, err := logging.New(logCfg)
	if err != nil {
		return fmt.Errorf("initializing logger: %w", err)
	}
	defer cleanup()

	logger.Info("oreon defense starting", "version", version)

	// Create and start daemon
	d := daemon.New(cfg, logger)

	// Register state change listener for logging
	d.StateManager().OnStateChange(func(old, new daemon.State) {
		logger.Info("state changed", "from", old, "to", new)
	})

	// TODO: Start IPC server

	fmt.Println("daemon ready, waiting for shutdown...")

	// notify systemd that we're ready
	if sent, err := sddaemon.SdNotify(false, sddaemon.SdNotifyReady); err != nil {
		logger.Warn("failed to notify systemd", "error", err)
	} else if sent {
		logger.Info("systemd notification sent")
	}

	// Run daemon (blocks until context is cancelled)
	if err := d.Run(ctx); err != nil {
		return fmt.Errorf("daemon error: %w", err)
	}

	logger.Info("daemon shutdown complete")

	// notify systemd that we're stopping
	sddaemon.SdNotify(false, sddaemon.SdNotifyStopping)

	return nil
}
