package dto

import (
	"errors"

	"github.com/Prabhat7saini/Basic-Setup/shared/constants"
	"github.com/Prabhat7saini/Basic-Setup/shared/constants/exception"
	"github.com/go-playground/validator/v10"
)

type CreateUserDto struct {
	Name     string `json:"name" binding:"required,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

var CreateUserValidationErrorMap = map[string]map[string]*constants.Exception{
	"Name": {
		"required": exception.GetException(exception.FIRST_NAME_REQUIRED),
	},
	"Email": {
		"required": exception.GetException(exception.EMAIL_REQUIRED),
	},
	"PASSWORD": {
		"required": exception.GetException(exception.PASSWORD_REQUIRED),
	},
}

func GetCreateUserDtoValidationError(err error) *constants.Exception {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		for _, fe := range ve {
			if tagMap, ok := CreateUserValidationErrorMap[fe.Field()]; ok {
				if exc, ok := tagMap[fe.Tag()]; ok {
					return exc
				}
			}
		}
	}
	return exception.GetException(exception.INVALID_PAYLOAD)
}
