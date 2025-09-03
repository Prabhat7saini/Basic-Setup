package user_repository

import (
	"context"

	"github.com/Prabhat7saini/Basic-Setup/internal/user/models"
	"github.com/Prabhat7saini/Basic-Setup/shared/constants"
	"github.com/Prabhat7saini/Basic-Setup/shared/constants/exception"
	"github.com/Prabhat7saini/Basic-Setup/shared/logger"

	// "github.com/Prabhat7saini/Basic-Setup/shared/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type userRepository struct {
	access UserRepositoryAccess
}

type UserRepositoryMethods interface {
	CreateUser(ctx context.Context, data map[string]any) (*models.User, *constants.Exception)
	FindUserByFields(ctx context.Context, conditions map[string]any, selectFields ...string) (*models.User, *constants.Exception)
}

func NewUserRepository(access UserRepositoryAccess) UserRepositoryMethods {
	return &userRepository{
		access: access,
	}
}

func (ul *userRepository) FindUserByFields(ctx context.Context, conditions map[string]any, selectFields ...string) (*models.User, *constants.Exception) {
	var user models.User
	db := ul.access.DB

	query := db.WithContext(ctx).Model(&models.User{})

	// Optional select
	if len(selectFields) > 0 {
		query = query.Select(selectFields)
	}

	// Apply map conditions safely
	query = query.Where(conditions)

	// Execute query
	err := query.First(&user).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		logger.Error("while find user", zap.Error(err))
		return nil, exception.GetException(exception.INTERNAL_SERVER_ERROR)
	}
	return &user, nil
}

func (ul *userRepository) CreateUser(ctx context.Context, data map[string]any) (*models.User, *constants.Exception) {
	user := &models.User{
		Email:    data["email"].(string),
		Password: data["password"].(string),
	}

	if createdBy, ok := data["created_by"].(int); ok {
		user.CreatedBy = &createdBy
	}
	if name, ok := data["name"].(string); ok {
		user.Name = &name
	}

	if err := ul.access.DB.WithContext(ctx).Create(user).Error; err != nil {
		logger.Error("while insert new user in database", zap.Error(err))
		return nil, exception.GetException(exception.INTERNAL_SERVER_ERROR)
	}
	return user, nil
}
