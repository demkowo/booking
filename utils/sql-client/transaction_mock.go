package sqlclient

import (
	"database/sql"
	"log"
)

type sqlTxMock struct {
	Tx    [][]interface{}
	Error error
	tx    *sql.Tx
}

func (t *sqlTxMock) Commit() error {
	log.Println("--- transactions/sqlclient/Commit() ---")
	return t.tx.Commit()
}

func (t *sqlTxMock) Exec(query string, args ...interface{}) (sql.Result, error) {
	log.Println("--- transactions/sqlclient/Exec() ---")

	return t.tx.Exec(query, args...)
}

func (t *sqlTxMock) Rollback() error {
	return t.tx.Rollback()
}
