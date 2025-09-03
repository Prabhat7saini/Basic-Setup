package initializer

import (
	"github.com/Prabhat7saini/Basic-Setup/config"
	"go.uber.org/zap"
)

type BaseHandler struct {
}

func NewBaseHandler(logger *zap.Logger, cfg *config.Env, baseService *BaseService) *BaseHandler {

	return &BaseHandler{}
}
