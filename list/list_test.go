package list

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
	item, err := NewItemFromUrl("http://www.google.com")
	c.Check(err, check.IsNil)

	res := s.list.Has(item)
	c.Check(res, check.Equals, false)
}

func (s *ListSuite) TestSingle(c *check.C) {
	item1, err1 := NewItemFromUrl("http://www.example.org")
	item2, err2 := NewItemFromUrl("http://www.EXAMPLE.org#hello")
	c.Check(err1, check.IsNil)
	c.Check(err2, check.IsNil)

	s.list.Add(item1)

	res := s.list.Has(item2)
	c.Check(res, check.Equals, true)
}
