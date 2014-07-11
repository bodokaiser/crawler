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
	hrefMacro  = []byte(`:`)
)

type HrefResult struct {
	macro []byte
	value []byte
}

func (r *HrefResult) Macro() string {
	return string(r.macro)
}

func (r *HrefResult) Value() string {
	return string(r.value)
}

type HrefParser struct {
	scanner *bufio.Scanner
}

func NewHrefParser(r io.Reader) parser.Parser {
	s := bufio.NewScanner(r)
	s.Split(ScanHref)

	return &HrefParser{s}
}

func (p *HrefParser) Next() parser.Result {
	s := p.scanner

	for s.Scan() {
		b := s.Bytes()

		var m []byte

		if i := bytes.Index(b, hrefMacro); i != -1 {
			m = b[:i]
			b = b[i+1:]
		}

		return &HrefResult{m, b}
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
