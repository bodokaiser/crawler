package gerenuk

import (
	"bufio"
	"net/http"
	"strings"

	"github.com/bodokaiser/gerenuk/html"
	"github.com/bodokaiser/gerenuk/pipe"
	"github.com/bodokaiser/gerenuk/worker"
)

type Crawler struct {
	Pool *worker.Pool
	Pipe *pipe.Pipeline
}

func NewCrawler() *Crawler {
	return &Crawler{
		Pool: worker.NewPool(),
		Pipe: pipe.NewPipeline(),
	}
}

func (c *Crawler) Put(p Page) {
	w := &worker.Work{
		Func:   DefaultWork,
		Done:   make(chan bool),
		Params: []interface{}{p},
	}

	go func(w *worker.Work, p *pipe.Pipeline) {
		<-w.Done

		if w.Error == nil {
			p.Emit(w.Result.(Page))
		}
	}(w, c.Pipe)

	c.Pool.Put(w)
}

func DefaultWork(params ...interface{}) (interface{}, error) {
	p := params[0].(Page)

	req, err := http.NewRequest("GET", p.Origin(), nil)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
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
			if !p.HasRefer(t) {
				p.AddRefer(t)
			}
		}
	}

	return p, nil
}
