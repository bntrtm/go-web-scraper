package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args[1:]) < 1 {
		fmt.Println("no website provided")
		os.Exit(1)
	} else if len(os.Args[1:]) > 1 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}
	fmt.Printf("starting crawl of: %s\n", os.Args[1])
	pages := map[string]int{}
	crawlPage(os.Args[1], os.Args[1], pages)
	fmt.Println("HTML Extracted:")
	fmt.Println("____________________")
	for k, v := range pages {
		fmt.Printf("%s: %d\n", k, v)
	}
}
