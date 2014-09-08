package store

import (
	"database/sql"
	"net/url"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

const (
	sqlCreateUrlTable = `
		CREATE TABLE IF NOT EXISTS url (
		  id   bigint(40)   unsigned NOT NULL AUTO_INCREMENT,
		  url  varchar(255)          NOT NULL,
		  done tinyint(1)            NOT NULL DEFAULT 0,
		  PRIMARY KEY (id),
		  UNIQUE KEY UNIQUE_url (url)
		) ENGINE=XtraDB DEFAULT CHARSET=utf8;
	`
	sqlCreateRefTable = `
		CREATE TABLE IF NOT EXISTS ref (
		  origin_id bigint(40) unsigned NOT NULL,
		  refer_id  bigint(40) unsigned NOT NULL,
		  PRIMARY KEY (origin_id, refer_id),
		  KEY refer (refer_id),
		  CONSTRAINT origin FOREIGN KEY (origin_id) REFERENCES url (id),
		  CONSTRAINT refer  FOREIGN KEY (refer_id)  REFERENCES url (id)
		) ENGINE=XtraDB DEFAULT CHARSET=utf8;
	`
	sqlDropTables = `
		DROP TABLE ref, url
	`
	sqlTruncateRefTable = `
		TRUNCATE TABLE ref
	`
	sqlSelectUrl = `
		SELECT id, url
		FROM url
		WHERE done = false
		LIMIT 1
		FOR UPDATE
	`
	sqlUpdateUrl = `
		UPDATE url
		SET done = ?
		WHERE id = ?
	`
	sqlInsertUrl = `
		INSERT into url (url)
		VALUES (?)
		ON DUPLICATE KEY UPDATE id = LAST_INSERT_ID(id)
	`
	sqlInsertRef = `
		INSERT into ref (origin_id, refer_id)
		VALUES (?, ?)
	`
)

func Open(url string) (*Store, error) {
	db, err := sql.Open("mysql", url)
	if err != nil {
		return nil, err
	}

	s := &Store{db: db}

	return s, s.create()
}

type Page struct {
	tx *sql.Tx

	id     int64
	origin string
	refers []string
}

func (p *Page) Origin() string {
	return p.origin
}

func (p *Page) Refers() []string {
	return p.refers
}

func (p *Page) HasRefer(ref string) bool {
	uri, err := parse(ref)
	if err != nil {
		return false
	}
	ref = uri.String()

	for _, r := range p.refers {
		if r == ref {
			return true
		}
	}

	return false
}

func (p *Page) AddRefer(ref string) {
	p.refers = append(p.refers, ref)
}

func (p *Page) Abort() error {
	return p.tx.Rollback()
}

func (p *Page) Commit() error {
	insUrlStmt, err := p.tx.Prepare(sqlInsertUrl)
	if err != nil {
		p.Abort()

		return err
	}

	insRefStmt, err := p.tx.Prepare(sqlInsertRef)
	if err != nil {
		p.Abort()

		return err
	}

	for _, ref := range p.refers {
		res, err := insUrlStmt.Exec(ref)
		if err != nil {
			p.Abort()

			return err
		}

		id, err := res.LastInsertId()
		if err != nil {
			p.Abort()

			return err
		}

		_, err = insRefStmt.Exec(p.id, id)
		if err != nil {
			p.Abort()

			return err
		}
	}

	_, err = p.tx.Exec(sqlUpdateUrl, p.id, true)
	if err != nil {
		p.Abort()

		return err
	}

	return p.tx.Commit()
}

type Store struct {
	db *sql.DB
}

func (s *Store) Get() (*Page, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}

	p := &Page{tx: tx}

	if err := tx.QueryRow(sqlSelectUrl).Scan(&p.id, &p.origin); err != nil {
		tx.Rollback()

		return nil, err
	}

	return p, nil
}

func (s *Store) Put(urlStr string) error {
	_, err := s.db.Exec(sqlInsertUrl, urlStr)

	return err
}

func (s *Store) create() error {
	_, err := s.db.Exec(sqlCreateUrlTable)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(sqlCreateRefTable)

	return err
}

func (s *Store) Reset() error {
	_, err := s.db.Exec(sqlDropTables)
	if err != nil {
		return err
	}

	return s.create()
}

func (s *Store) Close() error {
	return s.db.Close()
}

func parse(urlStr string) (*url.URL, error) {
	uri, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	uri.Host = strings.ToLower(uri.Host)
	uri.Scheme = strings.ToLower(uri.Scheme)
	uri.Fragment = ""

	return uri, nil
}
