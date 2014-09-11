package mysql

import (
	"database/sql"

	"github.com/bodokaiser/gerenuk/store"
)

const (
	sqlCreateTableUrl = `
		CREATE TABLE IF NOT EXISTS url (
		  id   bigint(40)   unsigned NOT NULL AUTO_INCREMENT,
		  url  varchar(255)          NOT NULL,
		  PRIMARY KEY (id),
		  UNIQUE KEY UNIQUE_url (url)
		) ENGINE=XtraDB DEFAULT CHARSET=utf8;
	`

	sqlCreateTableRef = `
		CREATE TABLE IF NOT EXISTS ref (
		  origin_id bigint(40) unsigned NOT NULL,
		  refer_id  bigint(40) unsigned NOT NULL,
		  PRIMARY KEY (origin_id, refer_id),
		  KEY refer (refer_id),
		  CONSTRAINT origin FOREIGN KEY (origin_id) REFERENCES url (id),
		  CONSTRAINT refer  FOREIGN KEY (refer_id)  REFERENCES url (id)
		) ENGINE=XtraDB DEFAULT CHARSET=utf8;
	`

	sqlSelect = `
		SELECT
			origin.url AS origin,
			GROUP_CONCAT(refer.url) AS refers
		FROM ref
			LEFT JOIN url origin
				ON origin_id = origin.id
			LEFT JOIN url refer
				ON refer_id = refer.id
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

	sqlDropTables = `
		DROP TABLE ref, url
	`
)

type PageStore struct {
	db *sql.DB
}

func (s *PageStore) create() error {
	_, err := s.db.Exec(sqlCreateTableUrl)
	if err != nil {
		return err
	}
	_, err = s.db.Exec(sqlCreateTableRef)

	return err
}

func (s *PageStore) Insert(p *store.Page) error {
	tx, err := newPageTx(s.db)
	if err != nil {
		return err
	}

	oid, err := tx.InsertUrl(p.Origin)
	if err != nil {
		return err
	}

	for _, ref := range p.Refers {
		rid, err := tx.InsertUrl(ref)
		if err != nil {
			return err
		}

		err = tx.InsertRef(oid, rid)
		if err != nil {
			return err
		}
	}

	return tx.tx.Commit()
}

func (s *PageStore) Reset() error {
	_, err := s.db.Exec(sqlDropTables)
	if err != nil {
		return err
	}

	return s.create()
}

type pageTx struct {
	tx  *sql.Tx
	url *sql.Stmt
	ref *sql.Stmt
}

func newPageTx(db *sql.DB) (*pageTx, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}

	url, err := tx.Prepare(sqlInsertUrl)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	ref, err := tx.Prepare(sqlInsertRef)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return &pageTx{
		tx:  tx,
		url: url,
		ref: ref,
	}, nil
}

func (tx *pageTx) InsertUrl(url string) (int64, error) {
	res, err := tx.url.Exec(url)
	if err != nil {
		tx.tx.Rollback()
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		tx.tx.Rollback()
		return id, err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		tx.tx.Rollback()
		return 0, err
	}

	if rows == 0 {
		tx.tx.Rollback()
		return 0, store.ErrDupRow
	}

	return id, nil
}

func (tx *pageTx) InsertRef(oid, rid int64) error {
	_, err := tx.ref.Exec(oid, rid)

	return err
}
