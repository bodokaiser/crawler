package main

import (
	"fmt"
	"log"
	"strings"

	. "github.com/bodokaiser/gerenuk"
	_ "github.com/bodokaiser/gerenuk/store/mysql"
)

var conf = &Config{}

func main() {
	conf.Parse()

	if err := conf.Check(); err != nil {
		log.Fatalf("Error parsing parameters: %s.\n", err)

		return
	}

	p := Page{}
	p.SetOrigin(conf.Origin)

	c, err := NewCrawler(conf)
	if err != nil {
		log.Fatalf("Error initializing crawler: %s.\n", err)

		return
	}
	c.Put(p)

	for r := range c.Results() {
		fmt.Printf("%s: %s\n", r.Origin(), strings.Join(r.Refers(), ", "))

		c.Put(r)
	}
}
