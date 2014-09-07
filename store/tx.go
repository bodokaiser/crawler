package store

import (
	"database/sql"
	"errors"
	"net/url"
	"strings"
)

var ErrRefExists = errors.New("reference exists already")

// Type Tx encapsulates a SQL transaction which sets up url relation for a
// single html site.
type Tx struct {
	origin  row
	refers  []row
	tx      *sql.Tx
	urlStmt *sql.Stmt
	refStmt *sql.Stmt
}

type row struct {
	Id   int64
	Url  string
	Done bool
}

var sqlSelectUrl = `
	SELECT id, url
	FROM url
	WHERE done = false
	LIMIT 1
	LOCK IN SHARE MODE
`

var sqlUpdateUrl = `
	UPDATE url
	SET done = ?
	WHERE id = ?
`

var sqlInsertUrl = `
	INSERT into url (url)
	VALUES (?)
	ON DUPLICATE KEY UPDATE id = LAST_INSERT_ID(id)
`

var sqlInsertRef = `
	INSERT into ref (origin_id, refer_id)
	VALUES (?, ?)
`

// Sets up a new url transaction by acquiring uncrawled site, locking it and
// preparing statements for insertion.
func newTx(tx *sql.Tx) (*Tx, error) {
	t := &Tx{
		tx:     tx,
		origin: row{},
		refers: make([]row, 0),
	}

	if err := t.selectUrl(&t.origin); err != nil {
		return nil, err
	}

	return t, nil
}

// Returns original url.
func (tx *Tx) Origin() string {
	return tx.origin.Url
}

// Returns referenced urls.
func (tx *Tx) Refers() []string {
	s := make([]string, len(tx.refers))

	for i, r := range tx.refers {
		s[i] = r.Url
	}

	return s
}

// Adds referenced url if not exists.
func (tx *Tx) AddRefer(r string) error {
	for _, ref := range tx.refers {
		if ref.Url == r {
			return ErrRefExists
		}
	}

	uri, err := tx.normalize(r)
	if err != nil {
		return err
	}

	row := row{Url: uri.String()}

	if err := tx.insertUrl(&row); err != nil {
		return err
	}
	if err := tx.insertRef(&row); err != nil {
		return err
	}

	tx.refers = append(tx.refers, row)

	return nil
}

func (tx *Tx) Abort() error {
	return tx.tx.Rollback()
}

// Commits updated to database.
func (tx *Tx) Commit() error {
	tx.origin.Done = true

	if err := tx.updateUrl(&tx.origin); err != nil {
		tx.tx.Rollback()

		return err
	}

	return tx.tx.Commit()
}

// Acquires an unlocked url row which was not crawled before.
// Maps data of that row to transaction and locks distinct access.
func (tx *Tx) selectUrl(r *row) error {
	err := tx.tx.QueryRow(sqlSelectUrl).Scan(&r.Id, &r.Url)

	if err != nil {
		tx.tx.Rollback()
	}

	return err
}

// Updates origin url to be marked as done.
// May include further meta data in future.
func (tx *Tx) updateUrl(r *row) error {
	_, err := tx.tx.Exec(sqlUpdateUrl, r.Id, r.Done)
	if err != nil {
		tx.tx.Rollback()

		return err
	}

	return err
}

// Executes prepared statement to insert url if not existent.
func (tx *Tx) insertUrl(r *row) error {
	if tx.urlStmt == nil {
		stmt, err := tx.tx.Prepare(sqlInsertUrl)
		if err != nil {
			tx.tx.Rollback()

			return err
		}

		tx.urlStmt = stmt
	}

	res, err := tx.urlStmt.Exec(r.Url)
	if err != nil {
		tx.tx.Rollback()

		return err
	}

	r.Id, err = res.LastInsertId()
	if err != nil {
		tx.tx.Rollback()
	}

	return err
}

// Executes prepared statement to insert url reference.
func (tx *Tx) insertRef(r *row) error {
	if tx.refStmt == nil {
		stmt, err := tx.tx.Prepare(sqlInsertRef)
		if err != nil {
			tx.tx.Rollback()

			return err
		}

		tx.refStmt = stmt
	}

	_, err := tx.refStmt.Exec(tx.origin.Id, r.Id)
	if err != nil {
		tx.tx.Rollback()
	}

	return err
}

// Returns normalized url string with lowercased host, scheme and removed
// fragments. Path stays untouched.
func (tx *Tx) normalize(ref string) (*url.URL, error) {
	uri, err := url.Parse(ref)
	if err != nil {
		return nil, err
	}

	uri.Host = strings.ToLower(uri.Host)
	uri.Scheme = strings.ToLower(uri.Scheme)
	uri.Fragment = ""

	return uri, nil
}
