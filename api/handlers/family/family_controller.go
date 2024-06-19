package family

import (
	"context"
	"fmt"

	"github.com/sebboness/yektaspoints/handlers"
	"github.com/sebboness/yektaspoints/storage"
)

type FamilyController struct {
	handlers.BaseController
	familyDB storage.IFamilyStorage
}

func NewFamilyController(ctx context.Context, env string) (*FamilyController, error) {
	storageCfg := storage.Config{Env: env}

	familyDB, err := storage.NewDynamoDbStorage(storageCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize family db: %w", err)
	}

	authContext, err := handlers.GetAuthContext()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize auth context: %w", err)
	}

	return &FamilyController{
		BaseController: handlers.BaseController{
			AuthContext: authContext,
		},
		familyDB: familyDB,
	}, nil
}
