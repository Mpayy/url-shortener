package http

import (
	"errors"
	"net/http"
	"url-shortener/internal/exception"
	"url-shortener/internal/model"
	"url-shortener/internal/usecase"
	"url-shortener/internal/util"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type AuthControllerImpl struct {
	AuthUsecase usecase.AuthUsecase
	Validator   *validator.Validate
	Log         *logrus.Logger
}

func NewAuthController(authUseCase usecase.AuthUsecase, validator *validator.Validate, log *logrus.Logger) AuthController {
	return &AuthControllerImpl{AuthUsecase: authUseCase, Validator: validator, Log: log}
}

func (c *AuthControllerImpl) Register(ctx *gin.Context) {
	var request model.RegisterUserRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		util.ResponseError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := c.Validator.Struct(&request); err != nil {
		util.ResponseError(ctx, http.StatusBadRequest, exception.ExtractValidationErrors(err))
		return
	}

	user, err := c.AuthUsecase.Register(ctx.Request.Context(), &request)
	if err != nil {
		if errors.Is(err, exception.ErrDuplicatedKeyEmail) || errors.Is(err, exception.ErrDuplicatedKeyUsername) {
			util.ResponseError(ctx, http.StatusConflict, err.Error())
			return
		}

		c.Log.WithError(err).Error("Unexpected error during registration")
		util.ResponseError(ctx, http.StatusInternalServerError, "internal server error")
		return
	}

	util.ResponseSuccess(ctx, http.StatusCreated, user)
}

func (c *AuthControllerImpl) Login(ctx *gin.Context) {
	var request model.LoginUserRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		util.ResponseError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := c.Validator.Struct(&request); err != nil {
		util.ResponseError(ctx, http.StatusBadRequest, exception.ExtractValidationErrors(err))
		return
	}

	user, err := c.AuthUsecase.Login(ctx.Request.Context(), &request)
	if err != nil {
		if errors.Is(err, exception.ErrUnauthorized) {
			util.ResponseError(ctx, http.StatusUnauthorized, err.Error())
			return
		}

		c.Log.WithError(err).Error("Unexpected error during login")
		util.ResponseError(ctx, http.StatusInternalServerError, "internal server error")
		return
	}

	util.ResponseSuccess(ctx, http.StatusOK, user)
}
