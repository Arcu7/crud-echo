package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Response struct {
	Status  bool              `json:"status"`
	Message string            `json:"message,omitempty"`
	Errors  map[string]string `json:"errors,omitempty"`
}

func CustomResponse(c echo.Context, code int, status bool, message string) error {
	resp := Response{
		Status:  status,
		Message: message,
	}

	return c.JSON(code, resp)
}

func ValidationErrorResponse(c echo.Context, message string, errors map[string]string) error {
	resp := Response{
		Status:  false,
		Message: message,
		Errors:  errors,
	}

	return c.JSON(http.StatusBadRequest, resp)
}
