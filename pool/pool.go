package pool

import (
	"sync"
	"sync/atomic"
)

// Default maximum amount of concurrent go routines.
var DefaultMaxWorker = 100

// Type Work defines a task to be executed in WorkerPool.
type Work struct {
	Func   func(...interface{}) (interface{}, error)
	Done   chan bool
	Error  error
	Result interface{}
	Params []interface{}
}

// Type WorkerPool allows concurrent execution of tasks.
type WorkerPool struct {
	work    *sync.Pool
	worker  int
	active  int32
	pending int32
}

// Returns new initialized worker pool.
func NewWorkerPool() *WorkerPool {
	return &WorkerPool{
		work:   &sync.Pool{},
		worker: DefaultMaxWorker,
	}
}

// Puts work into worker pool and spawns worker if capacity is not full.
func (p *WorkerPool) Put(w *Work) {
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
				}

				atomic.AddInt32(&p.pending, -1)
			}
		}(p.work)

		atomic.AddInt32(&p.active, 1)
	}

	p.work.Put(w)
}

// Sets function to automatically generate new work.
func (p *WorkerPool) SetNewFunc(fn func() *Work) {
	p.work.New = func() interface{} {
		return fn()
	}
}

// Sets maximum amount of parallel go routine workers.
func (p *WorkerPool) SetMaxWorker(n int) {
	if n > 0 {
		p.worker = n
	}
}
