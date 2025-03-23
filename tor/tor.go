package tor

import (
	"context"
	"fmt"
	"time"
	"net/http"
	"sync"
	"github.com/cretz/bine/tor"
	"golang.org/x/net/html"
	"bytes"
	"strings"
	"os"
)

func RunTor(port int, ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Starting proxy on port", port)
	stopChan := make(chan struct{})
	go func() {
		fmt.Println("Starting tor and fetching title of https://check.torproject.org, please wait a few seconds...")

		// Create a temporary torrc file with the specified SOCKS port
		torrcContent := fmt.Sprintf("SOCKSPort %d\n", port)
		torrcFile, err := os.CreateTemp("", "torrc-*.conf")
		if err != nil {
			fmt.Println("Error creating torrc file:", err)
			return
		}
		defer os.Remove(torrcFile.Name()) // Clean up the temporary file

		if _, err := torrcFile.WriteString(torrcContent); err != nil {
			fmt.Println("Error writing to torrc file:", err)
			return
		}
		if err := torrcFile.Close(); err != nil {
			fmt.Println("Error closing torrc file:", err)
			return
		}

		// Configure Tor to use the custom torrc file
		conf := &tor.StartConf{
			TorrcFile: torrcFile.Name(),
		}

		t, err := tor.Start(nil, conf)
		if err != nil {
			fmt.Println("Error starting Tor:", err)
			return
		}
		defer t.Close()

		// Wait at most a minute to start network and get
		dialCtx, dialCancel := context.WithTimeout(context.Background(), time.Minute)
		defer dialCancel()

		// Make connection
		dialer, err := t.Dialer(dialCtx, nil)
		if err != nil {
			fmt.Println("Error creating dialer:", err)
			return
		}

		httpClient := &http.Client{Transport: &http.Transport{DialContext: dialer.DialContext}}

		// Get /
		resp, err := httpClient.Get("https://check.torproject.org")
		if err != nil {
			fmt.Println("Error fetching URL:", err)
			return
		}
		defer resp.Body.Close()

		// Grab the <title>
		parsed, err := html.Parse(resp.Body)
		if err != nil {
			fmt.Println("Error parsing HTML:", err)
			return
		}

		title := getTitle(parsed)
		fmt.Printf("Title: %s\n", title)
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

func getTitle(n *html.Node) string {
	if n.Type == html.ElementNode && n.Data == "title" {
		var title bytes.Buffer
		if err := html.Render(&title, n.FirstChild); err != nil {
			panic(err)
		}
		return strings.TrimSpace(title.String())
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if title := getTitle(c); title != "" {
			return title
		}
	}
	return ""
}

