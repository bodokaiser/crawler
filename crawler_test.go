package crawler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"gopkg.in/check.v1"
)

func TestCrawler(t *testing.T) {
	check.Suite(&CrawlerSuite{})
	check.TestingT(t)
}

type CrawlerSuite struct {
	request *Request
	crawler *Crawler
	server  *httptest.Server
}

func (s *CrawlerSuite) SetUpSuite(c *check.C) {
	s.server = httptest.NewServer(http.HandlerFunc(handler))
}

func (s *CrawlerSuite) SetUpTest(c *check.C) {
	r, err := NewRequest(s.server.URL)
	c.Assert(err, check.IsNil)

	s.crawler = New()
	s.request = r
}

func (s *CrawlerSuite) TestAdd(c *check.C) {
	s.crawler.Add(s.request)

	_, ok := <-s.request.Done
	c.Assert(ok, check.Equals, false)
}

func (s *CrawlerSuite) TestGet(c *check.C) {
	s.crawler.Add(s.request)

	r := s.crawler.Get()
	c.Check(r, check.HasLen, 1)
	c.Check(r[0], check.DeepEquals, s.request)
}

func handler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, `
		<!DOCTYPE html>
		<html>
			<head></head>
			<body>
				<h1>Example</h1>
				<a href="foo.html"></a>
				<a href="/foo.html"></a>
				<a href="/bar.html"></a>
				<a href="//example.org"></a>
				<a href="http://example.org"></a>
				<a href="https://example.org"></a>
			</body>
		</html>
	`)
}
