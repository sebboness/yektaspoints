package handlers

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_PointsController_New(t *testing.T) {
	c, err := NewLambdaController(context.Background(), "local")
	assert.Nil(t, err)
	assert.NotNil(t, c)
}
