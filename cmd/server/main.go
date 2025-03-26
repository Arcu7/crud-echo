package main

import (
	"crud-echo/internal/inbound/handlers"
	"crud-echo/internal/inbound/routers"
	"crud-echo/internal/outbound/database"
	"crud-echo/internal/repository"
	"crud-echo/internal/usecase"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	db := database.NewPostgresDB()
	if err := db.Migrate(); err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Debug = true

	booksRepo := repository.NewBooksRepository(db.DB)
	validator := usecase.NewCustomValidator(validator.New())
	booksUseCase := usecase.NewBooksUseCase(booksRepo, validator)
	booksHandler := handlers.NewBooksHandler(booksUseCase)
	routers.RegisterRoutes(e, booksHandler)

	e.Logger.Fatal(e.Start(":1323"))
}
