package worker

import (
	"sync/atomic"
	"time"
)

var (
	DefaultTimeout   = time.Second
	DefaultMaxQueue  = 200
	DefaultMaxWorker = 100
)

type Work struct {
	Func   func(...interface{}) (interface{}, error)
	Done   chan bool
	Error  error
	Result interface{}
	Params []interface{}
}

type Pool struct {
	queue   chan *Work
	worker  int
	active  int32
	pending int32
	timeout time.Duration
}

func NewPool() *Pool {
	p := &Pool{}
	p.SetTimeout(DefaultTimeout)
	p.SetMaxQueue(DefaultMaxQueue)
	p.SetMaxWorker(DefaultMaxWorker)

	return p
}

func (p *Pool) Put(w *Work) {
	atomic.AddInt32(&p.pending, 1)
	pen := atomic.LoadInt32(&p.pending)
	act := atomic.LoadInt32(&p.active)

	if pen > act && act < int32(p.worker) {
		go func(queue <-chan *Work, timeout time.Duration) {
			for {
				select {
				case w := <-queue:
					w.Result, w.Error = w.Func(w.Params...)

					if w.Done != nil {
						w.Done <- true
						close(w.Done)
					}

					atomic.AddInt32(&p.pending, -1)
				case <-time.After(timeout):
					atomic.AddInt32(&p.active, -1)

					return
				}
			}
		}(p.queue, p.timeout)

		atomic.AddInt32(&p.active, 1)
	}

	p.queue <- w
}

func (p *Pool) SetTimeout(d time.Duration) {
	p.timeout = d
}

func (p *Pool) SetMaxQueue(n int) {
	p.queue = make(chan *Work, n)
}

func (p *Pool) SetMaxWorker(n int) {
	p.worker = n
}
