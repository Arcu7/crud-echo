package handlers

import "github.com/labstack/echo/v4"

type Response struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

func CustomResponse(c echo.Context, code int, status bool, message string) error {
	resp := Response{
		Status:  status,
		Message: message,
	}

	return c.JSON(code, resp)
}
