package work

import (
	"sync/atomic"
	"time"
)

// Default parameters to be used on new Worker.
var (
	DefaultTimeout         = time.Second
	DefaultMaxQueue  int64 = 200
	DefaultMaxWorker int64 = 100
)

// Type Work is an interface for tasks which can be added to a Worker.
type Work interface {
	// Do will be executed by the Worker.
	// Internally you may use channels to communicate results and errors.
	Do()
}

// Type Worker executes Work concurrently in goroutines. The more work is added
// the more goroutines are spawned until max worker amount is reached. Worker
// will timeout after max timeout to free resources.
type Worker struct {
	queue   chan Work
	worker  int64
	active  int64
	pending int64
	timeout time.Duration
}

// Initializes a new Worker with default parameters.
func New() *Worker {
	w := &Worker{
		queue: make(chan Work, DefaultMaxQueue),
	}
	w.SetTimeout(DefaultTimeout)
	w.SetMaxWorker(DefaultMaxWorker)

	return w
}

// Adds type which implements the Work interface to queue and spawns a
// goroutines for it if max capacity not reached.
func (w *Worker) Add(work Work) {
	atomic.AddInt64(&w.pending, 1)

	// TODO: Use of atomic counters is vulnerable for race conditions as values
	// may change while execution. As the worst case would be more goroutines to
	// be spawned and there does not seem to be an easy workaround - mutex will
	// block - we take this risk for now.
	if a, p := atomic.LoadInt64(&w.active), atomic.LoadInt64(&w.pending); p > a && a < w.worker {
		go func(queue <-chan Work, timeout time.Duration) {
			for {
				select {
				case work := <-queue:
					work.Do()

					atomic.AddInt64(&w.pending, -1)
				case <-time.After(timeout):
					atomic.AddInt64(&w.active, -1)

					return
				}
			}
		}(w.queue, w.timeout)

		atomic.AddInt64(&w.active, -1)
	}

	w.queue <- work
}

// Updates max timeout to duration.
// Only affects goroutines which are spawned post mortern.
func (w *Worker) SetTimeout(d time.Duration) {
	w.timeout = d
}

// Updates maximum amount of parallel running worker.
// Only affects policy for worker spawned post mortern.
func (w *Worker) SetMaxWorker(n int64) {
	w.worker = n
}
