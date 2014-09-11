package gerenuk

import (
	"bufio"
	"net/http"
	"strings"

	"github.com/bodokaiser/gerenuk/parser/html"
	"github.com/bodokaiser/gerenuk/pipeline"
	"github.com/bodokaiser/gerenuk/queue"
	"github.com/bodokaiser/gerenuk/store"
	"github.com/bodokaiser/gerenuk/worker"
)

type Config struct {
	DB string
}

type Crawler struct {
	active bool
	store  store.Store
	queue  *queue.Queue
	worker *worker.WorkerPool
	entry  *pipeline.Pipeline
	result *pipeline.Pipeline
}

func NewCrawler(c Config) (*Crawler, error) {
	s, err := store.Open("mysql", c.DB)
	if err != nil {
		return nil, err
	}

	q := queue.NewQueue()
	ep := pipeline.NewPipeline()
	rp := ep.Pipe(pipeline.StageFunc(func(in <-chan pipeline.Event) <-chan pipeline.Event {
		ch := make(chan pipeline.Event)

		go func(out chan<- pipeline.Event) {
			for e := range in {
				p := e.Result.(store.Page)

				for _, r := range p.Refers {
					ch <- pipeline.Event{
						Result: store.Page{
							Origin: r,
							Refers: make([]string, 0),
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
		worker: worker.NewWorkerPool(),
		entry:  ep,
		result: rp,
	}, nil
}

func (c *Crawler) Put(p store.Page) {
	c.queue.Push(p)

	if !c.active {
		go c.poll()
	}
}

func (c *Crawler) Results() <-chan store.Page {
	ch := make(chan store.Page)

	go func() {
		for e := range c.result.Listen() {
			ch <- e.Result.(store.Page)
		}
		close(ch)
	}()

	return ch
}

func (c *Crawler) poll() {
	for {
		x := c.queue.Pull()

		if x != nil {
			p := x.(store.Page)

			w := &worker.Work{
				Func:   crawl,
				Done:   make(chan bool),
				Params: []interface{}{p},
			}
			go func(w *worker.Work, pipe *pipeline.Pipeline) {
				<-w.Done

				if w.Error == nil {
					pipe.Emit(pipeline.Event{
						Result: w.Result.(store.Page),
					})
				}
			}(w, c.entry)

			c.worker.Put(w)
		}
	}
}

func crawl(params ...interface{}) (interface{}, error) {
	p := params[0].(store.Page)

	req, err := http.NewRequest("GET", p.Origin, nil)
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
