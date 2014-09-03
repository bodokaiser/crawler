package store

import (
	"testing"

	"gopkg.in/check.v1"
)

func TestList(t *testing.T) {
	check.Suite(&ListSuite{})
	check.TestingT(t)
}

type ListSuite struct {
	list *List
}

func (s *ListSuite) SetUpTest(c *check.C) {
	s.list = NewList()
}

func (s *ListSuite) TestEmpty(c *check.C) {
	res := s.list.Has("http://www.google.com")

	c.Check(res, check.Equals, false)
}

func (s *ListSuite) TestSingle(c *check.C) {
	s.list.Add("http://www.example.org")

	res := s.list.Has("http://www.EXAMPLE.org#hello")
	c.Check(res, check.Equals, true)
}
