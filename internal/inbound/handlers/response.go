package handlers

import (
	"github.com/labstack/echo/v4"
)

type Response struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"` // maybe for other things
}

func CustomResponse(c echo.Context, code int, status bool, message string, data any) error {
	resp := Response{
		Status:  status,
		Message: message,
		Data:    data,
	}

	return c.JSON(code, resp)
}
