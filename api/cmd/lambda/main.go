package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	awslambda "github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
	"github.com/sebboness/yektaspoints/handlers"
	userHandlers "github.com/sebboness/yektaspoints/handlers/user"
	"github.com/sebboness/yektaspoints/util/env"
	"github.com/sebboness/yektaspoints/util/log"
)

var lambdaCtrl *handlers.LambdaController
var userCtrl *userHandlers.UserController

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

	// intialize catchall lambda controller
	if lambdaCtrl == nil {
		logger.Infof("initializing new lambda controller")
		_c, err := handlers.NewLambdaController(env)
		if err != nil {
			logger.Fatalf("failed to initialize lambda controller: %v", err)
		}

		lambdaCtrl = _c
	}

	// initialize user controller
	if userCtrl == nil {
		logger.Infof("initializing new user controller")
		_c, err := userHandlers.NewUserController(env)
		if err != nil {
			logger.Fatalf("failed to initialize user controller: %v", err)
		}

		userCtrl = _c
	}

	if ginLambda == nil {
		logger.Infof("gin cold start")
		r := gin.Default()

		RegisterRoutes(r)

		ginLambda = ginadapter.New(r)
	}

	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	awslambda.Start(Handler)
}
