package util

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func ConvertErrorToString(oldErr *Error) string {
	errString, err := json.Marshal(oldErr)
	if err != nil {
		fmt.Printf("Error on Marshalling %v", errString)
		return ""
	}
	return string(errString)
}

func NewError(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

func HandleError(w http.ResponseWriter, code int, message string) {
	// Create an error response object
	errResp := Error{
		Code:    code,
		Message: message,
	}

	// Marshal the error response object into JSON
	errJSON, err := json.Marshal(errResp)
	if err != nil {
		http.Error(w, "Failed to marshal error response", code)
		return
	}

	// Set the content type header
	w.Header().Set("Content-Type", "application/json")

	// Set the HTTP status code for the error
	w.WriteHeader(code)

	// Write the JSON response to the response writer
	w.Write(errJSON)
}
