package main

import (
	"crud-echo/internal/inbound/handlers"
	"crud-echo/internal/inbound/routers"
	"crud-echo/internal/outbound/database"
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

	e.Validator = &vc.CustomValidator{Validator: validator.New()}

	booksHandler := &handlers.BooksHandler{DB: db}
	routers.RegisterRoutes(e, booksHandler)

	e.Logger.Fatal(e.Start(":1323"))
}
