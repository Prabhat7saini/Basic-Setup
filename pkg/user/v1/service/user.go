package user_service

import (
	"context"
	"net/http"

	"github.com/Prabhat7saini/Basic-Setup/pkg/user/models"
	"github.com/Prabhat7saini/Basic-Setup/pkg/user/repository"
	"github.com/Prabhat7saini/Basic-Setup/pkg/user/v1/dto"
	"github.com/Prabhat7saini/Basic-Setup/shared/constants"
	"github.com/Prabhat7saini/Basic-Setup/shared/constants/exception"
	"github.com/Prabhat7saini/Basic-Setup/shared/logger"
	"github.com/Prabhat7saini/Basic-Setup/shared/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserServiceMethods interface {
	CreateUser(ctx context.Context, db *gorm.DB, dto dto.CreateUserDTO) constants.ServiceOutput[*models.User]
}

type userService struct {
	repo user_repository.UserRepositoryMethods
}

func NewUserService() UserServiceMethods {
	return &userService{
		repo: user_repository.NewUserRepository(), // service handles repo creation
	}
}

func (s *userService) CreateUser(ctx context.Context, db *gorm.DB, dto dto.CreateUserDTO) constants.ServiceOutput[*models.User] {
	logger.Info("Creating User", zap.String("email", dto.Email))

	existingUser, exc := s.repo.FindUserByFields(ctx, db, map[string]any{"email": dto.Email}, "id")
	if exc != nil {
		return utils.HandleException[*models.User](*exc)
	}
	if existingUser != nil {
		return utils.ServiceError[*models.User](exception.USER_ALREADY_EXISTS)
	}

	hashPassword, err := utils.HashPassword(dto.Password)
	if err != nil {
		logger.Error("while hashing the password", zap.Error(err))
		return utils.ServiceError[*models.User](exception.INTERNAL_SERVER_ERROR)
	}


	user, incomingException:= s.repo.CreateUser(ctx,db, map[string]any{
		"email":         dto.Email,
		"name":          dto.Name,
		"password_hash": hashPassword,
		"created_by":    dto.CreatedBy,
	})

	if incomingException != nil {
		return utils.HandleException[*models.User](*incomingException)
	}

	return constants.ServiceOutput[*models.User]{
		Message:        "User created successfully",
		OutputData:     user,
		HttpStatusCode: http.StatusCreated,
		RespStatusCode: http.StatusCreated,
	}
}
