package user_service

import (
	"github.com/Prabhat7saini/Basic-Setup/config"
	"github.com/Prabhat7saini/Basic-Setup/shared/clients/redis"
)

type UserServiceAccess struct {
	Redis  redis.Client
	Config *config.Env
}

func NewUserRepositoryAcess(redis redis.Client, config *config.Env) *UserServiceAccess {
	return &UserServiceAccess{
		Redis:  redis,
		Config: config,
	}
}
