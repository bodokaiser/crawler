package main

import (
	"fmt"
	"log"

	"github.com/bodokaiser/gerenuk/parser"
	"github.com/bodokaiser/gerenuk/parser/html"
	"github.com/bodokaiser/gerenuk/robot"
)

func main() {
	r := robot.NewRobot(html.NewHrefParser)

	if err := r.Open("http://www.satisfeet.me", handle); err != nil {
		log.Fatal(err)
	}
}

func handle(r parser.Result) {
	switch r := r.(type) {
	case *html.HrefResult:
		fmt.Printf("result: %s with macro: %s\n", r.Value(), r.Macro())
		break
	default:
		fmt.Printf("result: %s\n", r.Value())
	}
}
