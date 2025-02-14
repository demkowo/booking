package sqlclient

import "errors"

type rowMock struct {
	Columns []string
	Row     [][]interface{}

	Error error
	Index int
}

func (r *rowMock) Next() bool {
	return r.Index < len(r.Row)
}

func (r *rowMock) Close() error {
	return nil
}

func (r *rowMock) Scan(destinations ...interface{}) error {
	row := r.Row[r.Index]
	if len(row) != len(destinations) {
		return errors.New("invalid destination length")
	}

	for index, value := range row {
		destinations[index] = value
	}

	return nil
}

func (r *rowMock) Err() error {
	return r.Error
}
