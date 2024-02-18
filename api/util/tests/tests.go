package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func AssertError(t *testing.T, err error, wantError string) {
	if err == nil {
		assert.Nil(t, err)
		assert.Empty(t, wantError, "expected error should not be nil")
	} else {
		assert.NotNil(t, err)
		assert.NotEmpty(t, wantError, "expected no error, but there was: "+err.Error())
		assert.Contains(t, err.Error(), wantError)
	}
}
