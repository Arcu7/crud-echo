package routers

import (
	"crud-echo/internal/inbound/handlers"
	"net/http"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, h *handlers.BooksHandler) {
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.POST("/book", h.CreateBook)
	e.GET("/books", h.GetAllBooks)
	e.GET("/book/:id", h.GetBookByID)
	e.PUT("/book", h.UpdateBook)
	e.DELETE("/book", h.DeleteBook)
}
