package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
	"testing"

	"github.com/sebboness/yektaspoints/util/result"
	"github.com/stretchr/testify/assert"
)

func AssertError(t *testing.T, err error, wantErr string) {
	if err == nil {
		assert.Nil(t, err)
		assert.Empty(t, wantErr, "expected error should not be nil")
	} else {
		assert.NotNil(t, err)
		assert.NotEmpty(t, wantErr, "expected no error, but there was: "+err.Error())
		assert.Contains(t, err.Error(), wantErr)
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

func AssertResultError(t *testing.T, res result.Result, wantErr string) {
	if res.IsSuccess() {
		assert.Empty(t, wantErr, "expected no errors")
	} else {
		assert.NotEmpty(t, wantErr, "expected no error, but there was")
		hasErr := false
		for _, e := range res.Errors {
			if strings.Contains(e, wantErr) {
				hasErr = true
				break
			}
		}
		assert.True(t, hasErr, "expected result to contain error")
	}
}
