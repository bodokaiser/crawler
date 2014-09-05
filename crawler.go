package gerenuk

import (
	"bufio"
	"errors"
	"net/http"
	"strings"

	"github.com/bodokaiser/gerenuk/list"
	"github.com/bodokaiser/gerenuk/parser/html"
	"github.com/bodokaiser/gerenuk/pool"
)

var (
	ErrBadInput  = errors.New("bad input argument")
	ErrBadOutput = errors.New("bad output argument")
)

type Crawler struct {
	list *list.List
	pool *pool.Pool
}

func NewCrawler(uri string) *Crawler {
	c := &Crawler{
		list: list.NewList(),
		pool: pool.NewPool(worker),
	}
	c.put(uri)

	return c
}

func (c *Crawler) put(uri string) {
	i, err := list.NewItemFromUrl(uri)

	if err != nil {
		return
	}

	if !c.list.Has(i) {
		c.list.Add(i)
		c.pool.Put(i)
	}
}

func (c *Crawler) Get() (string, error) {
	res := c.pool.Get()

	if res == nil {
		return "", nil
	}

	switch t := res.(type) {
	case error:
		return "", t
	case *list.Item:
		for _, ref := range t.Refers() {
			c.put(ref)
		}

		return t.Origin(), nil
	}

	return "", ErrBadOutput
}

func worker(v interface{}) interface{} {
	if i, ok := v.(*list.Item); ok {
		req, err := http.NewRequest("GET", i.Origin(), nil)
		if err != nil {
			return err
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}

		defer res.Body.Close()

		s := bufio.NewScanner(res.Body)
		s.Split(html.ScanHref)

		for s.Scan() {
			t := s.Text()

			switch {
			case strings.HasPrefix(t, "/"):
				req.URL.Path = t
				t = req.URL.String()

				fallthrough
			case strings.HasPrefix(t, "http"):
				i.AddRefer(t)
			}
		}

		return i
	}

	return ErrBadInput
}
