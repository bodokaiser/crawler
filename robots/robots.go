package robots

import (
	"io"
	"net/http"

	"github.com/bodokaiser/gerenuk/parser"
)

type Robot struct {
	results chan []string
	parsers []parser.Parser
}

func NewRobot(r chan []string) *Robot {
	p := []parser.Parser{}

	return &Robot{r, p}
}

func (r *Robot) Open(url string) error {
	res, err := http.Get(url)

	if err != nil {
		return err
	}

	for i := 0; i < len(r.parsers); i++ {
		io.Copy(r.parsers[i], res.Body)

		r.results <- r.parsers[i].Result()

		if i == len(r.parsers)-1 {
			close(r.results)
		}
	}

	return nil
}

func (r *Robot) RegisterParser(p parser.Parser) {
	r.parsers = append(r.parsers, p)
}
