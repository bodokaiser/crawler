package main

import (
	"flag"
	"fmt"
	"log"

	. "github.com/bodokaiser/gerenuk"
)

var conf = Config{}

func main() {
	flag.StringVar(&conf.DB, "db", "", "URL to mysql database.")
	flag.StringVar(&conf.Url, "url", "", "URL to use as crawler entry.")
	flag.Parse()

	c, err := NewCrawler(conf)
	if err != nil {
		log.Fatalf("Error initializing crawler: %s.\n", err)

		return
	}

	for {
		fmt.Printf("Visited: %s\n", c.Get())
	}
}
