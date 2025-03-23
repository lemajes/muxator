package tor

import (
	"context"
	"fmt"
	"log"
	"sync"

)
func RunTor(port int, ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Println("Starting proxy on port", port)
	stopChan := make(chan struct{})
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
