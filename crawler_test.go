package gerenuk

import (
	"fmt"
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
	empty    *httptest.Server
	single   *httptest.Server
	multiple *httptest.Server
}

func (s *CrawlerSuite) SetUpSuite(c *check.C) {
	s.empty = httptest.NewServer(http.HandlerFunc(emptyHandler))
	s.single = httptest.NewServer(http.HandlerFunc(singleHandler))
	s.multiple = httptest.NewServer(http.HandlerFunc(multipleHandler))
}

func (s *CrawlerSuite) TestEmpty(c *check.C) {
	crawl := NewCrawler(s.empty.URL)

	url1, err1 := crawl.Get()
	url2, err2 := crawl.Get()

	c.Check(err1, check.IsNil)
	c.Check(err2, check.IsNil)

	c.Check(url1, check.Equals, s.empty.URL)
	c.Check(url2, check.Equals, "")
}

func (s *CrawlerSuite) TestSingle(c *check.C) {
	crawl := NewCrawler(s.single.URL)

	url1, err1 := crawl.Get()
	url2, err2 := crawl.Get()
	url3, err3 := crawl.Get()

	c.Check(err1, check.IsNil)
	c.Check(err2, check.IsNil)
	c.Check(err3, check.IsNil)

	c.Check(url1, check.Equals, s.single.URL)
	c.Check(url2, check.Equals, s.single.URL+"/foo")
	c.Check(url3, check.Equals, "")
}

func (s *CrawlerSuite) TestMultiple(c *check.C) {
	crawler, err := NewCrawler()

	url1, err1 := crawl.Get()
	url2, err2 := crawl.Get()
	url3, err3 := crawl.Get()
	url4, err4 := crawl.Get()

	c.Check(err1, check.IsNil)
	c.Check(err2, check.IsNil)
	c.Check(err3, check.IsNil)
	c.Check(err4, check.IsNil)

	c.Check(url1, check.Equals, s.multiple.URL)
	c.Check(url2, check.Equals, s.multiple.URL+"/foo")
	c.Check(url3, check.Equals, s.multiple.URL+"/bar")
	c.Check(url4, check.Equals, "")
}

func (s *CrawlerSuite) TearDownSuite(c *check.C) {
	s.empty.Close()
	s.single.Close()
	s.multiple.Close()
}

func emptyHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}

func singleHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `<a href="/foo"></a>`)
}

func multipleHandler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		fmt.Fprintf(w, `<a href="/foo"></a>`)
	case "/foo":
		fmt.Fprintf(w, `<a href="/bar"></a>`)
	}
}
