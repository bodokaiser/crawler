package parser

import "bytes"

var (
	URL_OPEN  = []byte("href=\"/")
	URL_CLOSE = []byte("\"")
)

type URLParser struct {
	results []string
}

func (p *URLParser) Write(b []byte) (int, error) {
	size := len(b)

	for i := 0; i < size-len(URL_OPEN); i++ {
		offset := i + len(URL_OPEN)

		// check if the current bytes equal the open sequence
		if bytes.Equal(b[i:offset], URL_OPEN) {
			for n := offset + 1; n < size; n++ {
				if b[n] == URL_CLOSE[0] {
					p.results = append(p.results, string(b[offset-1:n]))

					break
				}
			}
		}
	}

	return size, nil
}

func (p *URLParser) Result() []string {
	return p.results
}
