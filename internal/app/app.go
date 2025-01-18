package app

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"

	"github.com/jonesrussell/goforms/internal/app/server"
	"github.com/jonesrussell/goforms/internal/config"
	"github.com/jonesrussell/goforms/internal/handlers"
	"github.com/jonesrussell/goforms/internal/logger"
	"github.com/jonesrussell/goforms/internal/middleware"
)

type App struct {
	server     *server.Server
	middleware *middleware.Manager
	handlers   *handlers.SubscriptionHandler
	logger     logger.Logger
}

func NewApp(
	lc fx.Lifecycle,
	log logger.Logger,
	echo *echo.Echo,
	cfg *config.Config,
	handler *handlers.SubscriptionHandler,
	healthHandler *handlers.HealthHandler,
	contactHandler *handlers.ContactHandler,
	marketingHandler *handlers.MarketingHandler,
) *App {
	mw := middleware.New(log, cfg)
	srv := server.New(lc, echo, log, &cfg.Server)

	app := &App{
		server:     srv,
		middleware: mw,
		handlers:   handler,
		logger:     log,
	}

	// Setup order: middleware -> handlers
	mw.Setup(echo)
	marketingHandler.Register(echo)
	handler.Register(echo)
	healthHandler.Register(echo)
	contactHandler.Register(echo)

	return app
}

// RegisterHooks sets up the application hooks
func RegisterHooks(app *App) {
	app.logger.Info("Application started successfully")
}
