package util

import (
	"github.com/gin-gonic/gin"
	"url-shortener/internal/model"
)

func ResponseError(ctx *gin.Context, code int, message any) {
	ctx.AbortWithStatusJSON(code, model.WebResponse[any]{
		Error: message,
	})
}

func ResponseSuccess[T any](ctx *gin.Context, code int, data T) {
	ctx.JSON(code, model.WebResponse[T]{
		Data: data,
	})
}
