package pool

import (
	"sync/atomic"
	"time"
)

// Defines the amount of work to queue.
var MaxQueue = 500

// Defines the amount of worker spawned by pool.
var MaxWorker = 100

// Defines the amount of time a worker is in idle.
var MaxTimeout = time.Second

// Type Pool implements a worker pool which is generic to the work done.
type Pool struct {
	in      chan interface{}
	out     chan interface{}
	active  int32
	pending int32
	worker  Worker
}

// Returns initialized worker pool with provided function as worker.
func NewPool(w Worker) *Pool {
	return &Pool{
		in:     make(chan interface{}, MaxQueue),
		out:    make(chan interface{}, MaxQueue),
		worker: w,
	}
}

// Returns an interface type value which was returned by a worker.
// Blocks when values are pending.
func (p *Pool) Get() interface{} {
	if atomic.LoadInt32(&p.pending) > 0 {
		atomic.AddInt32(&p.pending, -1)

		return <-p.out
	}

	return nil
}

// Puts an interface type value into the pool which will be passed to a worker.
// Blocks when all workers are busy until one is free again.
func (p *Pool) Put(i interface{}) {
	atomic.AddInt32(&p.pending, 1)

	pen := atomic.LoadInt32(&p.pending)
	act := atomic.LoadInt32(&p.active)

	if pen > act && act < int32(MaxWorker) {
		go func(in <-chan interface{}, out chan<- interface{}, w Worker) {
			for {
				select {
				case i := <-in:
					out <- w(i)
				case <-time.After(MaxTimeout):
					atomic.AddInt32(&p.active, -1)

					return
				}
			}
		}(p.in, p.out, p.worker)

		atomic.AddInt32(&p.active, 1)
	}

	p.in <- i
}

// The Worker type defins a function which receives a value put into the pool
// which then may be processed concurrently to return another result which can
// be received through get on the pool.
type Worker func(interface{}) interface{}
