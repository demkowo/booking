package sqlclient

import (
	"database/sql"
	"errors"
	"log"
	"os"
)

var (
	isMocked bool
)

type clientMock struct {
	mocks map[string]Mock
}

type Mock struct {
	Query   string
	Args    []interface{}
	Error   error
	Columns []string
	Rows    [][]interface{}
}

func (c *clientMock) Close() {
}

func (c *clientMock) Query(query string, args ...interface{}) (rows, error) {
	mock, exist := c.mocks[query]
	if !exist {
		return nil, errors.New("mock not found")
	}

	if mock.Error != nil {
		return nil, mock.Error
	}

	rows := rowsMock{
		Columns: mock.Columns,
		Rows:    mock.Rows,
	}

	return &rows, nil
}

func (c *clientMock) Exec(query string, args ...interface{}) (sql.Result, error) {
	log.Println("--- client_mock/sqlclient/Exec() ---")

	return nil, nil
}

func (c *clientMock) QueryRow(query string, args ...interface{}) row {
	mock := c.mocks[query]

	row := rowsMock{
		Columns: mock.Columns,
		Rows:    mock.Rows,
	}

	return &row
}

func isProduction() bool {
	log.Println("--- client/sqlclient/isProduction() ---")
	return os.Getenv(env) == production
}

func StartMockServer() {
	isMocked = true
}

func StopMockServer() {
	isMocked = false
}

func AddMock(mock Mock) {
	if dbClient == nil {
		return
	}
	client, expectedType := dbClient.(*clientMock)
	if !expectedType {
		return
	}
	if client.mocks == nil {
		client.mocks = make(map[string]Mock, 0)
	}
	client.mocks[mock.Query] = mock
}
