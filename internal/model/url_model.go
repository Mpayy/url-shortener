package model

type UrlResponse struct {
	ID          int64  `json:"id,omitempty"`
	ShortCode   string `json:"short_code,omitempty"`
	OriginalUrl string `json:"original_url,omitempty"`
	Hits        int64  `json:"hits,omitempty"`
}

type UrlCreateRequest struct {
	OriginalUrl string `json:"original_url" validate:"required,http_url"`
}
