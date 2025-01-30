package main

import (
	"context"
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"

	"github.com/jonesrussell/goforms/internal/application"
	"github.com/jonesrussell/goforms/internal/application/config"
	"github.com/jonesrussell/goforms/internal/application/handlers"
	"github.com/jonesrussell/goforms/internal/application/logging"
	"github.com/jonesrussell/goforms/internal/application/loggingconfig"
	"github.com/jonesrussell/goforms/internal/application/repositories"
	"github.com/jonesrussell/goforms/internal/application/router"
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

var appLogger logging.Logger

func main() {
	cfg := loggingconfig.NewConfig()
	appLogger = logging.NewLogger(cfg)
	defer appLogger.Sync()

	if err := run(); err != nil {
		appLogger.Fatal("Application failed to start", zap.Error(err))
	}
}

func run() error {
	loadEnvironment()
	versionInfo := createVersionInfo()
	app := createApp(versionInfo)
	return startApp(app)
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

func createApp(versionInfo handlers.VersionInfo) *fx.App {
	return fx.New(
		fx.Provide(
			func() handlers.VersionInfo {
				return versionInfo
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
		fx.Invoke(startServer),
	)
}

func startApp(app *fx.App) error {
	ctx := context.Background()
	if err := app.Start(ctx); err != nil {
		return fmt.Errorf("failed to start application: %w", err)
	}

	if err := app.Stop(ctx); err != nil {
		log.Printf("Error stopping application: %v", err)
	}
	<-app.Done()
	return nil
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

func startServer(p ServerParams) error {
	appLogger.Debug("starting server with handlers", zap.Int("handler_count", len(p.Handlers)))

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

	appLogger.Info("Starting server",
		zap.String("addr", addr),
		zap.String("env", p.Config.App.Env),
		zap.String("version", version),
		zap.String("gitCommit", gitCommit),
	)

	return p.Echo.Start(addr)
}
