package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	awslambda "github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
	"github.com/sebboness/yektaspoints/handlers"
	"github.com/sebboness/yektaspoints/util/env"
	"github.com/sebboness/yektaspoints/util/log"
)

var c *handlers.LambdaController
var ginLambda *ginadapter.GinLambda
var logger = log.NewLogger("mypoints_lambda")

// Handler is the main entry point for Lambda. Receives a proxy request and
// returns a proxy response
func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	env := env.GetEnv("ENV")

	logger.WithContext(ctx).WithFields(map[string]any{
		"env":              env,
		"method":           req.HTTPMethod,
		"path":             req.Path,
		"path_parameters":  req.PathParameters,
		"query_parameters": req.QueryStringParameters,
		"request_id":       req.RequestContext.RequestID,
	}).Infof("starting lambda")

	if c == nil {
		logger.Infof("initializing new lambda controller")

		_c, err := handlers.NewLambdaController(env)
		if err != nil {
			logger.Fatalf("failed to initialize lambda controller: %v", err)
		}

		c = _c
	}

	if ginLambda == nil {
		logger.Infof("gin cold start")
		r := gin.Default()

		// Auth
		r.POST("/auth/token", c.UserAuthHandler)

		// Health
		r.GET("/", c.HealthCheckHandler)
		r.GET("/health", c.HealthCheckHandler)
		r.GET("/v1/health", c.HealthCheckHandler)

		// Points
		r.GET("/v1/points", c.GetUserPointsHandler)
		r.GET("/v1/points/:point_id", c.GetUserPointsHandler)
		r.POST("/v1/points", c.RequestPointsHandler)

		ginLambda = ginadapter.New(r)
	}

	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	awslambda.Start(Handler)
}
