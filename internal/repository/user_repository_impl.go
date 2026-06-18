package repository

import (
	"context"
	"url-shortener/internal/entity"

	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	DB *gorm.DB
}

func NewUserRepositoryImpl(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{DB: db}
}

func (r *UserRepositoryImpl) Registration(ctx context.Context, user *entity.User) error {
	err := r.DB.WithContext(ctx).Create(user).Error
	if err != nil {
		return err
	}
	return nil
}