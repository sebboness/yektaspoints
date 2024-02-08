package request_points

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sebboness/yektaspoints/storage"
	"github.com/sebboness/yektaspoints/util/env"
	"github.com/sebboness/yektaspoints/util/log"
)

var logger = log.NewLogger("request_points_lambda")

func main() {
	env := env.GetEnv("ENV")
	storageCfg := storage.Config{Env: env}

	ctx := context.Background()
	logger.WithContext(ctx).WithField("env", env).Infof("Starting lambda")

	pointsDB, err := storage.NewDynamoDbStorage(storageCfg)
	if err != nil {
		logger.Fatalf("failed to initialize points db: %v", err)
	}

	c := RequestPointsController{
		pointsDB: pointsDB,
	}

	lambda.Start(c.RequestPointsHandler)
}
