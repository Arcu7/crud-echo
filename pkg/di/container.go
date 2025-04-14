package di

import (
	"crud-echo/internal/config"
	"crud-echo/internal/inbound/customvalidator"
	"crud-echo/internal/inbound/handlers"
	"crud-echo/internal/inbound/routers"
	"crud-echo/internal/inbound/server"
	"crud-echo/internal/outbound/database"
	"crud-echo/internal/usecase"
	"crud-echo/pkg/postgres"

	"github.com/go-playground/validator/v10"
	"go.uber.org/dig"
)

func BuildContainer(configPath string) (*dig.Container, error) {
	container := dig.New()

	// config
	if err := container.Provide(func() (*config.Config, error) {
		return config.LoadConfig(configPath)
	}); err != nil {
		return nil, err
	}

	// db
	if err := container.Provide(postgres.NewDB, dig.As(new(database.RepositoryDBConn))); err != nil {
		return nil, err
	}

	// server
	if err := container.Provide(server.NewServer); err != nil {
		return nil, err
	}

	// repo
	// using dig.As to implement/specify the interface
	if err := container.Provide(database.NewBooksRepository, dig.As(new(usecase.UsecaseBooksRepository))); err != nil {
		return nil, err
	}

	// usecase
	if err := container.Provide(usecase.NewBooksUseCase, dig.As(new(handlers.HandlerBookUsecase))); err != nil {
		return nil, err
	}

	// custom validator
	if err := container.Provide(func() *validator.Validate {
		return validator.New()
	}); err != nil {
		return nil, err
	}
	if err := container.Provide(customvalidator.NewCustomValidator); err != nil {
		return nil, err
	}

	// router
	if err := container.Provide(routers.NewRouter); err != nil {
		return nil, err
	}

	// handlers
	if err := container.Provide(handlers.NewBooksHandler); err != nil {
		return nil, err
	}

	return container, nil
}
