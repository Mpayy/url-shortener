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

func ToUserTokenResponse(user *entity.User) *model.UserResponse {
	return &model.UserResponse{
		Token: user.Token,
	}
}
