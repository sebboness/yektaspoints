package handlers

import (
	"context"
	"fmt"

	"github.com/sebboness/yektaspoints/storage"
	"github.com/sebboness/yektaspoints/util/auth"
)

type LambdaController struct {
	BaseController
	auth     auth.AuthController
	pointsDB storage.IPointsStorage
}

func NewLambdaController(ctx context.Context, env string) (*LambdaController, error) {
	storageCfg := storage.Config{Env: env}

	pointsDB, err := storage.NewDynamoDbStorage(storageCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize points db: %w", err)
	}

	authController, err := auth.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize auth controller: %w", err)
	}

	authContext, err := GetAuthContext()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize auth context: %w", err)
	}

	return &LambdaController{
		BaseController: BaseController{
			AuthContext: authContext,
		},
		auth:     authController,
		pointsDB: pointsDB,
	}, nil
}
