package shared

import (
	"errors"
	"fmt"
	"net/http"
)

type AppError struct {
	Code    int    `json:"-"`
	Message string `json:"error"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err == nil {
		return e.Message
	}
	return fmt.Sprintf("%s: %v", e.Message, e.Err)
}

func (e *AppError) Unwrap() error { return e.Err }

// Constructors

func BadRequest(msg string, err error) *AppError {
	return &AppError{Code: http.StatusBadRequest, Message: msg, Err: err}
}
func NotFound(msg string, err error) *AppError {
	return &AppError{Code: http.StatusNotFound, Message: msg, Err: err}
}
func Conflict(msg string, err error) *AppError {
	return &AppError{Code: http.StatusConflict, Message: msg, Err: err}
}
func Internal(msg string, err error) *AppError {
	return &AppError{Code: http.StatusInternalServerError, Message: msg, Err: err}
}

func UnAuth(msg string, err error) *AppError {
	return &AppError{Code: http.StatusUnauthorized, Message: msg, Err: err}
}

func IsAppError(err error) (*AppError, bool) {
	var ae *AppError
	if errors.As(err, &ae) {
		return ae, true
	}
	return nil, false
}
