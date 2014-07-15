package httpd

import (
	"net/http"
	"sync/atomic"
)

var (
	// Amount of requests which can be pendingd without
	// getting a runtime error.
	MaxPoolQueue = 1000

	// Amount of go routines running in parallel to
	// execute requests from the pending.
	MaxPoolWorker = 20
)

// Pool can take HTTP requests and execute them in parallel.
// By using Pool you avoid exponential go routine spawns which
// can crash your process also you get improved error handling
// and control of your capacities.
type Pool struct {
	// Client defines the default HTTP Client to use during
	// requests. It defaults to http.DefaultClient.
	Client *http.Client

	active  int32
	results chan *result
	pending chan *http.Request
}

// Returns an initialized Pool with defaults.
func NewPool() *Pool {
	return &Pool{
		Client:  http.DefaultClient,
		active:  0,
		results: make(chan *result, MaxPoolQueue),
		pending: make(chan *http.Request, MaxPoolQueue),
	}
}

// Get will block until a pending request was made which then is
// returned with the corresponding response and an optional error.
// If there are no pending requests it will return you nil values.
func (p *Pool) Get() (*http.Request, *http.Response, error) {
	for {
		select {
		case r := <-p.results:
			atomic.AddInt32(&p.active, -1)

			return r.Request, r.Response, r.Error
		default:
			if atomic.LoadInt32(&p.active) == 0 {
				return nil, nil, nil
			}
		}
	}
}

// Do adds a request to the pools pending.
func (p *Pool) Add(r *http.Request) {
	atomic.AddInt32(&p.active, 1)

	p.pending <- r
}

// Run spawns some workers which will process requests in pending.
func (p *Pool) Run() {
	for i := 0; i < MaxPoolWorker; i++ {
		go func() {
			for req := range p.pending {
				res, err := p.Client.Do(req)

				p.results <- &result{err, req, res}
			}
		}()
	}
}

type result struct {
	Error    error
	Request  *http.Request
	Response *http.Response
}
