package errs

import (
	"encoding/json"
	"log"
)

type Error struct {
	Message string        `json:"message"`
	Code    int           `json:"code"`
	Status  string        `json:"status"`
	Causes  []interface{} `json:"causes"`
}

const (
	NotFoundErrorCode = 404
)

func NewError(message string, code int, status string, causes []interface{}) *Error {
	log.Println("--- errs/NewError() ---")

	return &Error{
		Message: message,
		Code:    code,
		Status:  status,
		Causes:  causes,
	}
}

func Err(err *Error) string {
	log.Println("--- errs/Err() ---")

	jsonBytes, _ := json.Marshal(err)
	return string(jsonBytes)
}
