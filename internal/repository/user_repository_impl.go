package repository

import (
	"context"
	"errors"
	"strings"
	"url-shortener/internal/entity"
	"url-shortener/internal/exception"
	"url-shortener/internal/util"

	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{DB: db}
}

func (r *UserRepositoryImpl) GetTx(ctx context.Context) *gorm.DB {
	if tx, ok := util.GetTxFromContext(ctx); ok {
		return tx
	}
	return r.DB
}

func (r *UserRepositoryImpl) Create(ctx context.Context, user *entity.User) error {
	err := r.GetTx(ctx).Create(user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			errMsg := err.Error()

			// Cek nama constraint / kolom di dalam pesan error MySQL
			if strings.Contains(errMsg, "username") {
				return exception.ErrDuplicatedKeyUsername
			}
			if strings.Contains(errMsg, "email") {
				return exception.ErrDuplicatedKeyEmail
			}
		}
		return exception.ErrInternalServer
	}
	return nil
}
