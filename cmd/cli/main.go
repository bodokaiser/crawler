package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/bodokaiser/crawler"
)

func main() {
	var entry string

	flag.StringVar(&entry, "entry", "", "")
	flag.Parse()

	r, err := crawler.NewRequest(entry)
	if err != nil {
		log.Fatalf("Error on initial request: %s.\n", err)

		return
	}

	c := crawler.New()
	c.Add(r)

	for {
		for _, r := range c.Get() {
			fmt.Printf("%s\n", r.Origin)

			for _, u := range r.Refers {
				r, _ := crawler.NewRequest(u.String())

				c.Add(r)
			}
		}
	}
}
