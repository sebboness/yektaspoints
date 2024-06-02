package points

import (
	"context"
	"fmt"

	"github.com/sebboness/yektaspoints/storage"
	"github.com/sebboness/yektaspoints/util/log"
)

var logger = log.Get()

type PointsController struct {
	pointsDB storage.IPointsStorage
	userDB   storage.IUserStorage
}

func NewPointsController(ctx context.Context, env string) (*PointsController, error) {
	storageCfg := storage.Config{Env: env}

	db, err := storage.NewDynamoDbStorage(storageCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize db: %w", err)
	}

	return &PointsController{
		pointsDB: db,
		userDB:   db,
	}, nil
}
