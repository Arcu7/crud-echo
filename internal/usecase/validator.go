package usecase

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type CustomValidator struct {
	Validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.Validator.Struct(i); err != nil {
		fmterr := formatValidationErrors(err)
		return echo.NewHTTPError(http.StatusBadRequest, fmterr)
	}
	return nil
}

func formatValidationErrors(err error) map[string]string {
	errors := make(map[string]string)

	for _, err := range err.(validator.ValidationErrors) {
		field := err.Field()
		errors[field] = formatErrorMsg(err)
	}

	return errors
}

func formatErrorMsg(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "This field is required"
	case "min":
		return "Should be at least " + err.Param() + " characters long"
	case "max":
		return "Should be at most " + err.Param() + " characters long"
	case "lte":
		return "Should be less than " + err.Param()
	case "gte":
		return "Should be greater than " + err.Param()
	default:
		return "Invalid value"
	}
}
