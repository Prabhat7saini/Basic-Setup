package initializer

import (
	"github.com/Prabhat7saini/Basic-Setup/config"
	"github.com/Prabhat7saini/Basic-Setup/shared/clients/redis"
	"go.uber.org/zap"
)

type BaseService struct {
}

func NewBaseService(redis redis.Client, logger *zap.Logger, cfg *config.Env, baseRepo *BaseRepository) *BaseService {

	return &BaseService{}
}
