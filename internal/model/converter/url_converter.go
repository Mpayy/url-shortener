package converter

import (
	"url-shortener/internal/entity"
	"url-shortener/internal/model"
)

func ToUrlResponse(url *entity.Url) *model.UrlResponse {
	return &model.UrlResponse{
		ID:          url.ID,
		ShortCode:   url.ShortCode,
		OriginalUrl: url.OriginalUrl,
	}
}
