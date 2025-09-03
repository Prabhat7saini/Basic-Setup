package user_service

import (
	"context"
	"net/http"

	user_repository "github.com/Prabhat7saini/Basic-Setup/internal/user/repository"
	"github.com/Prabhat7saini/Basic-Setup/shared/constants"
	"github.com/Prabhat7saini/Basic-Setup/shared/constants/exception"
	"github.com/Prabhat7saini/Basic-Setup/shared/logger"
	"github.com/Prabhat7saini/Basic-Setup/shared/utils"
	"go.uber.org/zap"
)

type UserServiceMethods interface {
	CreateUser(ctx context.Context, name string, email string, password string) constants.ServiceOutput[*struct{}]
}

type userService struct {
	repo   user_repository.UserRepositoryMethods
	access UserServiceAccess
}

func NewUserService(repo user_repository.UserRepositoryMethods, access UserServiceAccess) UserServiceMethods {
	return &userService{
		repo:   repo,
		access: access,
	}
}

func (s *userService) CreateUser(ctx context.Context, name string, email string, password string) constants.ServiceOutput[*struct{}] {

	logger.Info("Creating User", zap.String("name", name), zap.String("email", email))
	existingUser, incomingException := s.repo.FindUserByFields(ctx, map[string]any{"email": email}, "id")

	if incomingException != nil {
		return utils.HandleException[*struct{}](*incomingException)

	}

	if existingUser != nil {
		logger.Info("User already exists", zap.String("email", email))
		return utils.ServiceError[*struct{}](exception.USER_ALREADY_EXISTS)
	}

	hashPassword, err := utils.HashPassword(password)

	if err != nil {
		logger.Error("while hashing the password", zap.Error(err))
		return utils.ServiceError[*struct{}](exception.USER_ALREADY_EXISTS)
	}
	_, incomingException = s.repo.CreateUser(ctx, map[string]any{
		"email":         email,
		"name":          name,
		"password_hash": hashPassword,
		"created_by":    -1,
	})

	if incomingException != nil {
		return utils.HandleException[*struct{}](*incomingException)
	}
	return constants.ServiceOutput[*struct{}]{
		Message:        "User created successfully",
		OutputData:     nil,
		HttpStatusCode: http.StatusCreated,
		RespStatusCode: http.StatusCreated,
	}

}
