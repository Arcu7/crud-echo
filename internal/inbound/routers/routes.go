package routers

import (
	"crud-echo/internal/inbound/handlers"
	"crud-echo/internal/inbound/server"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Router struct {
	srv *server.Server
	h   *handlers.BooksHandler
}

func NewRouter(srv *server.Server, h *handlers.BooksHandler) *Router {
	return &Router{
		srv: srv,
		h:   h,
	}
}

func (r *Router) RegisterRoutes() {
	e := r.srv.GetEcho()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.POST("/book", r.h.CreateBook)
	e.GET("/books", r.h.GetAllBooks)
	e.GET("/book/:id", r.h.GetBookByID)
	e.PUT("/book", r.h.UpdateBook)
	e.DELETE("/book", r.h.DeleteBook)
}
