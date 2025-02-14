package sqlclient

import (
	"database/sql"
	"errors"
	"log"
)

const (
	env        = "ENV"
	production = "prod"
)

var (
	dbClient SqlClient
)

type client struct {
	db *sql.DB
}

type SqlClient interface {
	Close()
	Exec(string, ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (rows, error)
	QueryRow(query string, args ...interface{}) row
}

func (c *client) Close() {
	c.db.Close()
}

func (c *client) Exec(query string, args ...interface{}) (sql.Result, error) {
	log.Println("--- client/sqlclient/Exec() ---")

	res, err := c.db.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *client) Query(query string, args ...interface{}) (rows, error) {
	log.Println("--- client/sqlclient/Query() ---")

	res, err := c.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	rows := sqlRows{
		rows: res,
	}

	return &rows, nil
}

func (c *client) QueryRow(query string, args ...interface{}) row {
	log.Println("--- client/sqlclient/QueryRow() ---")
	res := c.db.QueryRow(query, args...)

	row := sqlRow{
		row: res,
	}

	return &row
}

func Open(driverName, dataSourceName string) (SqlClient, error) {
	log.Println("--- client/sqlclient/Open() ---")

	if !isProduction() && isMocked {
		dbClient = &clientMock{}
		return dbClient, nil
	}
	if driverName == "" {
		return nil, errors.New("invalid driver name")
	}

	database, err := sql.Open(driverName, dataSourceName)
	if driverName == "" {
		return nil, err
	}

	dbClient := &client{
		db: database,
	}

	return dbClient, nil
}
