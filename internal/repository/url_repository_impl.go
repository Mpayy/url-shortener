package repository

import (
	"context"
	"errors"
	"fmt"
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
		return tx.WithContext(ctx)
	}

	return r.DB.WithContext(ctx)
}

func (r *UrlRepositoryImpl) Create(ctx context.Context, url *entity.Url) error {
	err := r.GetTx(ctx).Create(url).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return exception.ErrDuplicatedKeyShortCode
		}
		return fmt.Errorf("Create: %w", exception.ErrInternalServer)
	}

	return nil
}

func (r *UrlRepositoryImpl) FindByUserID(ctx context.Context, userID int64) ([]entity.Url, error) {
	var urls []entity.Url
	err := r.GetTx(ctx).Where("user_id = ?", userID).Find(&urls).Error
	if err != nil {
		return nil, fmt.Errorf("FindByUserID: %w", exception.ErrInternalServer)
	}

	return urls, nil
}

func (r *UrlRepositoryImpl) Delete(ctx context.Context, shortCode string, userID int64) error {
	result := r.GetTx(ctx).Where("short_code = ? AND user_id = ?", shortCode, userID).Delete(&entity.Url{})
	if result.Error != nil {
		return fmt.Errorf("Delete: %w", exception.ErrInternalServer)
	}
	if result.RowsAffected == 0 {
		return exception.ErrNotFound
	}

	return nil
}

func (r *UrlRepositoryImpl) FindByShortCode(ctx context.Context, shortCode string) (*entity.Url, error) {
	var url entity.Url
	err := r.GetTx(ctx).Where("short_code = ?", shortCode).First(&url).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.ErrNotFound
		}

		return nil, fmt.Errorf("FindByShortCode: %w", exception.ErrInternalServer)
	}

	return &url, nil
}

func (r *UrlRepositoryImpl) IncrementHits(ctx context.Context, shortCode string) error {
	err := r.GetTx(ctx).Model(&entity.Url{}).Where("short_code = ?", shortCode).UpdateColumn("hits", gorm.Expr("hits + ?", 1)).Error
	if err != nil {
		return fmt.Errorf("IncrementHits: %w", exception.ErrInternalServer)
	}

	return nil
}