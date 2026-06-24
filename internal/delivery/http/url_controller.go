package http

import "github.com/gin-gonic/gin"

type UrlController interface {
	CreateUrl(ctx *gin.Context)
	GetAll(ctx *gin.Context)
}