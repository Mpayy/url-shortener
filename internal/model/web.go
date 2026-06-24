package model

type WebResponse[T any] struct {
	Data T `json:"data,omitempty"`
}

type ErrorResponse struct {
	Error any `json:"error"`
}