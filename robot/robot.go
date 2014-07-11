package robot

import (
	"bufio"
	"bytes"
	"io"
	"net/http"

	"github.com/bodokaiser/gerenuk/parser"
)

type Robot struct {
	parsers []parser.New
}

type Handle func(parser.Result)

func NewRobot(p ...parser.New) *Robot {
	return &Robot{p}
}

func (r *Robot) Open(url string, h Handle) error {
	res, err := http.Get(url)

	if err != nil {
		return err
	}

	l := len(r.parsers)

	in := make([]chan []byte, l)
	out := make(chan parser.Result)

	for i := 0; i < l; i++ {
		in[i] = make(chan []byte)

		go spawn(r.parsers[i], in[i], out)
	}

	go read(res.Body, in)

	for r := range out {
		h(r)
	}

	return nil
}

func read(r io.Reader, in []chan []byte) {
	s := bufio.NewScanner(r)

	for s.Scan() {
		b := s.Bytes()

		for _, in := range in {
			in <- b
		}
	}
}

func spawn(p parser.New, in chan []byte, out chan parser.Result) {
	for b := range in {
		r := bytes.NewReader(b)

		parse(p(r), out)
	}

	close(out)
}

func parse(p parser.Parser, out chan parser.Result) {
	for {
		r := p.Next()

		if r == nil {
			break
		}

		out <- r
	}
}
