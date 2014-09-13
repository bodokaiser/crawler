package work

import (
	"testing"

	"gopkg.in/check.v1"
)

func TestWorker(t *testing.T) {
	check.Suite(&WorkerSuite{})
	check.TestingT(t)
}

type WorkerSuite struct {
	worker *Worker
	work1  *counter
	work2  *counter
}

func (s *WorkerSuite) SetUpTest(c *check.C) {
	s.worker = New()

	s.work1 = &counter{10, make(chan bool)}
	s.work2 = &counter{100, make(chan bool)}
}

func (s *WorkerSuite) TestAdd(c *check.C) {
	s.worker.Add(s.work1)
	s.worker.Add(s.work2)

	<-s.work1.Done
	<-s.work2.Done

	c.Check(s.work1.Value, check.Equals, 20)
	c.Check(s.work2.Value, check.Equals, 110)
}

type counter struct {
	Value int
	Done  chan bool
}

func (c *counter) Do() {
	c.Value += 10

	close(c.Done)
}
