package repository

import (
	"context"
	"errors"
	"url-shortener/internal/entity"
	"url-shortener/internal/exception"
	"url-shortener/internal/util"

	"gorm.io/gorm"
)

type UrlRepositoryImpl struct {
	DB *gorm.DB
}

func NewUrlRepository(db *gorm.DB) UrlRepository {
	return &UrlRepositoryImpl{DB: db}
}

func (r *UrlRepositoryImpl) GetTx(ctx context.Context) *gorm.DB {
	if tx, ok := util.GetTxFromContext(ctx); ok {
		return tx
	}

	return r.DB
}

func (r *UrlRepositoryImpl) Create(ctx context.Context, url *entity.Url) error {
	err := r.GetTx(ctx).Create(url).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return exception.ErrDuplicatedKeyShortCode
		}
		return err
	}

	return nil
}

func (r *UrlRepositoryImpl) FindByUserID(ctx context.Context, userID int64) ([]entity.Url, error) {
	var urls []entity.Url
	err := r.GetTx(ctx).Where("user_id = ?", userID).Find(&urls).Error
	if err != nil {
		return nil, err
	}

	return urls, nil
}
