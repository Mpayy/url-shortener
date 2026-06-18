package repository

import (
	"context"
	"url-shortener/internal/entity"
)

type UserRepository interface {
	Registration(ctx context.Context, user *entity.User) error
}
