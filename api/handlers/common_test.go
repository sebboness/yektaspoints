package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_PointsController_New(t *testing.T) {
	c, err := NewPointsController("local")
	assert.Nil(t, err)
	assert.NotNil(t, c)
}