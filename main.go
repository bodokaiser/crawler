package main

import (
	"fmt"
	"log"

	"github.com/bodokaiser/go-crawler/conf"
	"github.com/bodokaiser/go-crawler/robot"
)

func main() {
	c := conf.New()

	if err := c.Flags(); err != nil {
		log.Fatal(err)
	}

	r := robot.New(c)

	if err := r.Open(); err != nil {
		log.Fatal(err)
	}

	for i := range r.Results {
		fmt.Printf("Result: %s\n", r.Results[i])
	}
}
