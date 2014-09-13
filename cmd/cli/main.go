package main

import (
	"fmt"
	"log"

	"github.com/bodokaiser/crawler"
)

func main() {
	conf := &crawler.Config{}
	if err := conf.Parse(); err != nil {
		log.Fatalf("Error parsing parameters: %s.\n", err)

		return
	}

	c := crawler.New()
	c.Put(&crawler.Page{
		Origin: conf.Entry,
	})

	for {
		p := c.Get()

		fmt.Printf("%s\n", p.Origin)

		for _, r := range p.Refers {
			p := &crawler.Page{}
			p.SetOrigin(r)

			c.Put(p)
		}
	}
}
