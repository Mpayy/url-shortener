package repository

import (
	"context"
	"url-shortener/internal/entity"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
}
