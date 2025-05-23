// Package main is the entry point for the GoFormX application.
// It sets up the application using dependency injection (via fx),
// configures the server, and manages the application lifecycle.
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"

	"github.com/goformx/goforms/internal/application/handler"
	"github.com/goformx/goforms/internal/application/middleware"
	"github.com/goformx/goforms/internal/domain"
	"github.com/goformx/goforms/internal/domain/user"
	"github.com/goformx/goforms/internal/handlers"
	"github.com/goformx/goforms/internal/infrastructure"
	"github.com/goformx/goforms/internal/infrastructure/config"
	"github.com/goformx/goforms/internal/infrastructure/logging"
	"github.com/goformx/goforms/internal/infrastructure/server"
	"github.com/goformx/goforms/internal/infrastructure/version"
	"github.com/goformx/goforms/internal/infrastructure/web"
	"github.com/goformx/goforms/internal/presentation/view"
)

const (
	// ShutdownTimeout is the maximum time to wait for graceful shutdown
	ShutdownTimeout = 5 * time.Second
)

// main is the entry point of the application.
func main() {
	var appLogger logging.Logger

	// Create the application with fx
	app := fx.New(
		// Core dependencies that are required for basic functionality
		fx.Provide(
			GetVersion,
			func() (logging.Logger, error) {
				factory := logging.NewFactory()
				return factory.CreateLogger()
			},
			func(logger logging.Logger) *zap.Logger {
				if zapLogger, ok := logger.(*logging.ZapLogger); ok {
					return zapLogger.GetZapLogger()
				}
				// Fallback to development logger if type assertion fails
				devLogger, _ := zap.NewDevelopment()
				return devLogger
			},
		),
		// Infrastructure module for database, cache, etc.
		infrastructure.RootModule,
		// Domain module containing business logic
		domain.Module,
		// View module for template rendering
		view.Module,
		// Server setup with Echo framework
		fx.Provide(
			newServer,
		),
		// Custom logger for fx events
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
		fx.Invoke(web.InitializeAssets),
		// Start the server using fx.Invoke
		fx.Invoke(startServer),
		fx.Invoke(func(logger logging.Logger) {
			appLogger = logger
			// Application is ready
			logger.Info("Application started")
		}),
	)

	// Start the application
	if err := app.Start(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start application: %v\n", err)
		return
	}

	// Set up signal handling
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for interrupt signal
	sig := <-signalChan
	if appLogger != nil {
		appLogger.Info("Received shutdown signal", logging.StringField("signal", sig.String()))
	}

	// Shutdown the application with timeout
	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), ShutdownTimeout)
	defer cancelShutdown()

	// Start graceful shutdown
	if err := app.Stop(shutdownCtx); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to stop application: %v\n", err)
		return
	}
}

// newServer creates and configures a new Echo server instance.
// It sets up middleware, logging, and security features.
func newServer(
	cfg *config.Config,
	userService user.Service,
	log logging.Logger,
) (
	*echo.Echo,
	*middleware.Manager,
	error,
) {
	// Initialize Echo server
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	// Set up request validation
	e.Validator = middleware.NewValidator()

	// Configure middleware stack using Manager pattern
	mwManager := middleware.New(&middleware.ManagerConfig{
		Logger:      log,
		UserService: userService,
		Security:    &cfg.Security,
	})
	mwManager.Setup(e)

	return e, mwManager, nil
}

// ServerParams contains the dependencies required for starting the server.
// It uses fx.In to automatically inject dependencies.
type ServerParams struct {
	fx.In

	Server            *server.Server
	Config            *config.Config
	Logger            logging.Logger
	Handlers          []handlers.Handler `group:"handlers"`
	MiddlewareManager *middleware.Manager
}

// startServer registers all handlers with the server.
// It uses fx.In to automatically inject dependencies.
func startServer(params ServerParams) error {
	// Register all handlers with the middleware manager
	for _, h := range params.Handlers {
		if webHandler, ok := h.(*handler.WebHandler); ok {
			handler.WithMiddlewareManager(params.MiddlewareManager)(webHandler)
		}
		h.Register(params.Server.Echo())
	}

	return nil
}

// GetVersion returns the version information
func GetVersion() version.Info {
	return version.GetInfo()
}
