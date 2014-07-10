package main

import (
	"fmt"
	"log"

	"github.com/bodokaiser/gerenuk/robot"
	"github.com/bodokaiser/gerenuk/split/html"
)

func main() {
	r := robot.NewRobot(html.SplitHref, html.SplitURL)

	if err := r.Open("http://www.satisfeet.me"); err != nil {
		log.Fatal(err)
	}

	for result := range r.Results {
		fmt.Printf("result: %s\n", result)
	}
}
