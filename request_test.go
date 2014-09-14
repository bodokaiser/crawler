package crawler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
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
	s.server = httptest.NewServer(http.HandlerFunc(s.handler))
}

func (s *RequestSuite) SetUpTest(c *check.C) {
	s.url1, _ = url.Parse(s.server.URL)
	s.url2, _ = url.Parse("http://example.org")
	s.url3, _ = url.Parse("https://example.com")
}

func (s *RequestSuite) TestDo(c *check.C) {
	r, err := NewRequest(s.server.URL)
	c.Assert(err, check.IsNil)
	r.Do()

	c.Check(r.Refers, check.DeepEquals, []*url.URL{
		s.url2,
		s.url3,
	})
}

func (s *RequestSuite) TestHas(c *check.C) {
	r := Request{Refers: []*url.URL{s.url2}}

	s.url2.Fragment = "foobar"
	s.url2.Host = strings.ToTitle(s.url2.Host)

	c.Check(r.Has(s.url2), check.Equals, true)
	c.Check(r.Has(s.url2), check.Equals, true)
	c.Check(r.Has(s.url3), check.Equals, false)
}

func (s *RequestSuite) TestPush(c *check.C) {
	r := Request{Refers: []*url.URL{}}
	r.Push(s.url3)

	c.Check(r.Refers, check.DeepEquals, []*url.URL{s.url3})
}

func (s *RequestSuite) handler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, `
		<!DOCTYPE html>
		<html>
			<head></head>
			<body>
				<h1>Example</h1>
				<a href="http://example.org"></a>
				<a href="https://example.com"></a>
			</body>
		</html>
	`)
}
