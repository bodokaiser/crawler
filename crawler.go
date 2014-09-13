package crawler

import (
	"bufio"
	"net/http"

	"github.com/bodokaiser/crawler/scan/html"
	"github.com/bodokaiser/crawler/work"
)

type Crawler struct {
	pages  *Pages
	result chan *Page
	worker *work.Worker
}

func New() *Crawler {
	return &Crawler{
		pages:  new(Pages),
		result: make(chan *Page),
		worker: work.New(),
	}
}

func (c *Crawler) Put(p *Page) {
	if c.pages.Has(p) {
		return
	}

	w := &crawl{p, make(chan bool)}

	go func(in <-chan bool, out chan<- *Page, p *Page) {
		<-in

		out <- p
	}(w.Done, c.result, w.Page)

	c.pages.Add(p)
	c.worker.Add(w)
}

func (c *Crawler) Get() *Page {
	return <-c.result
}

type crawl struct {
	Page *Page
	Done chan bool
}

func (c *crawl) Do() {
	res, err := http.Get(c.Page.Origin)
	if err != nil {
		return
	}
	defer res.Body.Close()

	s := bufio.NewScanner(res.Body)
	s.Split(html.ScanHref)

	for s.Scan() {
		t := s.Text()

		if !c.Page.Has(t) {
			c.Page.Push(t)
		}
	}

	close(c.Done)
}
