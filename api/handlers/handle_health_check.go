package handlers

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/sebboness/yektaspoints/util/env"
)

// HandleHealthCheck is used to return a simple health check
func (c *LambdaController) HandleHealthCheck(ctx context.Context, event *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	resp := struct {
		Env     string `json:"env"`
		Version string `json:"version"`
		BuiltAt string `json:"build_at"`
	}{
		Env:     env.GetEnv("ENV"),
		Version: env.GetEnv("VERSION"),
		BuiltAt: env.GetEnv("BUILT_AT"),
	}

	logger.WithContext(ctx).Infof("hello from health check endpoint")

	return ApiResponseOK(resp), nil
}
