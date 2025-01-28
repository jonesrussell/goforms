package main

import (
	"context"
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"

	"github.com/jonesrussell/goforms/internal/application/config"
	"github.com/jonesrussell/goforms/internal/application/database"
	"github.com/jonesrussell/goforms/internal/application/handlers"
	"github.com/jonesrussell/goforms/internal/application/logging"
	"github.com/jonesrussell/goforms/internal/application/middleware"
	"github.com/jonesrussell/goforms/internal/application/router"
	"github.com/jonesrussell/goforms/internal/application/validator"
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
	// Load environment
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v\n", err)
	}

	// Create version info
	versionInfo := handlers.VersionInfo{
		Version:   version,
		BuildTime: buildTime,
		GitCommit: gitCommit,
		GoVersion: goVersion,
	}

	// Create app with DI
	app := fx.New(
		logging.Module,
		config.Module,
		database.Module,
		user.Module,
		contact.Module,
		// Provide the Echo instance
		fx.Provide(
			newServer, // This will provide *echo.Echo
		),
		// Version info first
		fx.Provide(
			func() handlers.VersionInfo {
				return versionInfo
			},
		),
		// View module for rendering
		view.Module,
		// Add debug logging for dependency injection
		fx.Invoke(func(log logging.Logger) {
			log.Debug("checking module initialization")
		}),
		// Start server last
		fx.Invoke(startServer), // Ensure the server is started
	)

	// Run app
	ctx := context.Background()
	if err := app.Start(ctx); err != nil {
		return fmt.Errorf("failed to start application: %w", err)
	}

	<-app.Done()
	if err := app.Stop(ctx); err != nil {
		return fmt.Errorf("failed to stop application: %w", err)
	}
	return nil
}

func newServer(cfg *config.Config, logFactory *logging.Factory, userService user.Service) (*echo.Echo, error) {
	// Create logger
	logger := logFactory.CreateFromConfig()

	// Create Echo instance
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	// Register validator
	e.Validator = validator.NewValidator()

	// Configure middleware
	middleware.Setup(e, &middleware.Config{
		Logger:      logger,
		JWTSecret:   cfg.Security.JWTSecret,
		UserService: userService,
		EnableCSRF:  cfg.Security.CSRF.Enabled,
	})

	return e, nil
}

// ServerParams contains the dependencies for starting the server
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

	for i, h := range p.Handlers {
		p.Logger.Debug("handler available",
			logging.Int("index", i),
			logging.String("type", fmt.Sprintf("%T", h)),
		)
	}

	// Configure routes
	router.Setup(p.Echo, &router.Config{
		Handlers: p.Handlers,
		Static: router.StaticConfig{
			Path: "/static",
			Root: "static",
		},
		Logger: p.Logger,
	})

	// Start server
	addr := fmt.Sprintf("%s:%d", p.Config.Server.Host, p.Config.Server.Port)
	if p.Config.Server.Port == 0 {
		addr = fmt.Sprintf("%s:8090", p.Config.Server.Host) // Default to 8090 if port is not set
	}

	p.Logger.Info("Starting server",
		logging.String("addr", addr),
		logging.String("env", p.Config.App.Env),
		logging.String("version", version),
		logging.String("gitCommit", gitCommit),
	)

	return p.Echo.Start(addr)
}
