package server

import (
	"crud-echo/internal/config"
	"crud-echo/internal/inbound/handlers"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	e   *echo.Echo
	cfg *config.Config
}

func NewServer(cfg *config.Config) *Server {
	e := echo.New()

	return &Server{
		e:   e,
		cfg: cfg,
	}
}

func (s *Server) Start() error {
	s.e.Use(middleware.Logger())
	s.e.Use(middleware.Recover())

	s.e.HTTPErrorHandler = handlers.CustomHTTPErrorHandler

	s.e.Debug = true

	// Start server
	addr := fmt.Sprintf("%s:%d", s.cfg.Server.Host, s.cfg.Server.Port)
	return s.e.Start(addr)
}

func (s *Server) GetEcho() *echo.Echo {
	return s.e
}
