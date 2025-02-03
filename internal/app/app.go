package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/aspirin100/finapi/internal/config"
)

type App struct {
	server *http.Server
}

func New(cfg *config.Config) *App {
	router := gin.Default()

	// TODO: storage constructor for service layer

	// TODO: service layer constructor
	// ...New(storage)

	// TODO: handler constructor

	srv := &http.Server{
		Addr:    cfg.Hostname + ":" + cfg.Port,
		Handler: router,
	}

	return &App{
		server: srv,
	}
}

func (app *App) Run() error {
	err := app.server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("failed to start application: %w", err)
	}

	return nil
}

func (app *App) Stop(ctx context.Context) error {
	return nil
}
