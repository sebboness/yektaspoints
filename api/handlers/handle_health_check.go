package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sebboness/yektaspoints/util/env"
)

type healthCheckResponse struct {
	Env     string `json:"env"`
	Version string `json:"version"`
	BuiltAt string `json:"build_at"`
}

// HealthCheckHandler is used to return a simple health check
// func (c *LambdaController) HealthCheckHandler(ctx context.Context, event *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
func (c *LambdaController) HealthCheckHandler(cgin *gin.Context) {
	resp := healthCheckResponse{
		Env:     env.GetEnv("ENV"),
		Version: env.GetEnv("VERSION"),
		BuiltAt: env.GetEnv("BUILT_AT"),
	}

	logger.WithContext(cgin.Request.Context()).Infof("hello from health check endpoint")

	cgin.JSON(http.StatusOK, SuccessResult(resp))
}
