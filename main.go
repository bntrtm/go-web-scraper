package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"sync"
)

type config struct {
	pages              map[string]PageData
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
}

func main() {
	if len(os.Args[1:]) < 1 {
		fmt.Println("no website provided")
		os.Exit(1)
	} else if len(os.Args[1:]) > 1 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}
	rawBaseURL := os.Args[1]
	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		log.Fatal(err)
	}
	cfg := &config{
		baseURL:            baseURL,
		pages:              map[string]PageData{},
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, 5),
		wg:                 &sync.WaitGroup{},
	}

	fmt.Printf("starting crawl of: %s\n", cfg.baseURL)
	cfg.wg.Go(func() {
		defer func() { <-cfg.concurrencyControl }()
		cfg.concurrencyControl <- struct{}{}
		cfg.crawlPage(rawBaseURL)
	})

	cfg.wg.Wait()
	fmt.Println("\nLinks found per webpage:")
	fmt.Println("____________________")
	for k, v := range cfg.pages {
		fmt.Printf("%s: %d\n", k, len(v.OutgoingLinks)+len(v.ImageURLs))
	}
}
