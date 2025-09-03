package initializer

import (
	"github.com/Prabhat7saini/Basic-Setup/config"
	"github.com/Prabhat7saini/Basic-Setup/shared/clients/redis"
	"go.uber.org/zap"
)

type BaseRepository struct {
}

func NewBaseRepository(redis redis.Client, logger *zap.Logger, cfg *config.Env, baseRepo *BaseRepository) *BaseRepository {

	return &BaseRepository{}
}
