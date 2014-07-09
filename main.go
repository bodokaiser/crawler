package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/bodokaiser/gerenuk/conf"
	"github.com/bodokaiser/gerenuk/parser"
	"github.com/bodokaiser/gerenuk/robots"
)

func main() {
	c := conf.New()

	if err := c.Flags(); err != nil {
		log.Fatal(err)
	}

	r := make(chan []string)

	go request(c["url"], r)

	for {
		result, ok := <-r

		if ok {
			fmt.Printf("\nResult: %s\n", strings.Join(result, ", "))
		} else {
			break
		}
	}
}

func request(url string, results chan []string) {
	r := robots.NewRobot(results)

	r.RegisterParser(&parser.URLParser{})
	r.RegisterParser(&parser.EmailParser{})

	if err := r.Open(url); err != nil {
		log.Fatal(err)
	}
}
