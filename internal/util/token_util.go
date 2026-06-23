package util

import (
	"time"
	"url-shortener/internal/exception"
	"url-shortener/internal/model"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

const TokenDuration = 24 * time.Hour * 30

type TokenUtil interface {
	CreateToken(auth *model.Auth) (string, error)
	ParseToken(token string) (*model.Auth, error)
}

type TokenUtilImpl struct {
	SecretKey string
}

func NewTokenUtil(config *viper.Viper) TokenUtil {
	secretKey := config.GetString("JWT_SECRET_KEY")
	return &TokenUtilImpl{
		SecretKey: secretKey,
	}
}

func (t *TokenUtilImpl) CreateToken(auth *model.Auth) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       auth.ID,
		"username": auth.Username,
		"email":    auth.Email,
		"exp":      time.Now().Add(TokenDuration).Unix(),
	})

	jwtToken, err := token.SignedString([]byte(t.SecretKey))
	if err != nil {
		return "", exception.ErrInternalServer
	}

	return jwtToken, nil
}

func (t *TokenUtilImpl) ParseToken(jwtToken string) (*model.Auth, error) {
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (any, error) {
		return []byte(t.SecretKey), nil
	})

	if err != nil {
		return nil, exception.ErrUnauthorized
	}

	claim, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return nil, exception.ErrUnauthorized
	}

	id, ok := claim["id"].(float64)
	if !ok {
		return nil, exception.ErrUnauthorized
	}

	username, ok := claim["username"].(string)
	if !ok {
		return nil, exception.ErrUnauthorized
	}

	email, ok := claim["email"].(string)
	if !ok {
		return nil, exception.ErrUnauthorized
	}

	auth := &model.Auth{
		ID:       int64(id),
		Username: username,
		Email:    email,
	}

	return auth, nil

}
