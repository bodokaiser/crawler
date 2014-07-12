package httpd

import "net/http"

var (
	// Amount of requests which can be queued without
	// getting a runtime error.
	MaxPoolQueue = 1000

	// Amount of go routines running in parallel to
	// execute requests from the queue.
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

	active int

	results chan *result
	queue   chan *http.Request
}

// Returns an initialized Pool with defaults.
func NewPool() *Pool {
	return &Pool{
		Client: http.DefaultClient,

		active: 0,

		results: make(chan *result, MaxPoolQueue),
		queue:   make(chan *http.Request, MaxPoolQueue),
	}
}

// Get will block until a pending request was made which then is
// returned with the corresponding response and an optional error.
// If there are no pending requests it will return you nil values.
func (p *Pool) Get() (*http.Request, *http.Response, error) {
	for {
		select {
		case r := <-p.results:
			return r.Request, r.Response, r.Error
		default:
			if p.active < 1 {
				return nil, nil, nil
			}
		}
	}
}

// Do adds a request to the pools queue.
func (p *Pool) Add(r *http.Request) {
	p.active++
	p.queue <- r
}

// Run spawns some workers which will process requests in queue.
func (p *Pool) Run() {
	for i := 0; i < MaxPoolWorker; i++ {
		go func() {
			for req := range p.queue {
				res, err := p.Client.Do(req)

				p.results <- &result{err, req, res}

				p.active--
			}
		}()
	}
}

type result struct {
	Error    error
	Request  *http.Request
	Response *http.Response
}
