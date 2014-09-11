package gerenuk

import (
	"bufio"
	"net/http"
	"strings"

	"github.com/bodokaiser/gerenuk/html"
	"github.com/bodokaiser/gerenuk/pipeline"
	"github.com/bodokaiser/gerenuk/queue"
	"github.com/bodokaiser/gerenuk/store"
	"github.com/bodokaiser/gerenuk/worker"
)

type Crawler struct {
	active bool
	store  store.Store
	queue  *queue.Queue
	worker *worker.Pool
	entry  *pipeline.Pipeline
	result *pipeline.Pipeline
}

func NewCrawler(c *Config) (*Crawler, error) {
	s, err := store.Open("mysql", c.Store.Url)
	if err != nil {
		return nil, err
	}

	q := queue.NewQueue()
	ep := pipeline.NewPipeline()
	rp := ep.Pipe(pipeline.StageFunc(func(in <-chan pipeline.Event) <-chan pipeline.Event {
		ch := make(chan pipeline.Event)

		go func(out chan<- pipeline.Event) {
			for e := range in {
				p := e.Result.(Page)

				for _, r := range p.Refers() {
					ch <- pipeline.Event{
						Result: Page{
							origin: r,
							refers: make([]string, 0),
						},
					}
				}
			}
		}(ch)

		return ch
	}))

	return &Crawler{
		store:  s,
		queue:  q,
		worker: worker.NewPool(),
		entry:  ep,
		result: rp,
	}, nil
}

func (c *Crawler) Put(p Page) {
	c.queue.Push(p)

	if !c.active {
		go c.poll()
	}
}

func (c *Crawler) Results() <-chan Page {
	ch := make(chan Page)

	go func() {
		for e := range c.result.Listen() {
			ch <- e.Result.(Page)
		}
		close(ch)
	}()

	return ch
}

func (c *Crawler) poll() {
	for {
		x := c.queue.Pull()

		if x != nil {
			p := x.(Page)

			w := &worker.Work{
				Func:   crawl,
				Done:   make(chan bool),
				Params: []interface{}{p},
			}
			go func(w *worker.Work, pipe *pipeline.Pipeline) {
				<-w.Done

				if w.Error == nil {
					pipe.Emit(pipeline.Event{
						Result: w.Result.(Page),
					})
				}
			}(w, c.entry)

			c.worker.Put(w)
		}
	}
}

func crawl(params ...interface{}) (interface{}, error) {
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
