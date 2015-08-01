package core

import (
	"net/http"
)

type HttpErrorer interface {
	error
	ErrCode() int
}

type httpError struct {
	Code    int    `json:"code"`
	Message string `json:"error"`
}

func (e httpError) Error() string {
	return e.Message
}

func (e httpError) ErrCode() int {
	return e.Code
}

func ErrToErrorer(err error) (int, HttpErrorer) {
	var (
		httpErr httpError
		ok      bool
	)

	if httpErr, ok = err.(httpError); !ok {
		httpErr = httpError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return httpErr.ErrCode(), httpErr
}

func NewHttpError(code int, message string) httpError {
	return httpError{
		Code:    code,
		Message: message,
	}
}

func NewNotFoundError(message string) httpError {
	return httpError{
		Code:    http.StatusNotFound,
		Message: message,
	}
}

func NewForbiddenError() httpError {
	return httpError{
		Code:    http.StatusForbidden,
		Message: http.StatusText(http.StatusForbidden),
	}
}

func NewBadRequestError(message string) httpError {
	return httpError{
		Code:    http.StatusBadRequest,
		Message: message,
	}
}

func NewInternalServerError(message string) httpError {
	return httpError{
		Code:    http.StatusInternalServerError,
		Message: message,
	}
}
