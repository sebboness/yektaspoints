package handlers

import (
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/sebboness/yektaspoints/util/env"
	"github.com/stretchr/testify/assert"
)

func Test_HandleHealthCheck(t *testing.T) {
	c, err := NewLambdaController(env.GetEnv("ENV"))
	if err != nil {
		panic("failed to initialize lambda controller: " + err.Error())
	}

	resp, err := c.HandleHealthCheck(context.Background(), &events.APIGatewayProxyRequest{})
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	assert.Contains(t, resp.Body, `"env":"local"`)
}
