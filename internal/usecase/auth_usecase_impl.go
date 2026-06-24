package usecase

import (
	"context"
	"errors"
	"url-shortener/internal/config"
	"url-shortener/internal/entity"
	"url-shortener/internal/exception"
	"url-shortener/internal/model"
	"url-shortener/internal/model/converter"
	"url-shortener/internal/repository"
	"url-shortener/internal/util"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecaseImpl struct {
	UserRepository repository.UserRepository
	Transaction    util.Transaction
	TokenUtil      util.TokenUtil
	RedisClient    config.RedisClient
	Log            *logrus.Logger
}

func NewAuthUsecase(userRepository repository.UserRepository, transaction util.Transaction, tokenUtil util.TokenUtil, redisClient config.RedisClient, log *logrus.Logger) AuthUsecase {
	return &AuthUsecaseImpl{
		UserRepository: userRepository,
		Transaction:    transaction,
		TokenUtil:      tokenUtil,
		RedisClient:    redisClient,
		Log:            log,
	}
}

func (u *AuthUsecaseImpl) Register(ctx context.Context, request *model.RegisterUserRequest) (*model.UserResponse, error) {
	u.Log.WithField("email", request.Email).Debug("Attempting to register user")

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		u.Log.WithError(err).Error("Failed to hash password")
		return nil, exception.ErrInternalServer
	}

	user := &entity.User{
		Username: request.Username,
		Email:    request.Email,
		Password: string(hashPassword),
	}

	err = u.Transaction.WithTransaction(ctx, func(ctx context.Context) error {
		return u.UserRepository.Create(ctx, user)
	})

	if err != nil {
		if errors.Is(err, exception.ErrDuplicatedKeyEmail) {
			u.Log.WithField("email", request.Email).Warn("Duplicate email registration attempt")
			return nil, err
		}
		if errors.Is(err, exception.ErrDuplicatedKeyUsername) {
			u.Log.WithField("username", request.Username).Warn("Duplicate username registration attempt")
			return nil, err
		}
		u.Log.WithError(err).Error("Failed to create user")
		return nil, err
	}

	u.Log.WithField("user_id", user.ID).Info("User registered successfully")
	return converter.ToUserResponse(user), nil
}

func (u *AuthUsecaseImpl) Login(ctx context.Context, request *model.LoginUserRequest) (*model.UserResponse, error) {
	u.Log.WithField("email", request.Email).Debug("Login attempt")

	user, err := u.UserRepository.FindByEmail(ctx, request.Email)
	if err != nil {
		if errors.Is(err, exception.ErrNotFound) {
			u.Log.WithField("email", request.Email).Warn("Login failed: user not found")
			return nil, exception.ErrUnauthorized
		}
		u.Log.WithError(err).Error("Failed to find user")
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		u.Log.WithField("email", request.Email).Warn("Login failed: wrong password")
		return nil, exception.ErrUnauthorized
	}

	auth := &model.Auth{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}

	token, err := u.TokenUtil.CreateToken(auth)
	if err != nil {
		u.Log.WithField("user_id", user.ID).Error("Failed to create token")
		return nil, err
	}

	err = u.RedisClient.Set(ctx, config.AuthPrefix+token, user.ID, util.TokenDuration)
	if err != nil {
		u.Log.WithField("user_id", user.ID).Error("Failed to set token")
		return nil, err
	}

	u.Log.WithField("user_id", user.ID).Info("User logged in successfully")
	return converter.ToUserTokenResponse(token), nil
}

func (u *AuthUsecaseImpl) Logout(ctx context.Context, token string) (bool, error) {
	u.Log.WithField("token", token).Debug("Logout attempt")

	err := u.RedisClient.Delete(ctx, config.AuthPrefix+token)
	if err != nil {
		u.Log.WithField("token", token).Error("Failed to delete token")
		return false, err
	}

	u.Log.WithField("token", token).Info("User logged out successfully")
	return true, nil
}
