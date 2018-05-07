package httputil

import (
	"fmt"
	"net/http"

	"github.com/go-openapi/runtime"
)

// Error is an http error.
// It implements both the error interface, so it can be returned as an error
// and also the WriteResponser interface so it can be used as a response for
// go-swagger applications.
type Error struct {
	Code    int    `json:"-"`
	Message string `json:"message"`
}

// Error implements the error interface
func (e Error) Error() string {
	return fmt.Sprintf("%s: %s", http.StatusText(e.Code), e.Message)
}

// WriteResponse builds the HTTP response with the necessary error Code
func (e Error) WriteResponse(w http.ResponseWriter, producer runtime.Producer) {
	w.WriteHeader(e.Code)
	if err := producer.Produce(w, e); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

// NewError returns a new error with Code and Message
func NewError(code int, msg string, args ...interface{}) Error {
	return Error{Code: code, Message: fmt.Sprintf(msg, args...)}
}

// NewErrNotFound create a new error which is mapped to http.StatusNotFound
func NewErrNotFound(msg string, args ...interface{}) Error {
	return NewError(http.StatusNotFound, msg, args...)
}

// NewErrConflict create a new error which is mapped to http.StatusConflict
func NewErrConflict(msg string, args ...interface{}) Error {
	return NewError(http.StatusConflict, msg, args...)
}

// NewErrInternalServer create a new error which is mapped to http.StatusInternalServerError
func NewErrInternalServer(msg string, args ...interface{}) Error {
	return NewError(http.StatusInternalServerError, msg, args...)
}

// NewErrBadRequest create a new error which is mapped to http.StatusInternalServerError
func NewErrBadRequest(msg string, args ...interface{}) Error {
	return NewError(http.StatusBadRequest, msg, args...)
}
