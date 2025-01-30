package main

import (
	"context"
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"

	"github.com/jonesrussell/goforms/internal/application/config"
	"github.com/jonesrussell/goforms/internal/application/database"
	"github.com/jonesrussell/goforms/internal/application/handlers"
	"github.com/jonesrussell/goforms/internal/application/logging"
	"github.com/jonesrussell/goforms/internal/application/router"
	"github.com/jonesrussell/goforms/internal/application/view"
	"github.com/jonesrussell/goforms/internal/domain/contact"
	"github.com/jonesrussell/goforms/internal/domain/user"
)

func main() {
	app := fx.New(
		// Config
		config.Module,

		// Logging
		logging.Module,
		fx.WithLogger(func(log logging.Logger) fxevent.Logger {
			return &logging.FxEventLogger{Logger: log}
		}),

		// Database
		database.Module,

		// Provide domain modules
		user.ProvideModule(),
		contact.ProvideModule(),

		// Invoke the server
		fx.Invoke(startAppAndServer),
	)

	app.Run()
}

func startAppAndServer(lc fx.Lifecycle, p ServerParams, logger logging.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			loadEnvironment()

			logger.Debug("starting server with handlers", logging.Int("handler_count", len(p.Handlers)))

			for _, handler := range p.Handlers {
				handler.Register(p.Echo)
			}

			router.Setup(p.Echo, &router.Config{
				Handlers: p.Handlers,
				Static: router.StaticConfig{
					Path: "/static",
					Root: "static",
				},
				Logger: p.Logger,
			})

			addr := fmt.Sprintf("%s:%d", p.Config.Server.Host, p.Config.Server.Port)
			if p.Config.Server.Port == 0 {
				addr = fmt.Sprintf("%s:8090", p.Config.Server.Host)
			}

			logger.Info("Starting server",
				logging.String("addr", addr),
				logging.String("env", p.Config.App.Env),
			)

			return p.Echo.Start(addr)
		},
		OnStop: func(ctx context.Context) error {
			// Handle any cleanup if necessary
			return nil
		},
	})
}

func loadEnvironment() {
	log.Println("Loading environment")
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v\n", err)
		return // Handle the error without panicking
	}
}

type ServerParams struct {
	fx.In

	Echo           *echo.Echo
	Config         *config.Config
	Logger         logging.Logger
	Renderer       *view.Renderer
	ContactService contact.Service
	UserService    user.Service
	Handlers       []handlers.Handler `group:"handlers"`
}
