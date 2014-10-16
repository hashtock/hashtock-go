package http_utils

import (
    "net/http"
)

type HttpErrorer interface {
    error
    ErrCode() int
}

type HttpError struct {
    Code    int    `json:"code"`
    Message string `json:"error"`
}

func (e HttpError) Error() string {
    return e.Message
}

func (e HttpError) ErrCode() int {
    return e.Code
}

func NewHttpError(code int, message string) HttpError {
    return HttpError{
        Code:    code,
        Message: message,
    }
}

func NewNotFoundError(message string) HttpError {
    return HttpError{
        Code:    http.StatusNotFound,
        Message: message,
    }
}

func NewInternalServerError(message string) HttpError {
    return HttpError{
        Code:    http.StatusInternalServerError,
        Message: message,
    }
}
