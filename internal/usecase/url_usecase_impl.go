package usecase

import (
	"context"
	"errors"
	"time"
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
	u.Log.WithField("user_id", userId).Debug("Creating URL")

	var url *entity.Url
	var lastErr error

	for attempt := range 3 {
		shortCode, err := util.GenerateShortCode()
		if err != nil {
			u.Log.WithError(err).Error("Failed to generate short code")
			return nil, exception.ErrInternalServer
		}

		urlEntity := &entity.Url{
			UserID:      userId,
			OriginalUrl: request.OriginalUrl,
			ShortCode:   shortCode,
		}

		err = u.UrlRepository.Create(ctx, urlEntity)
		if err == nil {
			url = urlEntity
			return converter.ToUrlResponse(url), nil
		}
		lastErr = err

		if !errors.Is(err, exception.ErrDuplicatedKeyShortCode) {
			break
		}

		u.Log.Warn("ShortCode collision, retrying... attempt: ", attempt+1)
	}

	u.Log.WithField("user_id", userId).WithError(lastErr).Error("Failed to create URL")
	return nil, exception.ErrInternalServer

}

func (u *UrlUsecaseImpl) GetUserUrls(ctx context.Context, userId int64) ([]model.UrlResponse, error) {
	u.Log.WithField("user_id", userId).Debug("Getting user URLs")

	urls, err := u.UrlRepository.FindByUserID(ctx, userId)
	if err != nil {
		u.Log.WithError(err).Error("Failed to get user URLs")
		return nil, exception.ErrInternalServer
	}

	u.Log.WithField("user_id", userId).Info("User URLs retrieved successfully")
	return converter.ToUrlResponses(urls), nil
}

func (u *UrlUsecaseImpl) DeleteUrl(ctx context.Context, shortCode string, userId int64) (bool, error) {
	u.Log.WithFields(logrus.Fields{"short_code": shortCode, "user_id": userId}).Debug("Deleting URL")

	err := u.Transaction.WithTransaction(ctx, func(tx context.Context) error {
		return u.UrlRepository.Delete(tx, shortCode, userId)
	})

	if err != nil {
		if errors.Is(err, exception.ErrNotFound) {
			u.Log.WithFields(logrus.Fields{"short_code": shortCode, "user_id": userId}).WithError(err).Warn("URL not found")
			return false, err
		}
		u.Log.WithFields(logrus.Fields{"short_code": shortCode, "user_id": userId}).WithError(err).Error("Failed to delete URL")
		return false, exception.ErrInternalServer
	}

	err = u.RedisClient.Delete(ctx, config.UrlCachePrefix+shortCode)
	if err != nil {
		u.Log.WithFields(logrus.Fields{"short_code": shortCode, "user_id": userId}).Warn("Failed to delete URL from Redis")
	}

	u.Log.WithFields(logrus.Fields{"short_code": shortCode, "user_id": userId}).Info("URL deleted successfully")
	return true, nil
}

func (u *UrlUsecaseImpl) asyncIncrementHits(shortCode string) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := u.UrlRepository.IncrementHits(ctx, shortCode)
		if err != nil {
			u.Log.WithField("short_code", shortCode).WithError(err).Warn("Failed to increment hits")
		}
	}()
}

func (u *UrlUsecaseImpl) Redirect(ctx context.Context, shortCode string) (string, error) {
	u.Log.WithField("short_code", shortCode).Debug("Redirecting URL")

	originalUrl, err := u.RedisClient.Get(ctx, config.UrlCachePrefix+shortCode)
	if err == nil {
		u.Log.WithField("short_code", shortCode).Info("URL found in cache")
		u.asyncIncrementHits(shortCode)
		return originalUrl, nil
	}

	u.Log.WithField("short_code", shortCode).Debug("Cache miss")
	url, err := u.UrlRepository.FindByShortCode(ctx, shortCode)
	if err != nil {
		if errors.Is(err, exception.ErrNotFound) {
			u.Log.WithField("short_code", shortCode).WithError(err).Warn("URL not found")
			return "", err
		}
		u.Log.WithField("short_code", shortCode).WithError(err).Error("Failed to find URL")
		return "", exception.ErrInternalServer
	}

	err = u.RedisClient.Set(ctx, config.UrlCachePrefix+shortCode, url.OriginalUrl, 24*time.Hour)
	if err != nil {
		u.Log.WithField("short_code", shortCode).WithError(err).Warn("Failed to set URL in cache")
	}

	u.asyncIncrementHits(shortCode)

	u.Log.WithField("short_code", shortCode).Info("URL redirected successfully")
	return url.OriginalUrl, nil
}
