package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	awslambda "github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
	"github.com/sebboness/yektaspoints/handlers"
	"github.com/sebboness/yektaspoints/handlers/family"
	userHandlers "github.com/sebboness/yektaspoints/handlers/user"
	"github.com/sebboness/yektaspoints/handlers/userauth"
	"github.com/sebboness/yektaspoints/util/env"
	"github.com/sebboness/yektaspoints/util/log"
)

var authCtrl *userauth.UserAuthController
var familyCtrl *family.FamilyController
var lambdaCtrl *handlers.LambdaController
var userCtrl *userHandlers.UserController

var ginLambda *ginadapter.GinLambda
var logger *log.Logger

// Handler is the main entry point for Lambda. Receives a proxy request and
// returns a proxy response
func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	logger = log.NewLogger("mypoints_lambda")

	_env := env.GetEnv("ENV")

	logger.WithContext(ctx).AddFields(map[string]any{
		"env":              _env,
		"method":           req.HTTPMethod,
		"path":             req.Path,
		"path_parameters":  req.PathParameters,
		"query_parameters": req.QueryStringParameters,
		"request_id":       req.RequestContext.RequestID,
		"authorizer":       req.RequestContext.Authorizer,
	}).Infof("starting lambda")

	// initialize auth user controller
	if authCtrl == nil {
		logger.Infof("initializing new user controller")
		_c, err := userauth.NewUserAuthController(ctx, _env)
		if err != nil {
			logger.Fatalf("failed to initialize user auth controller: %v", err)
		}

		authCtrl = _c
	}

	// initialize family controller
	if familyCtrl == nil {
		logger.Infof("initializing new family controller")
		_c, err := family.NewFamilyController(ctx, _env)
		if err != nil {
			logger.Fatalf("failed to initialize family controller: %v", err)
		}

		familyCtrl = _c
	}

	// intialize catchall lambda controller
	if lambdaCtrl == nil {
		logger.Infof("initializing new lambda controller")
		_c, err := handlers.NewLambdaController(ctx, _env)
		if err != nil {
			logger.Fatalf("failed to initialize lambda controller: %v", err)
		}

		lambdaCtrl = _c
	}

	// initialize user controller
	if userCtrl == nil {
		logger.Infof("initializing new user controller")
		_c, err := userHandlers.NewUserController(ctx, _env)
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

	// prepare context with authorizer info provided in lambda event
	ctx = handlers.PrepareAuthorizedContext(ctx, req)

	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	awslambda.Start(Handler)
}
