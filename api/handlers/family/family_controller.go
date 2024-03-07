package family

import (
	"context"
	"fmt"

	"github.com/sebboness/yektaspoints/storage"
)

type FamilyController struct {
	familyDB storage.IFamilyStorage
}

func NewFamilyController(ctx context.Context, env string) (*FamilyController, error) {
	storageCfg := storage.Config{Env: env}

	familyDB, err := storage.NewDynamoDbStorage(storageCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize family db: %w", err)
	}

	return &FamilyController{
		familyDB: familyDB,
	}, nil
}
