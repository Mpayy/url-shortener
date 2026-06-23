package http

import (
	"net/http"
	"url-shortener/internal/delivery/http/middleware"
	"url-shortener/internal/exception"
	"url-shortener/internal/model"
	"url-shortener/internal/usecase"
	"url-shortener/internal/util"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type UrlControllerImpl struct {
	UrlUsecase usecase.UrlUsecase
	Validator  *validator.Validate
	Log        *logrus.Logger
}

func NewUrlController(urlUsecase usecase.UrlUsecase, validator *validator.Validate, log *logrus.Logger) UrlController {
	return &UrlControllerImpl{UrlUsecase: urlUsecase, Validator: validator, Log: log}
}

func (c *UrlControllerImpl) CreateUrl(ctx *gin.Context) {
	auth, exists := middleware.GetAuthFromCtx(ctx)
	if !exists {
		util.ResponseError(ctx, http.StatusUnauthorized, "unauthorized")
		return
	}

	var request model.UrlCreateRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		util.ResponseError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := c.Validator.Struct(&request); err != nil {
		util.ResponseError(ctx, http.StatusBadRequest, exception.ExtractValidationErrors(err))
		return
	}

	url, err := c.UrlUsecase.CreateUrl(ctx.Request.Context(), &request, auth.ID)
	if err != nil {
		util.ResponseError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	util.ResponseSuccess(ctx, http.StatusCreated, url)
}
