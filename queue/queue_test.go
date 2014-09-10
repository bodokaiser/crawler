package queue

import (
	"sync"
	"testing"

	"gopkg.in/check.v1"
)

func TestQueue(t *testing.T) {
	check.Suite(&QueueSuite{})
	check.TestingT(t)
}

type QueueSuite struct {
	queue *Queue
}

func (s *QueueSuite) SetUpTest(c *check.C) {
	s.queue = NewQueue()
	s.queue.queue = append(s.queue.queue, 1, 2, 3)
	s.queue.length += 3
}

func (s *QueueSuite) TestPull(c *check.C) {
	x1 := s.queue.Pull()
	x2 := s.queue.Pull()
	x3 := s.queue.Pull()
	x4 := s.queue.Pull()

	c.Check(x1.(int), check.Equals, 1)
	c.Check(x2.(int), check.Equals, 2)
	c.Check(x3.(int), check.Equals, 3)
	c.Check(x4, check.IsNil)
}

func (s *QueueSuite) TestPush(c *check.C) {
	s.queue.Push(4)
	c.Check(s.queue.queue, check.DeepEquals, []interface{}{1, 2, 3, 4})

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(q *Queue, wg *sync.WaitGroup) {
			for i := 0; i < 10000; i++ {
				s.queue.Push(i)
			}
			wg.Done()
		}(s.queue, &wg)
	}
	wg.Wait()
	c.Check(s.queue.length, check.Equals, len(s.queue.queue))
	c.Check(s.queue.length, check.Equals, 100*10000+4)
}

func (s *QueueSuite) TestLen(c *check.C) {
	c.Check(s.queue.Len(), check.Equals, 3)
}
