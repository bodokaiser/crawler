package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	. "github.com/bodokaiser/gerenuk"
	"github.com/bodokaiser/gerenuk/store"
	_ "github.com/bodokaiser/gerenuk/store/mysql"
)

var url string

var conf Config

func main() {
	flag.StringVar(&url, "url", "", "url to crawl")
	flag.StringVar(&conf.DB, "db", "", "url to database")
	flag.Parse()

	if len(url) == 0 {
		log.Fatalf("Please provide an url parameter.\n")
		return
	}

	c, err := NewCrawler(conf)
	if err != nil {
		log.Fatalf("Error initializing crawler: %s.\n", err)
		return
	}
	c.Put(store.Page{
		Origin: url,
		Refers: make([]string, 0),
	})

	for r := range c.Results() {
		fmt.Printf("%s: %s\n", r.Origin, strings.Join(r.Refers, ", "))

		c.Put(r)
	}
}
