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
	pool *Pool
}

func (s *PoolSuite) SetUpTest(c *check.C) {
	s.pool = NewPool(worker)
}

func (s *PoolSuite) TestEmpty(c *check.C) {
	res := s.pool.Get()

	c.Check(res, check.IsNil)
}

func (s *PoolSuite) TestSingle(c *check.C) {
	s.pool.Put(2)

	res := s.pool.Get()
	c.Check(res, check.Equals, 4)
}

func (s *PoolSuite) TestMultiple(c *check.C) {
	s.pool.Put(1)
	s.pool.Put(2)

	n := []int{
		s.pool.Get().(int),
		s.pool.Get().(int),
	}

	c.Check(n[0]+n[1], check.Equals, 7)
}

func worker(value interface{}) interface{} {
	if n, ok := value.(int); ok {
		return n + 2
	}

	return 0
}
