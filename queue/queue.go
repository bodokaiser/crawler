package queue

import "sync"

// Type Queue is a concurrent safe list.
type Queue struct {
	mutex  *sync.Mutex
	queue  []interface{}
	length int
}

func NewQueue() *Queue {
	return &Queue{
		mutex: new(sync.Mutex),
		queue: make([]interface{}, 0),
	}
}

// Returns the last item from Queue.
func (q *Queue) Pull() interface{} {
	if q.length == 0 {
		return nil
	}
	q.mutex.Lock()
	x := q.queue[0]
	q.queue = q.queue[1:]
	q.length--
	q.mutex.Unlock()
	return x
}

// Pushes the provided item to Queue.
func (q *Queue) Push(x interface{}) {
	q.mutex.Lock()
	q.queue = append(q.queue, x)
	q.length++
	q.mutex.Unlock()
}

// Returns the total length of all items.
func (q *Queue) Len() int {
	return q.length
}
