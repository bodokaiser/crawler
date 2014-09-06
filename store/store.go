package store

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type Store struct {
	db *sql.DB
}

var sqlCreateUrlTable = `
	CREATE TABLE IF NOT EXISTS url (
	  id   bigint(40)   unsigned NOT NULL AUTO_INCREMENT,
	  url  varchar(255)          NOT NULL,
	  done tinyint(1)            NOT NULL DEFAULT 0,
	  PRIMARY KEY (id),
	  UNIQUE KEY UNIQUE_url (url)
	) ENGINE=XtraDB DEFAULT CHARSET=utf8;
`

var sqlCreateRefTable = `
	CREATE TABLE IF NOT EXISTS ref (
	  origin_id bigint(40) unsigned NOT NULL,
	  refer_id  bigint(40) unsigned NOT NULL,
	  PRIMARY KEY (origin_id, refer_id),
	  KEY refer (refer_id),
	  CONSTRAINT origin FOREIGN KEY (origin_id) REFERENCES url (id),
	  CONSTRAINT refer  FOREIGN KEY (refer_id)  REFERENCES url (id)
	) ENGINE=XtraDB DEFAULT CHARSET=utf8;
`

var sqlDropTables = `
	DROP TABLE IF EXISTS ref, url;
`

func NewStore(url string) (*Store, error) {
	db, err := sql.Open("mysql", url)
	if err != nil {
		return nil, err
	}

	return &Store{
		db: db,
	}, nil
}

func (s *Store) EnsureTables() error {
	_, err := s.db.Exec(sqlCreateUrlTable)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(sqlCreateRefTable)

	return err
}

func (s *Store) DropTables() error {
	_, err := s.db.Exec(sqlDropTables)

	return err
}

func (s *Store) Begin() (*Tx, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}

	return newTx(tx)
}

func (s *Store) Insert(url string) error {
	_, err := s.db.Exec(sqlInsertUrl, url)

	return err
}

func (s *Store) Close() error {
	return s.db.Close()
}
