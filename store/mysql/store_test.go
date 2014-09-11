package mysql

import (
	"database/sql"
	"testing"

	"gopkg.in/check.v1"
)

func TestStore(t *testing.T) {
	check.Suite(&StoreSuite{
		url: "root@/gerenuk",
	})
	check.TestingT(t)
}

type StoreSuite struct {
	url   string
	db    *sql.DB
	store *Store
}

func (s *StoreSuite) SetUpTest(c *check.C) {
	db, err := sql.Open("mysql", s.url)
	c.Assert(err, check.IsNil)

	s.db = db
	s.store = &Store{db: db}
}

func (s *StoreSuite) TestOpen(c *check.C) {
	store, err := (&Driver{}).Open(s.url)
	c.Assert(err, check.IsNil)

	err = store.Close()
	c.Assert(err, check.IsNil)
}

func (s *StoreSuite) TestPage(c *check.C) {
	_, err := s.store.Page()
	c.Assert(err, check.IsNil)
}

func (s *StoreSuite) TearDownTest(c *check.C) {
	err := s.db.Close()
	c.Assert(err, check.IsNil)
}
