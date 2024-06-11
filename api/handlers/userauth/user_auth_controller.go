package userauth

import (
	"context"
	"fmt"

	"github.com/sebboness/yektaspoints/handlers"
	"github.com/sebboness/yektaspoints/util/auth"
)

type UserAuthController struct {
	handlers.BaseController
	auth auth.AuthController
}

func NewUserAuthController(ctx context.Context, env string) (*UserAuthController, error) {
	authController, err := auth.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize auth controller: %w", err)
	}

	authContext, err := handlers.GetAuthContext()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize auth context: %w", err)
	}

	return &UserAuthController{
		BaseController: handlers.BaseController{
			AuthContext: authContext,
		},
		auth: authController,
	}, nil
}
