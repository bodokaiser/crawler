package gerenuk

import (
	"net/url"
	"strings"
)

type Page struct {
	origin string
	refers []string
}

func (p *Page) Origin() string {
	return p.origin
}

func (p *Page) SetOrigin(url string) {
	url = normalize(url)

	if len(url) > 0 {
		p.origin = url
	}
}

func (p *Page) Refers() []string {
	return p.refers
}

func (p *Page) HasRefer(url string) bool {
	url = normalize(url)

	for _, ref := range p.refers {
		if ref == url {
			return true
		}
	}

	return false
}

func (p *Page) AddRefer(url string) {
	if p.refers == nil {
		p.refers = make([]string, 0)
	}
	if url := normalize(url); len(url) > 0 {
		p.refers = append(p.refers, url)
	}
}

func normalize(urlStr string) string {
	uri, err := url.Parse(urlStr)
	if err != nil {
		return ""
	}

	uri.Host = strings.ToLower(uri.Host)
	uri.Scheme = strings.ToLower(uri.Scheme)
	uri.Fragment = ""

	return uri.String()
}
