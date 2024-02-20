package lambda

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

	logger.WithContext(ctx).WithField("env", env).Infof("starting lambda")

	if c == nil {
		logger.Infof("initializing new lambda controller")

		_c, err := handlers.NewLambdaController(env)
		if err != nil {
			logger.Fatalf("failed to initialize lambda controller: %v", err)
		}

		c = _c
	}

	if ginLambda == nil {
		// stdout and stderr are sent to AWS CloudWatch Logs
		logger.Infof("gin cold start")
		r := gin.Default()

		// Health
		r.GET("/", c.HealthCheckHandler)
		r.GET("/health", c.HealthCheckHandler)

		// Auth
		r.POST("/auth/token", c.UserAuthHandler)

		// Points
		r.GET("/points", c.GetUserPointsHandler)
		r.GET("/points/:point_id", c.GetUserPointsHandler)
		r.POST("/points", c.RequestPointsHandler)

		ginLambda = ginadapter.New(r)
	}

	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	awslambda.Start(Handler)
}
