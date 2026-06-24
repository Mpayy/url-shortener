package usecase

import (
	"context"
	"url-shortener/internal/model"
)

type AuthUsecase interface {
	Register(ctx context.Context, request *model.RegisterUserRequest) (*model.UserResponse, error)
	Login(ctx context.Context, request *model.LoginUserRequest) (*model.TokenResponse, error)
	Logout(ctx context.Context, token string) (bool, error)
}
