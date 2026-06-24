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

func ToUrlResponses(urls []entity.Url) []model.UrlResponse {
	responses := make([]model.UrlResponse, 0, len(urls))
	for _, url := range urls {
		responses = append(responses, *ToUrlResponse(&url))
	}
	return responses
}
