package main

import (
	"crud-echo/internal/inbound/routers"
	"crud-echo/internal/inbound/server"
	"crud-echo/internal/outbound/database"
	"crud-echo/pkg/di"
	"log"
)

func main() {
	configPath := "../../internal/config/"

	container, err := di.BuildContainer(configPath)
	if err != nil {
		log.Fatal("container error:", err)
	}

	if err := container.Invoke(func(dbConn database.RepositoryDBConn) {
		if err := dbConn.Migrate(); err != nil {
			log.Fatal("migrate error:", err)
		}
	}); err != nil {
		log.Fatal("migrate invoke error:", err)
	}

	if err := container.Invoke(func(router *routers.Router) {
		router.RegisterRoutes()
	}); err != nil {
		log.Fatal("router invoke error:", err)
	}

	if err := container.Invoke(func(server *server.Server) error {
		return server.Start()
	}); err != nil {
		log.Fatal("server invoke error:", err)
	}

	// e := echo.New()
	// // e.Use(middleware.Logger())
	//
	// logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	// e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
	// 	LogStatus: true,
	// 	LogURI:    true,
	// 	LogError:  true,
	// 	LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
	// 		if v.Error == nil {
	// 			logger.LogAttrs(context.Background(), slog.LevelInfo, "REQUEST",
	// 				slog.String("uri", v.URI),
	// 				slog.Int("status", v.Status),
	// 			)
	// 		} else {
	// 			logger.LogAttrs(context.Background(), slog.LevelError, "REQUEST_ERROR",
	// 				slog.String("uri", v.URI),
	// 				slog.Int("status", v.Status),
	// 				slog.String("err", v.Error.Error()),
	// 			)
	// 		}
	// 		return nil
	// 	},
	// }))
	//
	// e.Use(middleware.Recover())
	// e.HTTPErrorHandler = handlers.CustomHTTPErrorHandler
	//
	// e.Debug = true

	// booksRepo := database.NewBooksRepository(db.DB)
	// booksValidator := customvalidator.NewCustomValidator(validator.New())
	// booksUseCase := usecase.NewBooksUseCase(booksRepo)
	// booksHandler := handlers.NewBooksHandler(booksUseCase, booksValidator)
	// routers.RegisterRoutes(e, booksHandler)
	//
	// e.Logger.Fatal(e.Start(":1323"))
}
