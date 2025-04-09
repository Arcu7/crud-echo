package handlers

import (
	"crud-echo/internal/models"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

func CustomHTTPErrorHandler(err error, c echo.Context) {
	var code int
	var message string
	// maybe this can be used if needed?
	// var data any

	var httpError *echo.HTTPError
	if errors.As(err, &httpError) {
		code = httpError.Code
		message = httpError.Message.(string)
	} else {
		code = http.StatusInternalServerError
		message = models.InternalServerError
	}

	resp := Response{
		Status:  false,
		Message: message,
		Data:    nil,
	}

	c.JSON(code, resp)
}
