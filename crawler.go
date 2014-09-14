package crawler

import (
	"time"

	"github.com/bodokaiser/crawler/work"
)

// Crawler manages parallel execution of Requests collecting the results and
// ensuring that there are no Request doublications.
type Crawler struct {
	worker   *work.Worker
	visited  []*Request
	results  chan *Request
	results2 chan []*Request
}

// Returns a new initialized Crawler.
func New() *Crawler {
	c := &Crawler{
		worker:   work.New(),
		visited:  make([]*Request, 0),
		results:  make(chan *Request),
		results2: make(chan []*Request),
	}

	go func(in <-chan *Request, out chan<- []*Request) {
		rs, timer := make([]*Request, 0), time.NewTimer(0)

		var t <-chan time.Time
		var o chan<- []*Request

		for {
			select {
			case r := <-in:
				rs = append(rs, r)

				if t == nil {
					timer.Reset(100 * time.Millisecond)
					t = timer.C
				}
			case <-t:
				o = out
				t = nil
			case o <- rs:
				rs = make([]*Request, 0)
				o = nil
			}
		}
	}(c.results, c.results2)

	return c
}

// Adds a Request to queue if not visited.
// In the background it will wait for the request to be done to send it to a
// result queue.
func (c *Crawler) Add(r *Request) {
	if !c.has(r) {
		go func(r *Request, out chan *Request) {
			<-r.Done

			out <- r
		}(r, c.results)

		c.push(r)
		c.worker.Add(r)
	}
}

// Returns a slice of ready to read Requests.
// As a single consumer is too slow to read the complete data we are buffering
// up Requests in slices when direct passthrough is not possible.
func (c *Crawler) Get() []*Request {
	return <-c.results2
}

func (c *Crawler) has(r *Request) bool {
	for _, v := range c.visited {
		if v.Origin.String() == r.Origin.String() {
			return true
		}
	}

	return false
}

func (c *Crawler) push(r *Request) {
	c.visited = append(c.visited, r)
}
