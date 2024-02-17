package handlers

import (
	"context"
	"fmt"

	"github.com/sebboness/yektaspoints/storage"
	"github.com/sebboness/yektaspoints/util/auth"
)

type LambdaController struct {
	auth     auth.AuthController
	pointsDB storage.IPointsStorage
}

func NewLambdaController(env string) (*LambdaController, error) {
	storageCfg := storage.Config{Env: env}

	pointsDB, err := storage.NewDynamoDbStorage(storageCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize points db: %w", err)
	}

	authController, err := auth.New(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("failed to initialize auth controller: %w", err)
	}

	return &LambdaController{
		auth:     authController,
		pointsDB: pointsDB,
	}, nil
}
