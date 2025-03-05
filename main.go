package main

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("need two positional args: url and interval in seconds")
	}
	urlString := os.Args[1]
	curURL, err := url.Parse(urlString)
	if err != nil {
		log.Fatal("failed prepare request")
	}
	if curURL.Scheme == "" {
		curURL.Scheme = "https"
	}
	sec, err := strconv.ParseInt(os.Args[2], 10, 64)
	if err != nil {
		log.Fatal("failed get ticker from arg", err)
	}
	ticker := time.NewTicker(time.Duration(sec) * time.Second)
	done := make(chan bool)
	defaultTransport, ok := http.DefaultTransport.(*http.Transport)
	if !ok {
		log.Fatal("transport assertion error")
	}

	transport := defaultTransport.Clone()

	client := &http.Client{
		Transport: transport,
		Timeout:   3 * time.Second,
	}
	signalCtx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	req, err := http.NewRequestWithContext(signalCtx, http.MethodGet, curURL.String(), nil)
	if err != nil {
		log.Fatal("failed prepare request")
	}
	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				start := time.Now()
				curReq := req.Clone(context.Background())
				resp, err := client.Do(curReq)
				if err != nil {
					log.Printf("err: %v\n", err)
				}
				if resp.StatusCode != http.StatusOK {
					log.Println("status code not 200", resp.StatusCode)
				}
				log.Printf("response time: %s", time.Since(start))
			}
		}
	}()
	<-signalCtx.Done()
	done <- true
}
