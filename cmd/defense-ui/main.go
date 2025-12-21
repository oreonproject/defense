// oreon/defense · watchthelight <wtl>

package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/oreonproject/defense/internal/tray"
	"github.com/oreonproject/defense/pkg/ipc"
)

var version = "0.1.0-dev"

func main() {
	fmt.Printf("Oreon Defense v%s\n", version)

	// Create a channel to listen for interrupt signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	// Initialize the IPC client
	client, err := ipc.NewClient("/run/oreon/defense.sock")
	if err != nil {
		log.Fatalf("Failed to create IPC client: %v", err)
	}

	// Create and run the system tray
	trayApp := tray.New(client)

	// Run the tray in a goroutine so we can handle shutdown gracefully
	go func() {
		if err := trayApp.Run(); err != nil {
			log.Fatalf("Tray application error: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-sigCh
	log.Println("Shutting down...")
	// Any cleanup can be done here if needed

}
