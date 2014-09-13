package html

import (
	"bytes"
	"testing"

	"gopkg.in/check.v1"
)

func TestScanner(t *testing.T) {
	check.Suite(&ScannerSuite{})
	check.TestingT(t)
}

type ScannerSuite struct{}

func (s *ScannerSuite) TestHref(c *check.C) {
	data1 := []byte(`<h1>Hello World</h1>`)
	data2 := []byte(`<a href="foobar">Email</a>`)
	data3 := []byte(`<a href="foobar">Email</a><a href="helloworld">Hello</a>`)

	off1, sym1, err1 := ScanHref(data1, false)
	off2, sym2, err2 := ScanHref(data2, false)
	off3, sym3, err3 := ScanHref(data3, false)
	off4, sym4, err4 := ScanHref(data3[off3:], false)

	c.Check(err1, check.IsNil)
	c.Check(err2, check.IsNil)
	c.Check(err3, check.IsNil)
	c.Check(err4, check.IsNil)

	c.Check(sym1, check.IsNil)
	c.Check(sym2, check.DeepEquals, []byte("foobar"))
	c.Check(sym3, check.DeepEquals, []byte("foobar"))
	c.Check(sym4, check.DeepEquals, []byte("helloworld"))

	c.Check(off1, check.Equals, len(data1))
	c.Check(off2, check.Equals, bytes.LastIndex(data2, []byte(`">`))+1)
	c.Check(off3, check.Equals, bytes.Index(data3, []byte(`">`))+1)
	c.Check(off4, check.Equals, bytes.Index(data3[off3:], []byte(`">`))+1)
}
