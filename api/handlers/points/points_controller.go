package points

import (
	"context"
	"fmt"

	"github.com/sebboness/yektaspoints/handlers"
	"github.com/sebboness/yektaspoints/storage"
	"github.com/sebboness/yektaspoints/util/log"
)

var logger = log.Get()

type PointsController struct {
	handlers.BaseController
	pointsDB storage.IPointsStorage
	userDB   storage.IUserStorage
}

func NewPointsController(ctx context.Context, env string) (*PointsController, error) {
	storageCfg := storage.Config{Env: env}

	db, err := storage.NewDynamoDbStorage(storageCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize db: %w", err)
	}

	authContext, err := handlers.GetAuthContext()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize auth context: %w", err)
	}

	return &PointsController{
		BaseController: handlers.BaseController{
			AuthContext: authContext,
		},
		pointsDB: db,
		userDB:   db,
	}, nil
}
