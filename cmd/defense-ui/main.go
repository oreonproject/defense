// oreon/defense · watchthelight <wtl>

package main

import (
	"fmt"
	"os"
)

var version = "0.1.0-dev"

func main() {
	fmt.Printf("Oreon Defense v%s\n", version)

	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	fmt.Println("Qt bindings not yet configured")
	return nil
}
