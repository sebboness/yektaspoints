package points

import (
	"context"
	"fmt"

	"github.com/sebboness/yektaspoints/storage"
)

type PointsController struct {
	pointsDB storage.IPointsStorage
}

func NewPointsController(ctx context.Context, env string) (*PointsController, error) {
	storageCfg := storage.Config{Env: env}

	userDB, err := storage.NewDynamoDbStorage(storageCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize points db: %w", err)
	}

	return &PointsController{
		pointsDB: userDB,
	}, nil
}
