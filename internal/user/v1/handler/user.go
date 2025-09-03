package handler

import (
	"github.com/Prabhat7saini/Basic-Setup/internal/user/v1/dto"
	user_service "github.com/Prabhat7saini/Basic-Setup/internal/user/v1/service"
	"github.com/Prabhat7saini/Basic-Setup/shared/constants"
	"github.com/Prabhat7saini/Basic-Setup/shared/utils"
	"github.com/gin-gonic/gin"
)

type UserHandlerMethods interface {
	CreateUser(ctx *gin.Context)
}
type userHandler struct {
	userService user_service.UserServiceMethods
}

func NewAuthHandler(userService user_service.UserServiceMethods) UserHandlerMethods {
	return &userHandler{userService: userService}
}

func (ah *userHandler) CreateUser(ctx *gin.Context) {
	var req dto.CreateUserDto
	if err := ctx.ShouldBindJSON(&req); err != nil {
		exc := dto.GetCreateUserDtoValidationError(err)
		resp := constants.ServiceOutput[struct{}]{
			Exception: exc,
		}
		utils.SendRestResponse(ctx, resp)
		return
	}
	output := ah.userService.CreateUser(ctx, req.Name, req.Email, req.Password)
	utils.SendRestResponse(ctx, output)
}
