package main

import (
	"flag"
	"fmt"
	"log"

	. "github.com/bodokaiser/gerenuk"
)

var url string

func main() {
	flag.StringVar(&url, "url", "", "URL to crawl.")
	flag.Parse()

	c := NewCrawler(url)

	for {
		url, err := c.Get()
		if err != nil {
			log.Fatalf("Error working: %s\n", err)

			return
		}

		fmt.Printf("Visited: %s\n", url)
	}
}
