package html

import (
	"bufio"
	"bytes"
	"io"

	"github.com/bodokaiser/gerenuk/parser"
)

var (
	hrefPrefix = []byte(`href="`)
	hrefSuffix = []byte(`"`)
)

type HrefParser struct {
	scanner *bufio.Scanner
}

func NewHrefParser(r io.Reader) *HrefParser {
	s := bufio.NewScanner(r)
	s.Split(ScanHref)

	return &HrefParser{s}
}

func (p *HrefParser) Next() *parser.Result {
	s := p.scanner

	for s.Scan() {
		b := s.Bytes()

		return &parser.Result{b}
	}

	return nil
}

func ScanHref(b []byte, eof bool) (int, []byte, error) {
	i := bytes.Index(b, hrefPrefix)

	if i != -1 {
		i += len(hrefPrefix)
		n := bytes.Index(b[i:], hrefSuffix) + i

		// check if slice is in range
		if n != -1 && n > i && n < len(b) {
			return n + 1, b[i:n], nil
		}
	}

	return len(b), nil, nil
}
