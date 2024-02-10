package error

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Error_NewApiError(t *testing.T) {

	err := NewApiError("failure")
	assert.Error(t, err.Err, errors.New("failure"))
}

func Test_Error_Error(t *testing.T) {
	type state struct {
		msg string
		err string
	}
	type want struct {
		result string
	}
	type test struct {
		name string
		state
		want
	}

	cases := []test{
		{
			"error only",
			state{err: "failed to validate"},
			want{result: "failed to validate"},
		},
		{
			"error with additional messages",
			state{err: "failed to validate", msg: "missing username"},
			want{result: "failed to validate: missing username"},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := NewApiError(c.state.err)
			err.AppendError(c.state.msg)

			result := err.Error()
			assert.Equal(t, result, c.want.result)
		})
	}
}

func Test_Error_AppendError(t *testing.T) {
	type state struct {
		errors []string
	}
	type want struct {
		result string
	}
	type test struct {
		name string
		state
		want
	}

	cases := []test{
		{
			"one error",
			state{errors: []string{"missing name"}},
			want{result: "invalid input: missing name"},
		},
		{
			"two errors",
			state{errors: []string{"missing name", "missing email"}},
			want{result: "invalid input: missing name; missing email"},
		},
		{
			"four errors",
			state{errors: []string{"missing name", "missing email", "drei", "vier"}},
			want{result: "invalid input: missing name; missing email; drei; vier"},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := NewApiError("invalid input")

			for _, e := range c.state.errors {
				err.AppendError(e)
			}

			result := err.Error()
			assert.Equal(t, result, c.want.result)
		})
	}
}

func Test_Error_determineStatusCode(t *testing.T) {
	type state struct {
		err    error
		errors []string
	}
	type want struct {
		code int
	}
	type test struct {
		name string
		state
		want
	}

	cases := []test{
		{"no error", state{}, want{http.StatusInternalServerError}},
		{"bad request", state{err: BadRequest}, want{http.StatusBadRequest}},
		{"invalid input", state{err: InvalidInput}, want{http.StatusBadRequest}},
		{"unauthorized", state{err: Unauthorized}, want{http.StatusUnauthorized}},
		{"not found", state{err: NotFound}, want{http.StatusNotFound}},
		{"internal server error", state{err: InternalServerError}, want{http.StatusInternalServerError}},
		{"non-nil error", state{err: errors.New("fail")}, want{http.StatusBadRequest}},
		{"non-empty errors", state{errors: []string{"fail!"}}, want{http.StatusBadRequest}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := New(c.state.err)

			if len(c.state.errors) > 0 {
				err.errors = c.state.errors
			}

			result := err.determineStatusCode()
			assert.Equal(t, c.want.code, result)
		})
	}
}

func Test_Error_StatusCode(t *testing.T) {
	type state struct {
		hasErr      bool
		hasNoErrors bool
		statusCode  int
	}
	type want struct {
		code int
	}
	type test struct {
		name string
		state
		want
	}

	cases := []test{
		{"no error - some errors", state{}, want{http.StatusBadRequest}},
		{"has error - some errors", state{hasErr: true}, want{http.StatusBadRequest}},
		{"has error - no errors", state{hasNoErrors: true}, want{http.StatusBadRequest}},
		{"custom status code", state{statusCode: http.StatusForbidden}, want{http.StatusForbidden}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := NewApiError("wah").WithStatus(c.state.statusCode)

			if c.state.hasErr {
				err.Err = errors.New("fail")
			}
			if c.state.hasNoErrors {
				err.errors = []string{}
			}

			result := err.StatusCode()
			assert.Equal(t, c.want.code, result)
		})
	}
}

func Test_Error_WithError(t *testing.T) {
	err := New(NotFound).WithError("test")
	assert.Equal(t, "resource not found: test", err.Error())
}

func Test_Error_AsError(t *testing.T) {
	err := New(NotFound)
	var apiErr *ApiError
	assert.ErrorAs(t, err, &apiErr)
}

func Test_Error_IsApiError(t *testing.T) {
	type state struct {
		err error
	}
	type want struct {
		shouldBeApiErr bool
	}
	type test struct {
		name string
		state
		want
	}

	cases := []test{
		{"is api error #1", state{NewApiError("something went wrong")}, want{true}},
		{"is api error #2", state{New(InvalidInput)}, want{true}},
		{"is not api error", state{errors.New("something else went wrong")}, want{}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			apiErr := IsApiError(c.state.err)
			if c.want.shouldBeApiErr {
				assert.NotNil(t, apiErr)
			} else {
				assert.Nil(t, apiErr)
			}
		})
	}
}
