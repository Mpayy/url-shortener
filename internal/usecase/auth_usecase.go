package usecase

import (
	"context"
	"url-shortener/internal/model"
)

type AuthUsecase interface {
	Register(ctx context.Context, request *model.RegisterUserRequest) (*model.UserResponse, error)
}