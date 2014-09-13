package mysql

import (
	"database/sql"

	"github.com/bodokaiser/crawler/store"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	store.Register("mysql", &Driver{})
}

type Driver struct{}

func (d *Driver) Open(url string) (store.Store, error) {
	db, err := sql.Open("mysql", url)
	if err != nil {
		return nil, err
	}

	return &Store{db: db}, nil
}

type Store struct {
	db *sql.DB
}

func (s *Store) Page() (store.PageStore, error) {
	p := &PageStore{db: s.db}

	return p, p.create()
}

func (s *Store) Close() error {
	return s.db.Close()
}
