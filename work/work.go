package work

import (
	"sync/atomic"
	"time"
)

var (
	DefaultTimeout         = time.Second
	DefaultMaxQueue  int32 = 200
	DefaultMaxWorker int32 = 100
)

type Work interface {
	Do()
}

type Worker struct {
	queue   chan Work
	worker  int32
	active  int32
	pending int32
	timeout time.Duration
}

func NewWorker() *Worker {
	w := &Worker{}
	w.SetTimeout(DefaultTimeout)
	w.SetMaxQueue(DefaultMaxQueue)
	w.SetMaxWorker(DefaultMaxWorker)

	return w
}

func (w *Worker) Add(work Work) {
	atomic.AddInt32(&w.pending, 1)

	if a, p := atomic.LoadInt32(&w.active), atomic.LoadInt32(&w.pending); p > a && a < int32(w.worker) {
		go func(queue <-chan Work, timeout time.Duration) {
			for {
				select {
				case work := <-queue:
					work.Do()

					atomic.AddInt32(&w.pending, -1)
				case <-time.After(timeout):
					atomic.AddInt32(&w.active, -1)

					return
				}
			}
		}(w.queue, w.timeout)

		atomic.AddInt32(&w.active, 1)
	}

	w.queue <- work
}

func (w *Worker) SetTimeout(d time.Duration) {
	w.timeout = d
}

func (w *Worker) SetMaxQueue(n int32) {
	w.queue = make(chan Work, n)
}

func (w *Worker) SetMaxWorker(n int32) {
	w.worker = n
}
