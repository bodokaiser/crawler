package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/bodokaiser/crawler"
)

var wg = new(sync.WaitGroup)

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

	for i := 0; i < 1000; i++ {
		wg.Add(1)

		go crawl(c)
	}

	wg.Wait()
}

func crawl(c *crawler.Crawler) {
	for {
		p := c.Get()

		fmt.Printf("%s\n", p.Origin)

		for _, r := range p.Refers {
			p := &crawler.Page{}
			p.SetOrigin(r)

			c.Put(p)
		}
	}

	wg.Done()
}
