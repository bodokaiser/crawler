package pool

import (
	"sync"
	"sync/atomic"
)

// Type Work defines a job which should be executed in the work pool.
type Work struct {
	// Channel which indicates if work was done.
	Done chan bool
	// Function to execute in pool.
	Func func(...interface{}) (interface{}, error)
	// Error worker returned.
	Error error
	// Result worker returned.
	Result interface{}
	// Arguments to pass to worker function.
	Params []interface{}
}

// Defines maximum amount of parallel worker routines.
var DefaultMaxWorker = 100

// Type Config defines settings for a WorkPool.
type Config struct {
	// New can be used a work factory.
	New func() *Work
	// Defines custom amount of maximum parallel workers.
	MaxWorker int
}

// Type Pool defines a facility to execute work in parallel.
type WorkPool struct {
	// Stores undone work.
	work *sync.Pool
	// Amount of max worker.
	worker int
	// Amount of active workers.
	active int32
	// Amount of pending work.
	pending int32
}

// Returns initialized WorkPool with settings from Config.
func NewWorkPool(c Config) *WorkPool {
	p := &WorkPool{
		work:   &sync.Pool{},
		worker: DefaultMaxWorker,
	}

	if c.MaxWorker != 0 {
		p.worker = c.MaxWorker
	}
	if c.New != nil {
		p.work.New = func() interface{} {
			return c.New()
		}

		p.Put(c.New())
	}

	return p
}

// Puts work into the worker pool.
// Will spawn more workers if capacity is open and go routines busy.
func (p *WorkPool) Put(w *Work) {
	atomic.AddInt32(&p.pending, 1)

	pen := atomic.LoadInt32(&p.pending)
	act := atomic.LoadInt32(&p.active)

	if pen > act && act < int32(p.worker) {
		go func(work *sync.Pool) {
			for {
				w, ok := work.Get().(*Work)

				if !ok || w == nil {
					atomic.AddInt32(&p.active, -1)

					return
				}

				w.Result, w.Error = w.Func(w.Params...)

				if w.Done != nil {
					w.Done <- true

					close(w.Done)
				}

				atomic.AddInt32(&p.pending, -1)
			}
		}(p.work)

		atomic.AddInt32(&p.active, 1)
	}

	p.work.Put(w)
}
