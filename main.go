package main

import (
//	"muxator/proxy"
	"muxator/tor"
    "sync"
    "context"
    "os"
    "os/signal"
    "syscall"
    "fmt"
)

func main() {
    // var ports = []int{3000,3001,3002,3003,3004}
    var ports = []int{9051}
    var wg sync.WaitGroup
    ctx, cancel := context.WithCancel(context.Background())
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan // Wait for signal
		fmt.Println("\nReceived shutdown signal, stopping all listeners...")
		cancel()  // Send cancellation signal to all goroutines
	}()
    for _, id := range ports {
        wg.Add(1)
        port := id
//        go proxy.Runproxy(port, ctx, &wg)
        go tor.RunTor(port, ctx, &wg)

    }
    wg.Wait()
}

// https://stackoverflow.com/questions/40328025/tcp-connection-over-tor-in-golang
// https://jonathanmh.com/p/golang-proxy-http-requests-via-tor/
