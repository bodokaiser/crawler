package httpd

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"gopkg.in/check.v1"
)

func TestPool(t *testing.T) {
	check.Suite(&PoolSuite{})
	check.TestingT(t)
}

type PoolSuite struct {
	pool   *Pool
	server *httptest.Server
}

func (s *PoolSuite) SetUpSuite(c *check.C) {
	s.server = httptest.NewServer(http.HandlerFunc(handle))
}

func (s *PoolSuite) SetUpTest(c *check.C) {
	s.pool = NewPool()
}

func (s *PoolSuite) TestEmpty(c *check.C) {
	req, res, err := s.pool.Get()
	c.Check(err, check.IsNil)
	c.Check(req, check.IsNil)
	c.Check(res, check.IsNil)
}

func (s *PoolSuite) TestSingle(c *check.C) {
	request, _ := http.NewRequest("GET", s.server.URL, nil)

	s.pool.Add(request)
	s.pool.Run()

	req, res, err := s.pool.Get()
	c.Check(err, check.IsNil)
	c.Check(req.URL.String(), check.Equals, s.server.URL)
	c.Check(res.StatusCode, check.Equals, http.StatusOK)

	req, res, err = s.pool.Get()
	c.Check(err, check.IsNil)
	c.Check(req, check.IsNil)
	c.Check(res, check.IsNil)
}

func (s *PoolSuite) TestMultiple(c *check.C) {
	request1, _ := http.NewRequest("GET", s.server.URL, nil)
	request2, _ := http.NewRequest("GET", s.server.URL, nil)

	s.pool.Add(request1)
	s.pool.Add(request2)
	s.pool.Run()

	req, res, err := s.pool.Get()
	c.Check(err, check.IsNil)
	c.Check(req.URL.String(), check.Equals, s.server.URL)
	c.Check(res.StatusCode, check.Equals, http.StatusOK)

	req, res, err = s.pool.Get()
	c.Check(err, check.IsNil)
	c.Check(req.URL.String(), check.Equals, s.server.URL)
	c.Check(res.StatusCode, check.Equals, http.StatusOK)

	req, res, err = s.pool.Get()
	c.Check(err, check.IsNil)
	c.Check(req, check.IsNil)
	c.Check(res, check.IsNil)
}

func (s *PoolSuite) TestInvalid(c *check.C) {
	request, _ := http.NewRequest("GET", "error.dns", nil)

	s.pool.Add(request)
	s.pool.Run()

	req, res, err := s.pool.Get()
	c.Check(err, check.NotNil)
	c.Check(req, check.NotNil)
	c.Check(res, check.IsNil)
}

func handle(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}
