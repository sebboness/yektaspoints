package points

import (
	"context"
	"errors"
	"testing"

	"github.com/sebboness/yektaspoints/util/tests"
	"github.com/stretchr/testify/assert"
)

var errFail = errors.New("fail")

func Test_Controller_New(t *testing.T) {
	c, err := NewPointsController(context.Background(), "local")
	tests.AssertError(t, err, "")
	assert.NotNil(t, c)
}
