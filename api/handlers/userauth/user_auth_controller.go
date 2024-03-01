package userauth

import (
	"context"
	"fmt"

	"github.com/sebboness/yektaspoints/util/auth"
)

type UserAuthController struct {
	auth auth.AuthController
}

func NewUserAuthController(ctx context.Context, env string) (*UserAuthController, error) {
	authController, err := auth.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize auth controller: %w", err)
	}

	return &UserAuthController{
		auth: authController,
	}, nil
}
