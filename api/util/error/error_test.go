package error

import (
	"errors"
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
		err := NewApiError(c.state.err)
		err.AppendError(c.state.msg)

		result := err.Error()
		assert.Equal(t, result, c.want.result)
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
		err := NewApiError("invalid input")

		for _, e := range c.state.errors {
			err.AppendError(e)
		}

		result := err.Error()
		assert.Equal(t, result, c.want.result)
	}
}

func Test_Error_IsError(t *testing.T) {
	err := New(NotFound, "item (123) not found")
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
		{"is api error #2", state{New(InvalidInput, "failed to validate the request")}, want{true}},
		{"is not api error", state{errors.New("something else went wrong")}, want{}},
	}

	for _, c := range cases {
		apiErr := IsApiError(c.state.err)
		if c.want.shouldBeApiErr {
			assert.NotNil(t, apiErr)
		} else {
			assert.Nil(t, apiErr)
		}
	}
}
