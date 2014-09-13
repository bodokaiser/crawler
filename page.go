package crawler

import (
	"net/url"
	"strings"
)

type Page struct {
	Origin string
	Refers []string
}

func (p *Page) SetOrigin(url string) {
	url = normalize(url).String()

	if len(url) > 0 {
		p.Origin = url
	}
}

func (p *Page) HasRefer(url string) bool {
	url = normalize(url).String()

	for _, ref := range p.Refers {
		if ref == url {
			return true
		}
	}

	return false
}

func (p *Page) AddRefer(url string) {
	switch {
	case strings.HasPrefix(url, "/"):
		o := normalize(p.Origin)
		o.Path = url
		url = o.String()

		fallthrough
	case strings.HasPrefix(url, "http"):
		p.Refers = append(p.Refers, url)
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
