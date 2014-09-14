package crawler

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"gopkg.in/check.v1"
)

func TestRequest(t *testing.T) {
	check.Suite(&RequestSuite{})
	check.TestingT(t)
}

type RequestSuite struct {
	url1   *url.URL
	url2   *url.URL
	url3   *url.URL
	server *httptest.Server
}

func (s *RequestSuite) SetUpSuite(c *check.C) {
	s.server = httptest.NewServer(http.HandlerFunc(handler))

	s.url1, _ = url.Parse(s.server.URL)
	s.url2, _ = url.Parse("http://example.org")
	s.url3, _ = url.Parse("https://example.com")
}

func (s *RequestSuite) estDo(c *check.C) {
	r, err := NewRequest(s.server.URL)
	c.Assert(err, check.IsNil)
	r.Do()

	c.Check(r.Origin, check.Equals, s.server.URL)
	c.Check(r.Refers, check.DeepEquals, []*url.URL{})
}

func (s *RequestSuite) TestHas(c *check.C) {
	r := Request{Refers: []*url.URL{s.url2}}

	c.Check(r.Has(s.url2), check.Equals, true)
	c.Check(r.Has(s.url3), check.Equals, false)
}

func (s *RequestSuite) TestPush(c *check.C) {
	r := Request{Refers: []*url.URL{}}
	r.Push(s.url3)

	c.Check(r.Refers, check.DeepEquals, []*url.URL{s.url3})
}
