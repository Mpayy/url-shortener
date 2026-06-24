package util

import (
	"url-shortener/internal/model"

	"github.com/gin-gonic/gin"
)

func ResponseError(ctx *gin.Context, code int, message any) {
	ctx.AbortWithStatusJSON(code, model.ErrorResponse{
		Error: message,
	})
}

func ResponseSuccess[T any](ctx *gin.Context, code int, data T) {
	ctx.JSON(code, model.WebResponse[T]{
		Data: data,
	})
}
