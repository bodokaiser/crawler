package main

import (
	"flag"
	"fmt"
	"log"

	. "github.com/bodokaiser/gerenuk"
)

var db, url string

func main() {
	flag.StringVar(&db, "db", "", "URL to mysql database.")
	flag.StringVar(&url, "url", "", "URL to use as crawler entry.")
	flag.Parse()

	c, err := NewCrawler(db)
	if err != nil {
		log.Fatalf("Error initializing crawler: %s.\n", err)

		return
	}
	if err := c.Put(url); err != nil {
		log.Fatalf("Error putting url to store: %s.\n", err)

		return
	}

	for {
		fmt.Printf("Visited: %s\n", c.Get())
	}
}
