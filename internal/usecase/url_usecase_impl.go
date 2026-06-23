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
)

type UrlUsecaseImpl struct {
	UrlRepository repository.UrlRepository
	Transaction   util.Transaction
	Log           *logrus.Logger
	RedisClient   config.RedisClient
}

func NewUrlUsecase(urlRepository repository.UrlRepository, transaction util.Transaction, log *logrus.Logger, redisClient config.RedisClient) UrlUsecase {
	return &UrlUsecaseImpl{
		UrlRepository: urlRepository,
		Transaction:   transaction,
		Log:           log,
		RedisClient:   redisClient,
	}
}

func (u *UrlUsecaseImpl) CreateUrl(ctx context.Context, request *model.UrlCreateRequest, userId int64) (*model.UrlResponse, error) {
	var url *entity.Url
	err := u.Transaction.WithTransaction(ctx, func(ctx context.Context) error {
		for attempt := range 3 {
			shortCode, err := util.GenerateShortCode()
			if err != nil {
				return err
			}

			urlEntity := &entity.Url{
				UserID:      userId,
				OriginalUrl: request.OriginalUrl,
				ShortCode:   shortCode,
			}

			err = u.UrlRepository.Create(ctx, urlEntity)
			if err == nil {
				url = urlEntity
				return nil
			}

			if !errors.Is(err, exception.ErrDuplicatedKeyShortCode) {
				return err
			}

			u.Log.Warn("ShortCode collision, retrying... attempt: ", attempt+1)
		}

		return exception.ErrInternalServer
	})

	if err != nil {
		return nil, err
	}

	return converter.ToUrlResponse(url), nil
}
