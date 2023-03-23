package utils

import (
	"fmt"
	"net/http"
	"encoding/json"
)

// APIError represents an error returned by an API.
type APIError struct {
	Message string `json:"message"`
}

// Error returns the error message.
func (e *APIError) Error() string {
	return e.Message
}

// NewAPIError returns a new instance of APIError with the given message.
func NewAPIError(message string) *APIError {
	return &APIError{Message: message}
}

// DBError represents an error related to a database operation.
type DBError struct {
	Message string `json:"message"`
}

// Error returns the error message.
func (e *DBError) Error() string {
	return e.Message
}

// NewDBError returns a new instance of DBError with the given message.
func NewDBError(message string) *DBError {
	return &DBError{Message: message}
}

// ValidationError represents an error caused by invalid input.
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Error returns the error message.
func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// NewValidationError returns a new instance of ValidationError with the given field and message.
func NewValidationError(field string, message string) *ValidationError {
	return &ValidationError{Field: field, Message: message}
}

func RespondWithError(w http.ResponseWriter,status int, err error){
	w.WriteHeader(status)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "error": err.Error(),
    })
}

func RespondWithJSON(w http.ResponseWriter, status int, data interface{}) {
    w.WriteHeader(status)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(data)
}
