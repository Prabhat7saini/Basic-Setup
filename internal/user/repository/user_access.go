package user_repository

import (
	"github.com/Prabhat7saini/Basic-Setup/config"
	"github.com/Prabhat7saini/Basic-Setup/shared/clients/redis"
	"gorm.io/gorm"
)

type UserRepositoryAccess struct {
	DB     *gorm.DB
	Redis  redis.Client
	Config *config.Env
}

func NewUserRepositoryAccess(db *gorm.DB, redis redis.Client, config *config.Env) *UserRepositoryAccess {
	return &UserRepositoryAccess{
		DB:     db,
	}
}
