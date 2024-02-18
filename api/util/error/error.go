package error

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// Custom errors
var (
	BadRequest          = errors.New("bad request")
	InvalidInput        = errors.New("invalid input")
	Unauthorized        = errors.New("unauthorized")
	NotFound            = errors.New("resource not found")
	InternalServerError = errors.New("internal server error")
)

type ApiError struct {
	Err        error
	errors     []string
	statusCode int
}

func New(err error) *ApiError {
	return &ApiError{
		Err:    err,
		errors: []string{},
	}
}

func NewApiError(err string) *ApiError {
	return &ApiError{
		Err:    errors.New(err),
		errors: []string{},
	}
}

func (e *ApiError) WithError(err string) *ApiError {
	e.AppendError(err)
	return e
}

func (e *ApiError) WithStatus(statusCode int) *ApiError {
	e.statusCode = statusCode
	return e
}

func (e *ApiError) Error() string {
	joinedErrors := strings.Join(e.errors, "; ")
	if len(joinedErrors) > 0 {
		joinedErrors = ": " + joinedErrors
	}

	if e.Err != nil {
		return fmt.Sprintf("%v%s", e.Err, joinedErrors)
	}

	return joinedErrors
}

func (e *ApiError) Errors() []string {
	return e.errors
}

func (e *ApiError) StatusCode() int {
	if e.statusCode == 0 {
		e.statusCode = e.determineStatusCode()
	}
	return e.statusCode
}

func (e *ApiError) AppendError(errors ...string) {
	e.errors = append(e.errors, errors...)
}

func (e *ApiError) AppendErrorf(msg string, args ...any) {
	err := fmt.Sprintf(msg, args...)
	e.AppendError(err)
}

func (e *ApiError) Unwrap() error {
	return e.Err
}

func (e *ApiError) Is(err error) bool {
	return errors.Is(e.Err, err)
}

func (e *ApiError) determineStatusCode() int {
	if e.Is(BadRequest) || e.Is(InvalidInput) {
		return http.StatusBadRequest
	} else if e.Is(Unauthorized) {
		return http.StatusUnauthorized
	} else if e.Is(NotFound) {
		return http.StatusNotFound
	} else if e.Is(InternalServerError) {
		return http.StatusInternalServerError
	} else if e.Err != nil || len(e.errors) > 0 {
		return http.StatusBadRequest
	}
	return http.StatusInternalServerError
}

func IsApiError(err error) *ApiError {
	var apiErr *ApiError
	if errors.As(err, &apiErr) {
		return apiErr
	}
	return nil
}
