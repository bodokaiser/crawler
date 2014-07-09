package robots

import (
	"net/http"

	"github.com/bodokaiser/gerenuk/parser"
)

type Robot struct {
	parsers []parser.Parser
}

func NewRobot() *Robot {
	return &Robot{}
}

func (r *Robot) Open(url string) error {
	res, err := http.Get(url)

	if err != nil {
		return err
	}

	for i := 0; i < len(r.parsers); i++ {
		r.parsers[i](res)
	}

	return nil
}

func (r *Robot) Register(p parser.Parser) {
	r.parsers = append(r.parsers, p)
}
