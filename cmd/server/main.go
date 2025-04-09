package main

import (
	"context"
	"crud-echo/internal/inbound/customvalidator"
	"crud-echo/internal/inbound/handlers"
	"crud-echo/internal/inbound/routers"
	"crud-echo/internal/outbound/database"
	"crud-echo/internal/pkg/postgres"
	"crud-echo/internal/usecase"
	"log"
	"log/slog"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	db := postgres.NewPostgresDB()
	if err := db.Migrate(); err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	// e.Use(middleware.Logger())

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogError:    true,
		HandleError: true, // forwards error to the global error handler, so it can decide appropriate status code
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				logger.LogAttrs(context.Background(), slog.LevelInfo, "REQUEST",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
				)
			} else {
				logger.LogAttrs(context.Background(), slog.LevelError, "REQUEST_ERROR",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
					slog.String("err", v.Error.Error()),
				)
			}
			return nil
		},
	}))

	e.Use(middleware.Recover())
	e.HTTPErrorHandler = handlers.CustomHTTPErrorHandler

	e.Debug = true

	booksRepo := database.NewBooksRepository(db.DB)
	booksValidator := customvalidator.NewCustomValidator(validator.New())
	booksUseCase := usecase.NewBooksUseCase(booksRepo)
	booksHandler := handlers.NewBooksHandler(booksUseCase, booksValidator)
	routers.RegisterRoutes(e, booksHandler)

	e.Logger.Fatal(e.Start(":1323"))
}
