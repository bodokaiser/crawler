package worker

import (
	"testing"

	"gopkg.in/check.v1"
)

func TestWorker(t *testing.T) {
	check.Suite(&WorkerSuite{})
	check.TestingT(t)
}

type WorkerSuite struct {
	wpool *WorkerPool
	work1 *Work
	work2 *Work
}

func (s *WorkerSuite) SetUpTest(c *check.C) {
	s.wpool = NewWorkerPool()
	s.work1 = &Work{
		Done:   make(chan bool),
		Func:   work,
		Params: []interface{}{1},
	}
	s.work2 = &Work{
		Done:   make(chan bool),
		Func:   work,
		Params: []interface{}{1},
	}
}

func (s *WorkerSuite) TestSingle(c *check.C) {
	s.wpool.Put(s.work1)

	<-s.work1.Done

	c.Check(s.work1.Result.(int), check.Equals, 3)
}

func (s *WorkerSuite) TestMultiple(c *check.C) {
	s.wpool.Put(s.work1)
	s.wpool.Put(s.work2)

	<-s.work1.Done
	<-s.work2.Done

	c.Check(s.work1.Result.(int), check.Equals, 3)
	c.Check(s.work2.Result.(int), check.Equals, 3)
}

func work(params ...interface{}) (interface{}, error) {
	if n, ok := params[0].(int); ok {
		return n + 2, nil
	}
	return 0, nil
}
