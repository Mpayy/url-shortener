package repository

import (
	"context"
	"errors"
	"strings"
	"url-shortener/internal/entity"
	"url-shortener/internal/exception"
	"url-shortener/internal/util"

	"github.com/go-sql-driver/mysql"
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
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			msg := strings.ToLower(mysqlErr.Message)
			if strings.Contains(msg, "username") {
				return exception.ErrDuplicatedKeyUsername
			}
			if strings.Contains(msg, "email") {
				return exception.ErrDuplicatedKeyEmail
			}
		}
		return exception.ErrInternalServer
	}
	return nil
}

func (r *UserRepositoryImpl) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	err := r.DB.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.ErrNotFound
		}
		return nil, exception.ErrInternalServer
	}
	return &user, nil
}
