package repository

import (
	"context"
	"url-shortener/internal/entity"
)

type UrlRepository interface {
	Create(ctx context.Context, url *entity.Url) error
	FindByUserID(ctx context.Context, userID int64) ([]entity.Url, error)
}
