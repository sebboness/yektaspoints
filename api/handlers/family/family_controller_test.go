package family

import (
	"context"
	"testing"

	"github.com/sebboness/yektaspoints/util/tests"
	"github.com/stretchr/testify/assert"
)

func Test_Controller_New(t *testing.T) {
	c, err := NewFamilyController(context.Background(), "local")
	tests.AssertError(t, err, "")
	assert.NotNil(t, c)
}
