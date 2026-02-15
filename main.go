package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"sync"
)

type config struct {
	pages              map[string]PageData
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
	maxPages           int
}

func main() {
	if len(os.Args[1:]) < 1 {
		fmt.Println("no website provided")
		os.Exit(1)
	} else if len(os.Args[1:]) > 3 {
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
		concurrencyControl: make(chan struct{}, 1),
		wg:                 &sync.WaitGroup{},
		maxPages:           3,
	}
	if len(os.Args[1:]) > 1 {
		i, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatalf("error: bad input for concurrency control: %s (must be integer)", os.Args[2])
		}
		cfg.concurrencyControl = make(chan struct{}, i)
	}
	if len(os.Args[1:]) > 2 {
		i, err := strconv.Atoi(os.Args[3])
		if err != nil {
			log.Fatalf("error: bad input for max page crawl count: %s (must be integer)", os.Args[3])
		}
		cfg.maxPages = i
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
