package store

import (
	"testing"

	"gopkg.in/check.v1"

	_ "github.com/go-sql-driver/mysql"
)

func TestStore(t *testing.T) {
	check.Suite(&StoreSuite{
		url:  "root@/gerenuk",
		url1: "http://example.org",
		url2: "http://company.com",
		url3: "http://fooobar.net",
	})
	check.TestingT(t)
}

type StoreSuite struct {
	url   string
	url1  string
	url2  string
	url3  string
	store *Store
}

func (s *StoreSuite) SetUpTest(c *check.C) {
	store, err := Open(s.url)
	c.Assert(err, check.IsNil)
	c.Assert(store, check.NotNil)

	err = store.Put(s.url1)
	c.Assert(err, check.IsNil)

	s.store = store
}

func (s *StoreSuite) TestPut(c *check.C) {
	err := s.store.Put(s.url2)
	c.Assert(err, check.IsNil)

	var n int

	row := s.store.db.QueryRow(`SELECT COUNT(*) FROM url WHERE url = ?`, s.url2)
	err = row.Scan(&n)
	c.Assert(err, check.IsNil)
	c.Check(n, check.Equals, 1)
}

func (s *StoreSuite) TestGet(c *check.C) {
	p, err := s.store.Get()
	c.Assert(err, check.IsNil)

	c.Check(p.Origin(), check.Equals, s.url1)
	c.Check(p.Refers(), check.HasLen, 0)

	p.AddRefer(s.url2)
	p.AddRefer(s.url3)

	c.Check(p.HasRefer(s.url2), check.Equals, true)
	c.Check(p.HasRefer(s.url3), check.Equals, true)
	c.Check(p.Refers(), check.DeepEquals, []string{s.url2, s.url3})

	err = p.Commit()
	c.Assert(err, check.IsNil)

	rows, err := s.store.db.Query(`SELECT id, url, done FROM url`)
	c.Assert(err, check.IsNil)

	var id int64
	var url string
	var done bool

	c.Assert(rows.Next(), check.Equals, true)
	err = rows.Scan(&id, &url, &done)
	c.Assert(err, check.IsNil)
	c.Check(id, check.Equals, int64(1))
	c.Check(url, check.Equals, s.url1)
	c.Check(done, check.Equals, true)

	c.Assert(rows.Next(), check.Equals, true)
	err = rows.Scan(&id, &url, &done)
	c.Assert(err, check.IsNil)
	c.Check(id, check.Equals, int64(2))
	c.Check(url, check.Equals, s.url2)
	c.Check(done, check.Equals, false)

	c.Assert(rows.Next(), check.Equals, true)
	err = rows.Scan(&id, &url, &done)
	c.Assert(err, check.IsNil)
	c.Check(id, check.Equals, int64(3))
	c.Check(url, check.Equals, s.url3)
	c.Check(done, check.Equals, false)

	c.Assert(rows.Next(), check.Equals, false)
	c.Assert(rows.Close(), check.IsNil)

	var oid, rid int64

	rows, err = s.store.db.Query(`SELECT * FROM ref`)
	c.Assert(err, check.IsNil)

	c.Assert(rows.Next(), check.Equals, true)
	err = rows.Scan(&oid, &rid)
	c.Assert(err, check.IsNil)
	c.Check(oid, check.Equals, int64(1))
	c.Check(rid, check.Equals, int64(2))

	c.Assert(rows.Next(), check.Equals, true)
	err = rows.Scan(&oid, &rid)
	c.Assert(err, check.IsNil)
	c.Check(oid, check.Equals, int64(1))
	c.Check(rid, check.Equals, int64(3))

	c.Assert(rows.Next(), check.Equals, false)
	c.Assert(rows.Close(), check.IsNil)
}

func (s *StoreSuite) TearDownTest(c *check.C) {
	err := s.store.Reset()
	c.Assert(err, check.IsNil)
}

func (s *StoreSuite) TearDownSuite(c *check.C) {
	err := s.store.Close()
	c.Assert(err, check.IsNil)
}
