package error

import (
	"errors"
	"fmt"
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
	Errors     []string
	statusCode int
}

func New(err error, msg string) *ApiError {
	return &ApiError{
		Err:    err,
		Errors: []string{msg},
	}
}

func NewApiError(err string) *ApiError {
	return &ApiError{
		Err:    errors.New(err),
		Errors: []string{},
	}
}

func (e *ApiError) WithStatus(statusCode int) *ApiError {
	e.statusCode = statusCode
	return e
}

func (e *ApiError) Error() string {
	joinedErrors := strings.Join(e.Errors, "; ")
	if len(joinedErrors) > 0 {
		joinedErrors = ": " + joinedErrors
	}

	if e.Err != nil {
		return fmt.Sprintf("%v%s", e.Err, joinedErrors)
	}

	return joinedErrors
}

func (e *ApiError) StatusCode() int {
	return e.statusCode
}

func (e *ApiError) AppendError(errors ...string) {
	e.Errors = append(e.Errors, errors...)
}

func (e *ApiError) Is(err error) bool {
	return errors.Is(e.Err, err)
}

func (e *ApiError) Unwrap() error {
	return e.Err
}

func IsApiError(err error) *ApiError {
	var apiErr *ApiError
	if errors.As(err, &apiErr) {
		return apiErr
	}
	return nil
}
