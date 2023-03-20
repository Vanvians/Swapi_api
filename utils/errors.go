package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

//Represents an error with the specific HTTP status code
type ApiError struct{
	Code int `json:"code"`
	Message string `json:"message"`
}

//Creates a new instance if ApiError with the HTTP status code
func NewApiError(code int, message string) *ApiError{
	return &ApiError{
		Code: code,
		Message: message,
	}
}

//return string representation of ApiError
func (e *ApiError) Error() string{
	errmsg := fmt.Sprintf("ApiError{code=%d, message=%s}", e.Code, e.Message)
	return errmsg
}

//writes data to an HTTP respone in JSON format
func WriteJson(w http.ResponseWriter, code int, data interface{}){
	w.Header().Set("Content-Type", "apllication/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}

//writes an APiError to a HTTP response
func WriteError(w http.ResponseWriter, err *ApiError){
	WriteJson(w,err.Code, err)
}