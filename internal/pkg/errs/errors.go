package errs

import (
	"errors"
	"net/http"
)

type AppError struct {
	Code    int
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}

func NewNotFound(message string) *AppError {
	return &AppError{Code: http.StatusNotFound, Message: message}
}

func NewConflict(message string) *AppError {
	return &AppError{Code: http.StatusConflict, Message: message}
}

func NewBadRequest(message string) *AppError {
	return &AppError{Code: http.StatusBadRequest, Message: message}
}

func NewUnauthorized(message string) *AppError {
	return &AppError{Code: http.StatusUnauthorized, Message: message}
}

func NewForbidden(message string) *AppError {
	return &AppError{Code: http.StatusForbidden, Message: message}
}

func NewInternal(message string) *AppError {
	return &AppError{Code: http.StatusInternalServerError, Message: message}
}

func ToAppError(err error) *AppError {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}
	return NewInternal("internal server error")
}

