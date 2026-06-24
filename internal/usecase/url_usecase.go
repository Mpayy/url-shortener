package usecase

import (
	"context"
	"url-shortener/internal/model"
)

type UrlUsecase interface {
	CreateUrl(ctx context.Context, request *model.UrlCreateRequest, userId int64) (*model.UrlResponse, error)
	GetUserUrls(ctx context.Context, userId int64) ([]model.UrlResponse, error)
	DeleteUrl(ctx context.Context, shortCode string, userId int64) (bool, error)
	Redirect(ctx context.Context, shortCode string) (string, error)
}
