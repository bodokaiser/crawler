package gerenuk

import (
	"bufio"
	"net/http"
	"strings"

	"github.com/bodokaiser/gerenuk/pipe"
	"github.com/bodokaiser/gerenuk/scan/html"
	"github.com/bodokaiser/gerenuk/work"
)

type Crawler struct {
	Worker *work.Worker
	Pipe   *pipe.Pipeline
}

func NewCrawler() *Crawler {
	return &Crawler{
		Worker: work.NewWorker(),
		Pipe:   pipe.NewPipeline(),
	}
}

func (c *Crawler) Put(p Page) {
	w := &Crawl{
		Page: p,
		Done: make(chan bool),
	}

	go func(w *Crawl, p *pipe.Pipeline) {
		<-w.Done

		if w.Error == nil {
			p.Emit(w.Page)
		}
	}(w, c.Pipe)

	c.Worker.Add(w)
}

type Crawl struct {
	Page  Page
	Done  chan bool
	Error error
}

func (c *Crawl) Do() {
	req, err := http.NewRequest("GET", c.Page.Origin(), nil)
	if err != nil {
		return
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return
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
			if !c.Page.HasRefer(t) {
				c.Page.AddRefer(t)
			}
		}
	}

	close(c.Done)
}
