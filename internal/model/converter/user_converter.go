package converter

import (
	"url-shortener/internal/entity"
	"url-shortener/internal/model"
)

func ToUserResponse(user *entity.User) *model.UserResponse {
	return &model.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}
}

func ToTokenResponse(token string) *model.TokenResponse {
	return &model.TokenResponse{
		Token: token,
	}
}
