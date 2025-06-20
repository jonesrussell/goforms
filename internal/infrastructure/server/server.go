package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"

	"github.com/goformx/goforms/internal/application/middleware"
	"github.com/goformx/goforms/internal/infrastructure/config"
	"github.com/goformx/goforms/internal/infrastructure/logging"
	"github.com/goformx/goforms/internal/infrastructure/version"
	"github.com/goformx/goforms/internal/infrastructure/web"
)

const (
	// StartupTimeout is the timeout for server startup
	StartupTimeout = 5 * time.Second
	// ShutdownTimeout is the timeout for graceful shutdown
	ShutdownTimeout = 10 * time.Second
)

// Server handles HTTP server lifecycle and configuration
type Server struct {
	echo   *echo.Echo
	logger logging.Logger
	config *config.Config
	server *http.Server
}

// Address returns the server's address in host:port format
func (s *Server) Address() string {
	return net.JoinHostPort(s.config.App.Host, strconv.Itoa(s.config.App.Port))
}

// URL returns the server's full HTTP URL
func (s *Server) URL() string {
	return fmt.Sprintf("%s://%s", s.config.App.Scheme, s.Address())
}

// Start starts the server and returns when it's ready to accept connections
func (s *Server) Start() error {
	addr := s.Address()
	s.server = &http.Server{
		Addr:              addr,
		Handler:           s.echo,
		ReadTimeout:       s.config.App.ReadTimeout,
		WriteTimeout:      s.config.App.WriteTimeout,
		IdleTimeout:       s.config.App.IdleTimeout,
		ReadHeaderTimeout: s.config.App.ReadTimeout,
	}

	// Create channels for server startup coordination
	started := make(chan struct{})
	errored := make(chan error, 1)

	// Start server in a goroutine
	go func() {
		// Create a listener to check if the server can bind to the port
		listener, err := net.Listen("tcp", addr)
		if err != nil {
			errored <- fmt.Errorf("failed to create listener: %w", err)
			return
		}

		// Signal that the server is ready to accept connections
		close(started)

		// Start serving
		if serveErr := s.server.Serve(listener); serveErr != nil && serveErr != http.ErrServerClosed {
			errored <- fmt.Errorf("server error: %w", serveErr)
		}
	}()

	// Wait for the server to be ready or fail
	select {
	case err := <-errored:
		return fmt.Errorf("server failed to start: %w", err)
	case <-started:
		versionInfo := version.GetInfo()
		s.logger.Info("server started",
			"host", s.config.App.Host,
			"port", s.config.App.Port,
			"environment", s.config.App.Env,
			"version", versionInfo.Version,
			"build_time", versionInfo.BuildTime,
			"git_commit", versionInfo.GitCommit)
		return nil
	case <-time.After(StartupTimeout):
		return errors.New("server startup timed out after 5 seconds")
	}
}

// New creates a new server instance with the provided dependencies
func New(
	lc fx.Lifecycle,
	logger logging.Logger,
	cfg *config.Config,
	e *echo.Echo,
	middlewareManager *middleware.Manager,
	assetServer web.AssetServer,
) *Server {
	srv := &Server{
		echo:   e,
		logger: logger,
		config: cfg,
	}

	// Log server configuration
	logger.Info("initializing server",
		"host", cfg.App.Host,
		"port", cfg.App.Port,
		"environment", cfg.App.Env,
		"server_type", "echo")

	// Add health check endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// Register asset routes
	if err := assetServer.RegisterRoutes(e); err != nil {
		logger.Error("failed to register asset routes", "error", err)
	}

	// Register lifecycle hooks
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return nil // Server will be started after middleware is registered
		},
		OnStop: func(ctx context.Context) error {
			if srv.server == nil {
				return nil
			}

			srv.logger.Info("shutting down server")

			shutdownCtx, cancel := context.WithTimeout(ctx, ShutdownTimeout)
			defer cancel()

			if err := srv.server.Shutdown(shutdownCtx); err != nil {
				srv.logger.Error("server shutdown error", "error", err, "timeout", ShutdownTimeout)
				return fmt.Errorf("server shutdown error: %w", err)
			}

			srv.logger.Info("server stopped gracefully")
			return nil
		},
	})

	return srv
}

// Echo returns the underlying echo instance
func (s *Server) Echo() *echo.Echo {
	return s.echo
}

// Config returns the server configuration
func (s *Server) Config() *config.Config {
	return s.config
}
