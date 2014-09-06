package store

import (
	"database/sql"
	"testing"

	"gopkg.in/check.v1"

	_ "github.com/go-sql-driver/mysql"
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

func (s *StoreSuite) SetUpSuite(c *check.C) {
	db, err := sql.Open("mysql", s.url)
	c.Assert(err, check.IsNil)

	s.db = db
	s.store = &Store{
		db: db,
	}
}

func (s *StoreSuite) SetUpTest(c *check.C) {
	_, err := s.db.Exec(sqlCreateUrlTable)
	c.Assert(err, check.IsNil)

	_, err = s.db.Exec(sqlCreateRefTable)
	c.Assert(err, check.IsNil)

	_, err = s.db.Exec(sqlInsertUrl, "http://github.com/bodokaiser")
	c.Assert(err, check.IsNil)
}

func (s *StoreSuite) TestEnsureTables(c *check.C) {
	s.TearDownTest(c)

	err := s.store.EnsureTables()
	c.Assert(err, check.IsNil)

	var n int

	err = s.db.QueryRow(`
		SELECT COUNT(*)
		FROM information_schema.tables
		WHERE table_schema = ?
			AND (table_name = ? OR table_name = ?)
	`, "gerenuk", "ref", "url").Scan(&n)

	c.Assert(err, check.IsNil)
	c.Check(n, check.Equals, 2)
}

func (s *StoreSuite) TestDropTables(c *check.C) {
	err := s.store.DropTables()
	c.Assert(err, check.IsNil)

	var n int

	err = s.db.QueryRow(`
		SELECT COUNT(*)
		FROM information_schema.tables
		WHERE table_schema = ?
			AND (table_name = ? OR table_name = ?)
	`, "gerenuk", "ref", "url").Scan(&n)

	c.Assert(err, check.IsNil)
	c.Check(n, check.Equals, 0)
}

func (s *StoreSuite) TestInsert(c *check.C) {
	err := s.store.Insert("http://www.satisfeet.me")
	c.Assert(err, check.IsNil)

	var n int

	err = s.db.QueryRow(`
		SELECT COUNT(*)
		FROM url
		WHERE url = ?
	`, "http://www.satisfeet.me").Scan(&n)
	c.Assert(err, check.IsNil)
}

func (s *StoreSuite) TestBegin(c *check.C) {
	tx, err := s.store.Begin()
	c.Assert(err, check.IsNil)

	origin := tx.Origin()
	c.Check(origin, check.Equals, "http://github.com/bodokaiser")

	err = tx.AddRefer("http://www.satisfeet.me")
	c.Assert(err, check.IsNil)
	err = tx.AddRefer("http://www.satisfeet.me/products")
	c.Assert(err, check.IsNil)

	refers := tx.Refers()
	c.Check(refers, check.DeepEquals, []string{
		"http://www.satisfeet.me",
		"http://www.satisfeet.me/products",
	})

	err = tx.Commit()
	c.Assert(err, check.IsNil)

	rows, err := s.db.Query(`SELECT id, url, done FROM url`)
	c.Assert(err, check.IsNil)

	var id int64
	var url string
	var done bool

	c.Assert(rows.Next(), check.Equals, true)
	err = rows.Scan(&id, &url, &done)
	c.Assert(err, check.IsNil)
	c.Check(id, check.Equals, int64(1))
	c.Check(url, check.Equals, "http://github.com/bodokaiser")
	c.Check(done, check.Equals, true)

	c.Assert(rows.Next(), check.Equals, true)
	err = rows.Scan(&id, &url, &done)
	c.Assert(err, check.IsNil)
	c.Check(id, check.Equals, int64(2))
	c.Check(url, check.Equals, "http://www.satisfeet.me")
	c.Check(done, check.Equals, false)

	c.Assert(rows.Next(), check.Equals, true)
	err = rows.Scan(&id, &url, &done)
	c.Assert(err, check.IsNil)
	c.Check(id, check.Equals, int64(3))
	c.Check(url, check.Equals, "http://www.satisfeet.me/products")
	c.Check(done, check.Equals, false)

	c.Assert(rows.Next(), check.Equals, false)
	c.Assert(rows.Close(), check.IsNil)

	var originId, referId int64

	rows, err = s.db.Query(`SELECT * FROM ref`)
	c.Assert(err, check.IsNil)

	c.Assert(rows.Next(), check.Equals, true)
	err = rows.Scan(&originId, &referId)
	c.Assert(err, check.IsNil)
	c.Check(originId, check.Equals, int64(1))
	c.Check(referId, check.Equals, int64(2))

	c.Assert(rows.Next(), check.Equals, true)
	err = rows.Scan(&originId, &referId)
	c.Assert(err, check.IsNil)
	c.Check(originId, check.Equals, int64(1))
	c.Check(referId, check.Equals, int64(3))

	c.Assert(rows.Next(), check.Equals, false)
	c.Assert(rows.Close(), check.IsNil)
}

func (s *StoreSuite) TearDownTest(c *check.C) {
	_, err := s.db.Exec(sqlDropTables)
	c.Assert(err, check.IsNil)
}

func (s *StoreSuite) TearDownSuite(c *check.C) {
	err := s.db.Close()
	c.Assert(err, check.IsNil)
}
