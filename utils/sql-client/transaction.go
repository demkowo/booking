package sqlclient

import (
	"database/sql"
	"log"
)

type sqlTx struct {
	tx *sql.Tx
}

type Tx interface {
	Commit() error
	Exec(string, ...interface{}) (sql.Result, error)
	Rollback() error
}

func (t *sqlTx) Commit() error {
	log.Println("--- transactions/sqlclient/Commit() ---")
	return t.tx.Commit()
}

func (t *sqlTx) Exec(query string, args ...interface{}) (sql.Result, error) {
	log.Println("--- transactions/sqlclient/Exec() ---")

	return t.tx.Exec(query, args...)
}

func (t *sqlTx) Rollback() error {
	return t.tx.Rollback()
}
