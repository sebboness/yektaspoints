package user

import (
	"context"
	"fmt"

	"github.com/sebboness/yektaspoints/storage"
	"github.com/sebboness/yektaspoints/util/auth"
)

type UserController struct {
	auth   auth.AuthController
	userDB storage.IUserStorage
}

func NewUserController(ctx context.Context, env string) (*UserController, error) {
	storageCfg := storage.Config{Env: env}

	userDB, err := storage.NewDynamoDbStorage(storageCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize user db: %w", err)
	}

	authController, err := auth.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize auth controller: %w", err)
	}

	return &UserController{
		auth:   authController,
		userDB: userDB,
	}, nil
}
