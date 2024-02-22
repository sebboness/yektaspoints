package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"testing"

	"github.com/sebboness/yektaspoints/util/result"
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

func AssertResult(t *testing.T, b *bytes.Buffer) result.Result {
	body, err := io.ReadAll(b)
	assert.Nil(t, err, "there should be no error reading result body")

	var res result.Result
	err = json.Unmarshal(body, &res)
	assert.Nil(t, err, "there should be no error unmarshaling result from body")
	assert.NotEmpty(t, res.Status, "result status should not be empty")
	return res
}
