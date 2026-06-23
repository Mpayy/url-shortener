package usecase

import (
	"context"
	"url-shortener/internal/model"
)

type UrlUsecase interface {
	CreateUrl(ctx context.Context, request *model.UrlCreateRequest, userId int64) (*model.UrlResponse, error)	
}