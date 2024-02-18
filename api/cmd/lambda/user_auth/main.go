package user_auth

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sebboness/yektaspoints/handlers"
	"github.com/sebboness/yektaspoints/util/env"
	"github.com/sebboness/yektaspoints/util/log"
)

var logger = log.NewLogger("user_auth_lambda")

func main() {
	env := env.GetEnv("ENV")

	ctx := context.Background()
	logger.WithContext(ctx).WithField("env", env).Infof("Starting lambda")

	c, err := handlers.NewLambdaController(env)
	if err != nil {
		logger.Fatalf("failed to initialize lambda controller: %v", err)
	}

	lambda.Start(c.UserAuthHandler)
}
