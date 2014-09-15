package crawler

import "github.com/bodokaiser/crawler/work"

// Crawler manages parallel execution of Requests collecting the results and
// ensuring that there are no Request doublications.
type Crawler struct {
	visits map[string]*Request
	worker *work.Worker
}

// Returns a new initialized Crawler.
func New() *Crawler {
	return &Crawler{
		worker: work.New(),
		visits: make(map[string]*Request),
	}
}

// Does a crawl request and saves it as vistited. If request with same
// origin was already done it will be ignored.
func (c *Crawler) Do(r *Request) {
	u := r.Origin.String()

	if v, ok := c.visits[u]; ok {
		r.Refers = v.Refers

		return
	}

	c.visits[u] = r
	c.worker.Do(r)
}

// Executes crawl requests in worker queue.
func (c *Crawler) Run(n int) {
	c.worker.Run(n)
}

func (c *Crawler) Kill() {
	c.worker.Kill()
}
