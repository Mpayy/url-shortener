package repository

import (
	"context"
	"url-shortener/internal/entity"
)

type UrlRepository interface {
	Create(ctx context.Context, url *entity.Url) error
}
