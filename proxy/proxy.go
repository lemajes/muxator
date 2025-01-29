package proxy

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/things-go/go-socks5"
)
func Runproxy(port int, ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Println("Starting proxy on port", port)

	// Create the SOCKS5 server
	server := socks5.NewServer(
		socks5.WithLogger(socks5.NewLogger(log.New(os.Stdout, "socks5: ", log.LstdFlags))),
	)

	thisPort := fmt.Sprintf(":%d", port)

	// Use a channel to stop the server
	stopChan := make(chan struct{})

	// Start the server in a separate goroutine
	go func() {
		if err := server.ListenAndServe("tcp", thisPort); err != nil {
			log.Printf("Error starting SOCKS5 server on port %d: %v", port, err)
		}
		close(stopChan) // Ensure the stop signal is sent if the server exits
	}()

	// Wait for the shutdown signal
	select {
	case <-ctx.Done(): // Shutdown signal received
		fmt.Printf("Proxy %d received stop signal. Shutting down...\n", port)
	case <-stopChan: // Server exited unexpectedly
		fmt.Printf("Proxy %d stopped unexpectedly.\n", port)
	}
}
