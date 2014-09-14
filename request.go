package crawler

import (
	"bufio"
	"net/http"
	"net/url"
	"strings"

	"github.com/bodokaiser/crawler/scan/html"
)

// Requests represents a crawl request to a HTTP server which sits behind the
// Origin url. It supports parallel execution by implementing the work.Work
// interface where it will do a HTTP request and push extracted hrefs from the
// HTML code into the Refers slice.
type Request struct {
	Done   chan bool
	Origin *url.URL
	Refers []*url.URL
}

// Returns a new request which will be done against the provided url string.
// Returns an error if url package fails to parse url string.
func NewRequest(urlStr string) (*Request, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	normalize(u)

	return &Request{
		Done:   make(chan bool),
		Origin: u,
		Refers: make([]*url.URL, 0),
	}, nil
}

// Implemention of work.Work interface to execute the request and extract links.
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

// Returns true if request has saved provided url as reference. To be consistent
// the provided url will be normalize before.
func (r *Request) Has(u *url.URL) bool {
	normalize(u)

	for _, ref := range r.Refers {
		if ref.String() == u.String() {
			return true
		}
	}

	return false
}

// Appends the normalized url to the Refers slice.
func (r *Request) Push(u *url.URL) {
	normalize(u)

	r.Refers = append(r.Refers, u)
}

func normalize(u *url.URL) {
	u.Fragment = ""
	u.Host = strings.ToLower(u.Host)
	u.Scheme = strings.ToLower(u.Scheme)
}
