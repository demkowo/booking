package sqlclient

import "database/sql"

type sqlRows struct {
	rows *sql.Rows
}

type rows interface {
	Next() bool
	Close() error
	Scan(...interface{}) error
	Err() error
}

func (r *sqlRows) Next() bool {
	return r.rows.Next()
}

func (r *sqlRows) Close() error {
	return r.rows.Close()
}

func (r *sqlRows) Scan(destinations ...interface{}) error {
	return r.rows.Scan(destinations...)
}

func (r *sqlRows) Err() error {
	return r.rows.Err()
}
