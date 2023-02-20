package errs

import (
	"errors"
	"log"
	"net/http"
)

type ApiError struct {
	Err        error
	StatusCode int
}

func (e ApiError) Error() string {
	return e.Err.Error()
}

func (e ApiError) Status() int {
	return e.StatusCode
}

func NewInternal(err error) IApiError {
	return newApiError(err, http.StatusInternalServerError)
}

func NewBadRequest(err error) IApiError {
	return newApiError(err, http.StatusBadRequest)
}

func NewNotFound(err error) IApiError {
	return newApiError(err, http.StatusNotFound)
}

func NewUnauthorized(err error) IApiError {
	return newApiError(err, http.StatusUnauthorized)
}

func NewForbidden(err error) IApiError {
	return newApiError(err, http.StatusForbidden)
}

func NewNotImplemented() IApiError {
	return newApiError(errors.New("not implemented"), http.StatusNotImplemented)
}

func newApiError(err error, statusCode int) ApiError {
	log.Println(err.Error())
	return ApiError{
		Err:        err,
		StatusCode: statusCode,
	}
}
