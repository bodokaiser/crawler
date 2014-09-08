package pool

import (
	"testing"

	"gopkg.in/check.v1"
)

func TestPool(t *testing.T) {
	check.Suite(&PoolSuite{})
	check.TestingT(t)
}

type PoolSuite struct {
	pool  *WorkerPool
	work1 *Work
	work2 *Work
}

func (s *PoolSuite) SetUpTest(c *check.C) {
	s.pool = NewWorkerPool()
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

func (s *PoolSuite) TestSingle(c *check.C) {
	s.pool.Put(s.work1)

	<-s.work1.Done

	c.Check(s.work1.Result.(int), check.Equals, 3)
}

func (s *PoolSuite) TestMultiple(c *check.C) {
	s.pool.Put(s.work1)
	s.pool.Put(s.work2)

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
