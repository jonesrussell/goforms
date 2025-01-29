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
	"github.com/jonesrussell/goforms/internal/application/database"
	"github.com/jonesrussell/goforms/internal/application/handlers"
	"github.com/jonesrussell/goforms/internal/application/logging"
	"github.com/jonesrussell/goforms/internal/application/middleware"
	"github.com/jonesrussell/goforms/internal/application/router"
	"github.com/jonesrussell/goforms/internal/application/validator"
	"github.com/jonesrussell/goforms/internal/domain"
	"github.com/jonesrussell/goforms/internal/domain/contact"
	"github.com/jonesrussell/goforms/internal/domain/user"
	"github.com/jonesrussell/goforms/internal/presentation/view"
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
		database.Module,
		domain.Module,
		application.Module,
		fx.Provide(newServer),
		fx.Provide(func() handlers.VersionInfo { return versionInfo }),
		fx.Provide(
			func(logger logging.Logger, renderer *view.Renderer, contactService contact.Service) handlers.Handler {
				h := handlers.NewWebHandler(logger, handlers.WithRenderer(renderer), handlers.WithContactService(contactService))
				logger.Debug("WebHandler created successfully")
				return h
			},
		),
		fx.Invoke(func(log logging.Logger) {
			log.Debug("checking module initialization")
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

func newServer(cfg *config.Config, logFactory *logging.Factory, userService *user.Service) (*echo.Echo, error) {
	logger := logFactory.CreateFromConfig()

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	e.Validator = validator.NewValidator()

	middleware.Setup(e, &middleware.Config{
		Logger:      logger,
		JWTSecret:   cfg.Security.JWTSecret,
		UserService: userService,
		EnableCSRF:  cfg.Security.CSRF.Enabled,
	})

	return e, nil
}

type ServerParams struct {
	fx.In

	Echo     *echo.Echo
	Config   *config.Config
	Logger   logging.Logger
	Handlers []handlers.Handler `group:"handlers"`
}

func startServer(p ServerParams) error {
	p.Logger.Debug("starting server with handlers",
		logging.Int("handler_count", len(p.Handlers)),
	)

	for _, h := range p.Handlers {
		h.Register(p.Echo)
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

	p.Logger.Info("Starting server",
		logging.String("addr", addr),
		logging.String("env", p.Config.App.Env),
		logging.String("version", version),
		logging.String("gitCommit", gitCommit),
	)

	return p.Echo.Start(addr)
}
