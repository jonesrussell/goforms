package main

import (
	"context"
	"log"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"

	"github.com/jonesrussell/goforms/internal/application"
	"github.com/jonesrussell/goforms/internal/application/config"
	"github.com/jonesrussell/goforms/internal/application/handlers"
	"github.com/jonesrussell/goforms/internal/application/logging"
	"github.com/jonesrussell/goforms/internal/application/loggingconfig"
	"github.com/jonesrussell/goforms/internal/application/repositories"
	"github.com/jonesrussell/goforms/internal/application/view"
	"github.com/jonesrussell/goforms/internal/domain"
	"github.com/jonesrussell/goforms/internal/domain/contact"
	"github.com/jonesrussell/goforms/internal/domain/user"
)

//nolint:gochecknoglobals // These variables are populated by -ldflags at build time
var (
	version   = "dev"
	buildTime = "unknown"
	gitCommit = "unknown"
	goVersion = "unknown"
)

func main() {
	app := fx.New(
		fx.Provide(
			loggingconfig.NewConfig,
			logging.NewLogger,
			func() handlers.VersionInfo {
				return createVersionInfo()
			},
		),
		logging.Module,
		config.Module,
		domain.Module,
		application.Module,
		repositories.Module,
		user.Module,
		fx.Provide(
			user.NewService,
			user.NewUserRepository,
			user.NewTokenRepository,
		),
		fx.Provide(
			fx.Annotate(func(h *handlers.AuthHandler) handlers.Handler { return h }, fx.As(new(handlers.Handler))),
			fx.Annotate(func(h *handlers.WebHandler) handlers.Handler { return h }, fx.As(new(handlers.Handler))),
		),
		fx.WithLogger(func(log logging.Logger) fxevent.Logger {
			return &logging.FxEventLogger{Logger: log}
		}),
		fx.Invoke(startApp),
	)

	app.Run()
}

func startApp(lc fx.Lifecycle) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			loadEnvironment()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})
}

func loadEnvironment() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v\n", err)
	}
}

func createVersionInfo() handlers.VersionInfo {
	return handlers.VersionInfo{
		Version:   version,
		BuildTime: buildTime,
		GitCommit: gitCommit,
		GoVersion: goVersion,
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

// func startServer(p ServerParams, logger logging.Logger) error {
// 	logger.Debug("starting server with handlers", logging.Int("handler_count", len(p.Handlers)))

// 	for _, handler := range p.Handlers {
// 		handler.Register(p.Echo)
// 	}

// 	router.Setup(p.Echo, &router.Config{
// 		Handlers: p.Handlers,
// 		Static: router.StaticConfig{
// 			Path: "/static",
// 			Root: "static",
// 		},
// 		Logger: p.Logger,
// 	})

// 	addr := fmt.Sprintf("%s:%d", p.Config.Server.Host, p.Config.Server.Port)
// 	if p.Config.Server.Port == 0 {
// 		addr = fmt.Sprintf("%s:8090", p.Config.Server.Host)
// 	}

// 	logger.Info("Starting server",
// 		logging.String("addr", addr),
// 		logging.String("env", p.Config.App.Env),
// 		logging.String("version", version),
// 		logging.String("gitCommit", gitCommit),
// 	)

// 	return p.Echo.Start(addr)
// }
