package main

import (
	"crud-echo/internal/inbound/handlers"
	"crud-echo/internal/inbound/routers"
	"crud-echo/internal/outbound/database"
	"crud-echo/internal/repository"
	uc "crud-echo/internal/usecase"
	vc "crud-echo/internal/usecase/validators_custom"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}

	if err := database.Migrate(db); err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Debug = true

	// e.Validator = &vc.CustomValidator{Validator: validator.New()}

	repo := repository.NewBookRepository(db)
	validate2 := &vc.CustomValidator{Validator: validator.New()}
	usecase := uc.NewBookUseCase(repo, validate2)
	booksHandler := &handlers.BooksHandler{BUC: usecase}
	routers.RegisterRoutes(e, booksHandler)

	e.Logger.Fatal(e.Start(":1323"))
}
