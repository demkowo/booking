package resp

import (
	"encoding/json"
	"log"
)

type Response struct {
	Success bool          `json:"success"`
	Message string        `json:"message"`
	Code    int           `json:"code"`
	Status  string        `json:"status"`
	Data    interface{}   `json:"data,omitempty"`
	Causes  []interface{} `json:"causes,omitempty"`
}

func Error(message string, code int, status string, causes []interface{}) *Response {
	log.Println("--- resp/Error() ---")
	return &Response{
		Success: false,
		Message: message,
		Code:    code,
		Status:  status,
		Causes:  causes,
	}
}

func Success(message string, code int, status string, data interface{}) *Response {
	log.Println("--- resp/Success() ---")
	return &Response{
		Success: true,
		Message: message,
		Code:    code,
		Status:  status,
		Data:    data,
	}
}

func JSON(r *Response) string {
	jsonBytes, err := json.Marshal(r)
	if err != nil {
		log.Println("Error marshaling response:", err)
		return ""
	}
	return string(jsonBytes)
}
