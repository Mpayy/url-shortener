package model

type WebResponse[T any] struct {
	Data  T   `json:"data,omitempty"`
	Error any `json:"error,omitempty"`
}