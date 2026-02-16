package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
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
	cFlag := flag.Int("c", 1, "Number of concurrent goroutines to use")
	pFlag := flag.Int("p", 3, "Maximum number of pages to crawl before exit")
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("no website provided")
		os.Exit(1)
	} else if len(args) > 1 {
		fmt.Println("too many arguments provided")
		fmt.Println(args)
		os.Exit(1)
	}
	rawBaseURL := args[0]
	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		log.Fatal(err)
	}
	cfg := &config{
		baseURL:            baseURL,
		pages:              map[string]PageData{},
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, *cFlag),
		wg:                 &sync.WaitGroup{},
		maxPages:           *pFlag,
	}

	fmt.Printf("starting crawl of: %s\n", cfg.baseURL)
	cfg.wg.Go(func() {
		cfg.crawlPage(rawBaseURL)
	})

	cfg.wg.Wait()
	if len(cfg.pages) > 0 {
		filename := fmt.Sprintf("%s_report.csv", strings.ReplaceAll(cfg.baseURL.Hostname(), ".", "-"))
		err := writeCSVReport(cfg.pages, filename)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("generated report: %s\n", filename)
	} else {
		fmt.Println("No links found on webpage given.")
	}
}
