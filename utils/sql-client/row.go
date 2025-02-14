package sqlclient

import "database/sql"

type sqlRow struct {
	row *sql.Row
}

type row interface {
	Scan(...interface{}) error
	Err() error
}

func (r *sqlRow) Err() error {
	return r.row.Err()
}

func (r *sqlRow) Scan(destinations ...interface{}) error {
	return r.row.Scan(destinations...)
}
