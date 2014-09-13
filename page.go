package crawler

import (
	"net/url"
	"strings"
)

type Page struct {
	origin string
	refers []string
}

func NewPage(urlStr string) *Page {
	return &Page{
		origin: urlStr,
		refers: make([]string, 0),
	}
}

func (p *Page) Origin() string {
	return p.origin
}

func (p *Page) SetOrigin(url string) {
	url = normalize(url).String()

	if len(url) > 0 {
		p.origin = url
	}
}

func (p *Page) Refers() []string {
	return p.refers
}

func (p *Page) HasRefer(url string) bool {
	url = normalize(url).String()

	for _, ref := range p.refers {
		if ref == url {
			return true
		}
	}

	return false
}

func (p *Page) AddRefer(url string) {
	switch {
	case strings.HasPrefix(url, "/"):
		o := normalize(p.origin)
		o.Path = url
		url = o.String()

		fallthrough
	case strings.HasPrefix(url, "http"):
		p.refers = append(p.refers, url)
	}
}

func normalize(urlStr string) *url.URL {
	uri, err := url.Parse(urlStr)
	if err != nil {
		return nil
	}

	uri.Host = strings.ToLower(uri.Host)
	uri.Scheme = strings.ToLower(uri.Scheme)
	uri.Fragment = ""

	return uri
}
