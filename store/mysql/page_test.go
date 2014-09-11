package mysql

import (
	"database/sql"
	"strings"
	"testing"

	"gopkg.in/check.v1"

	"github.com/bodokaiser/gerenuk/store"
)

func TestPage(t *testing.T) {
	check.Suite(&PageSuite{
		url: "root@/gerenuk?timeout=20s",
	})
	check.TestingT(t)
}

type PageSuite struct {
	url   string
	db    *sql.DB
	store *PageStore
}

func (s *PageSuite) SetUpSuite(c *check.C) {
	db, err := sql.Open("mysql", s.url)
	c.Assert(err, check.IsNil)

	s.db = db
}

func (s *PageSuite) SetUpTest(c *check.C) {
	s.store = &PageStore{db: s.db}

	_, err := s.db.Exec(sqlCreateTableUrl)
	c.Assert(err, check.IsNil)
	_, err = s.db.Exec(sqlCreateTableRef)
	c.Assert(err, check.IsNil)
}

func (s *PageSuite) TestInsert(c *check.C) {
	p := page{
		origin: "http://example.org",
		refers: []string{
			"http://example1.org",
			"http://example2.org",
		},
	}

	err := s.store.Insert(&p)
	c.Assert(err, check.IsNil)

	var origin, refers string
	err = s.db.QueryRow(sqlSelect).Scan(&origin, &refers)
	c.Assert(err, check.IsNil)
	c.Check(origin, check.Equals, p.Origin())
	c.Check(refers, check.Equals, strings.Join(p.Refers(), ","))

	err = s.store.Insert(&p)
	c.Assert(err, check.Equals, store.ErrDupRow)
}

func (s *PageSuite) TestReset(c *check.C) {
	_, err := s.db.Exec(sqlInsertUrl, "http://foo.bar")
	c.Assert(err, check.IsNil)

	err = s.store.Reset()
	c.Assert(err, check.IsNil)

	var count int
	err = s.db.QueryRow("SELECT COUNT(*) FROM url").Scan(&count)
	c.Assert(err, check.IsNil)
	c.Check(count, check.Equals, 0)
}

func (s *PageSuite) TearDownTest(c *check.C) {
	_, err := s.db.Exec(sqlDropTables)
	c.Assert(err, check.IsNil)
}

func (s *PageSuite) TearDownSuite(c *check.C) {
	s.db.Close()
}

type page struct {
	origin string
	refers []string
}

func (p page) Origin() string {
	return p.origin
}

func (p page) Refers() []string {
	return p.refers
}
