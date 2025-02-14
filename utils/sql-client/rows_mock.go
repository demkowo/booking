package sqlclient

import "errors"

type rowsMock struct {
	Columns []string
	Rows    [][]interface{}

	Error error
	Index int
}

func (r *rowsMock) Next() bool {
	return r.Index < len(r.Rows)
}

func (r *rowsMock) Close() error {
	return nil
}

func (r *rowsMock) Scan(destinations ...interface{}) error {
	row := r.Rows[r.Index]
	if len(row) != len(destinations) {
		return errors.New("invalid destination length")
	}

	for index, value := range row {
		destinations[index] = value
	}

	return nil
}

func (r *rowsMock) Err() error {
	return r.Error
}
