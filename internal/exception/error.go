package exception

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

var (
	ErrUnauthorized           = errors.New("unauthorized")
	ErrInternalServer         = errors.New("internal server error")
	ErrNotFound               = errors.New("data not found")
	ErrDuplicatedKeyEmail     = errors.New("email already exists")
	ErrDuplicatedKeyUsername  = errors.New("username already exists")
	ErrDuplicatedKeyShortCode = errors.New("short code already exists")
)

func ExtractValidationErrors(err error) map[string]string {
	errorReport := make(map[string]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		errorReport = TranslateValidationError(validationErrors)
	}

	return errorReport
}

func TranslateValidationError(valErr validator.ValidationErrors) map[string]string {
	fieldError := make(map[string]string)
	for _, e := range valErr {
		var message string
		switch e.Tag() {
		case "required":
			message = "must be filled"
		case "email":
			message = "must be a valid email"
		case "min":
			message = "must be at least " + e.Param() + " characters long"
		case "max":
			message = "must be at most " + e.Param() + " characters long"
		default:
			message = "invalid input value"
		}
		fieldError[strings.ToLower(e.Field())] = message
	}

	return fieldError
}