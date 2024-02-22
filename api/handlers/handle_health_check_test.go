package handlers

import (
	"context"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sebboness/yektaspoints/util/env"
	"github.com/stretchr/testify/assert"
)

func Test_Controller_HealthCheckHandler(t *testing.T) {
	c, err := NewLambdaController(context.Background(), env.GetEnv("ENV"))
	if err != nil {
		panic("failed to initialize lambda controller: " + err.Error())
	}

	w := httptest.NewRecorder()
	cgin, _ := gin.CreateTestContext(w)
	cgin.Request = httptest.NewRequest("GET", "/health", nil)

	c.HealthCheckHandler(cgin)

	assert.Equal(t, 200, w.Code)

	body, err := io.ReadAll(w.Body)
	assert.Nil(t, err, "response body read should have no error")
	assert.Contains(t, string(body), `"env":"local"`)
}
