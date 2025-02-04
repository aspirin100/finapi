package app

import (
	"context"
	"fmt"

	"github.com/aspirin100/finapi/internal/config"
	"github.com/aspirin100/finapi/internal/handler"
	"github.com/aspirin100/finapi/internal/repository"
	"github.com/aspirin100/finapi/internal/service"
)

type App struct {
	requestHandler *handler.Handler
}

func New(ctx context.Context, cfg *config.Config) (*App, error) {
	repo, err := repository.NewConnection(ctx, cfg.PostgresDSN)
	if err != nil {
		return nil, fmt.Errorf("failed to create app instance: %w", err)
	}

	err = repo.UpMigrations("postgres", cfg.PostgresDSN)
	if err != nil {
		return nil, fmt.Errorf("failed to up migrations: %w", err)
	}

	srvc := service.New(repo)

	requestHandler := handler.New(cfg.Hostname, cfg.Port, srvc)

	return &App{
		requestHandler: requestHandler,
	}, nil
}

func (app *App) Run() error {
	err := app.requestHandler.Server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("failed to start application: %w", err)
	}

	return nil
}

func (app *App) Stop(ctx context.Context) error {
	return nil
}
