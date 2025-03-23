package tor

import (
	"context"
	"fmt"
    "time"
    "net/http"
	"sync"
	"github.com/cretz/bine/tor"
    "golang.org/x/net/html"
)
func RunTor(port int, ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Println("Starting proxy on port", port)
	stopChan := make(chan struct{})
	go func() error {
    	fmt.Println("Starting tor and fetching title of https://check.torproject.org, please wait a few seconds...")
    	t, err := tor.Start(nil, nil)
    	if err != nil {
    		return err
    	}
    	defer t.Close()
    	// Wait at most a minute to start network and get
    	dialCtx, dialCancel := context.WithTimeout(context.Background(), time.Minute)
    	defer dialCancel()
    	// Make connection
    	dialer, err := t.Dialer(dialCtx, nil)
    	if err != nil {
    		return err
    	}
    	httpClient := &http.Client{Transport: &http.Transport{DialContext: dialer.DialContext}}
    	// Get /
    	resp, err := httpClient.Get("https://check.torproject.org")
    	if err != nil {
    		return err
    	}
    	defer resp.Body.Close()
    	// Grab the <title>
    	parsed, err := html.Parse(resp.Body)
    	if err != nil {
    		return err
    	}
    	fmt.Printf("Title: %v\n", getTitle(parsed))
    	return nil
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
