package crawler

import (
	"bufio"
	"net/http"
	"net/url"

	"github.com/bodokaiser/crawler/scan/html"
)

type Request struct {
	Done   chan bool
	Origin *url.URL
	Refers []*url.URL
}

func NewRequest(urlStr string) (*Request, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	return &Request{
		Done:   make(chan bool),
		Origin: u,
		Refers: make([]*url.URL, 0),
	}, nil
}

func (r *Request) Do() {
	res, err := http.Get(r.Origin.String())
	if err != nil {
		return
	}
	defer res.Body.Close()

	s := bufio.NewScanner(res.Body)
	s.Split(html.ScanHref)

	for s.Scan() {
		uri, err := url.Parse(s.Text())

		if err == nil && !r.Has(uri) {
			r.Push(uri)
		}
	}

	close(r.Done)
}

func (r *Request) Has(u *url.URL) bool {
	for _, ref := range r.Refers {
		if ref.String() == u.String() {
			return true
		}
	}

	return false
}

func (r *Request) Push(u *url.URL) {
	r.Refers = append(r.Refers, u)
}
