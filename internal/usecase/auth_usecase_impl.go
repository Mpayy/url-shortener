package usecase

import (
	"context"
	"url-shortener/internal/entity"
	"url-shortener/internal/model"
	"url-shortener/internal/model/converter"
	"url-shortener/internal/repository"
	"url-shortener/internal/util"

	"golang.org/x/crypto/bcrypt"
)

type AuthUsecaseImpl struct {
	UserRepository repository.UserRepository
	Transaction    util.Transaction
}

func NewAuthUsecase(userRepository repository.UserRepository, transaction util.Transaction) AuthUsecase {
	return &AuthUsecaseImpl{
		UserRepository: userRepository,
		Transaction:    transaction,
	}
}

func (u *AuthUsecaseImpl) Register(ctx context.Context, request *model.RegisterUserRequest) (*model.UserResponse, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &entity.User{
		Username: request.Username,
		Email:    request.Email,
		Password: string(hashPassword),
	}

	err = u.Transaction.WithTransaction(ctx, func(ctx context.Context) error {
		errCreate := u.UserRepository.Create(ctx, user)
		if errCreate != nil {
			return errCreate
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return converter.ToUserResponse(user), nil
}
