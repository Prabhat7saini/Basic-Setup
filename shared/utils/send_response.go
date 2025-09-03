package utils

import (
	"github.com/Prabhat7saini/Basic-Setup/shared/constants"
	"github.com/gin-gonic/gin"
	// "gitlab.com/truemeds-dev-team/truemeds-dev-doctor/truemeds-dev-service/doctorportal-auth-service/shared/constants"
)

func SendRestResponse[T any](ctx *gin.Context, output constants.ServiceOutput[T]) {
	if output.Exception != nil {
		ctx.JSON(output.Exception.HttpStatusCode, constants.ApiResponse[any]{
			Code:    output.Exception.Code,
			Message: fallbackIfEmpty(output.Message, output.Exception.Message),
		})
		return
	}

	ctx.JSON(output.HttpStatusCode, constants.ApiResponse[T]{
		Code:    output.RespStatusCode,
		Message: fallbackIfEmpty(output.Message, "SUCCESS"),
		Data:    output.OutputData,
	})
}

func fallbackIfEmpty(preferred string, fallback string) string {
	if preferred != "" {
		return preferred
	}
	return fallback
}
