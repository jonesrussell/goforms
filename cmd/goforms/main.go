package main

import (
	"context"
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"

	"github.com/jonesrussell/goforms/internal/application"
	"github.com/jonesrussell/goforms/internal/application/config"
	"github.com/jonesrussell/goforms/internal/application/handlers"
	"github.com/jonesrussell/goforms/internal/application/logging"
	"github.com/jonesrussell/goforms/internal/application/repositories"
	"github.com/jonesrussell/goforms/internal/application/router"
	"github.com/jonesrussell/goforms/internal/domain"
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
	if err := run(); err != nil {
		log.Fatal(err)
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
		logging.Module,
		config.Module,
		domain.Module,
		application.Module,
		repositories.Module,
		user.Module,
		fx.Provide(func() handlers.VersionInfo { return versionInfo }),
		fx.Provide(user.NewService),
		fx.Provide(user.NewUserRepository),
		fx.Provide(user.NewTokenRepository),
		fx.Provide(newServer),
		fx.Provide(func(logger logging.Logger, userService user.Service) *handlers.AuthHandler {
			return handlers.NewAuthHandler(logger, userService)
		}),
		fx.Provide(func(logger logging.Logger) *handlers.WebHandler {
			return handlers.NewWebHandler(logger)
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

func newServer(userService user.Service) *echo.Echo {
	e := echo.New()
	// Set up routes and middleware using userService
	return e
}

type ServerParams struct {
	fx.In

	Echo        *echo.Echo
	Config      *config.Config
	Logger      logging.Logger
	AuthHandler *handlers.AuthHandler
	WebHandler  *handlers.WebHandler
}

func startServer(p ServerParams) error {
	p.Logger.Debug("starting server with handlers",
		logging.Int("handler_count", 2),
	)

	p.AuthHandler.Register(p.Echo)
	p.WebHandler.Register(p.Echo)

	router.Setup(p.Echo, &router.Config{
		Handlers: []handlers.Handler{p.AuthHandler, p.WebHandler},
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

	p.Logger.Info("Starting server",
		logging.String("addr", addr),
		logging.String("env", p.Config.App.Env),
		logging.String("version", version),
		logging.String("gitCommit", gitCommit),
	)

	return p.Echo.Start(addr)
}
