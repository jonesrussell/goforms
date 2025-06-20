package middleware

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"golang.org/x/time/rate"

	"github.com/goformx/goforms/internal/application/constants"
	"github.com/goformx/goforms/internal/application/middleware/access"
	"github.com/goformx/goforms/internal/application/middleware/context"
	"github.com/goformx/goforms/internal/application/middleware/session"
	formdomain "github.com/goformx/goforms/internal/domain/form"
	"github.com/goformx/goforms/internal/domain/user"
	appconfig "github.com/goformx/goforms/internal/infrastructure/config"
	"github.com/goformx/goforms/internal/infrastructure/logging"
	"github.com/goformx/goforms/internal/infrastructure/sanitization"
	"github.com/goformx/goforms/internal/infrastructure/version"
)

// Manager handles middleware configuration and setup
type Manager struct {
	logger            logging.Logger
	config            *ManagerConfig
	contextMiddleware *context.Middleware
}

// ManagerConfig represents the configuration for the middleware manager
type ManagerConfig struct {
	Logger         logging.Logger
	Config         *appconfig.Config // Single source of truth
	UserService    user.Service
	FormService    formdomain.Service
	SessionManager *session.Manager
	AccessManager  *access.AccessManager
	Sanitizer      sanitization.ServiceInterface
}

// NewManager creates a new middleware manager
func NewManager(cfg *ManagerConfig) *Manager {
	if cfg == nil {
		panic("config is required")
	}

	if cfg.UserService == nil {
		panic("user service is required")
	}

	if cfg.SessionManager == nil {
		panic("session manager is required")
	}

	if cfg.Sanitizer == nil {
		panic("sanitizer is required")
	}

	return &Manager{
		logger:            cfg.Logger,
		config:            cfg,
		contextMiddleware: context.NewMiddleware(cfg.Logger, cfg.Config.App.RequestTimeout),
	}
}

// GetSessionManager returns the session manager
func (m *Manager) GetSessionManager() *session.Manager {
	return m.config.SessionManager
}

// Setup registers all middleware with the Echo instance
func (m *Manager) Setup(e *echo.Echo) {
	versionInfo := version.GetInfo()
	m.logger.Info("setting up middleware",
		"app", "goforms",
		"version", versionInfo.Version,
		"environment", m.config.Config.App.Env,
		"build_time", versionInfo.BuildTime,
		"git_commit", versionInfo.GitCommit,
	)

	// Set Echo's logger to use our custom logger
	e.Logger = &EchoLogger{logger: m.logger, config: m.config}

	// Enable debug mode and set log level
	e.Debug = m.config.Config.Security.Debug
	if m.config.Config.App.IsDevelopment() {
		e.Logger.SetLevel(log.DEBUG)
		m.logger.Info("development mode enabled",
			"app", "goforms",
			"version", versionInfo.Version,
			"environment", m.config.Config.App.Env,
			"build_time", versionInfo.BuildTime,
			"git_commit", versionInfo.GitCommit)
	} else {
		e.Logger.SetLevel(log.INFO)
	}

	// Add recovery middleware first to catch panics
	e.Use(Recovery(m.logger, m.config.Sanitizer))

	// Add context middleware to handle request context
	e.Use(m.contextMiddleware.WithContext())

	// Add request tracking middleware for debugging
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if m.config.Config.App.IsDevelopment() {
				c.Logger().Debug("Request received",
					"path", c.Request().URL.Path,
					"method", c.Request().Method,
					"user_agent", c.Request().UserAgent())
			}
			return next(c)
		}
	})

	// Register basic middleware
	if m.config.Config.App.IsDevelopment() {
		// Use console format in development
		e.Use(echomw.LoggerWithConfig(echomw.LoggerConfig{
			Format: "${time_rfc3339} ${status} ${method} ${uri} ${latency_human}\n",
			Output: os.Stdout,
		}))
	} else {
		// Use JSON format in production
		e.Use(echomw.Logger())
	}

	// Use PerFormCORS middleware for form-specific CORS handling
	// This middleware will handle CORS for form routes and fallback to global CORS for other routes
	perFormCORSConfig := NewPerFormCORSConfig(m.config.FormService, m.logger, &m.config.Config.Security)
	e.Use(PerFormCORS(perFormCORSConfig))

	// Register security middleware
	e.Use(echomw.SecureWithConfig(echomw.SecureConfig{
		XSSProtection:         m.config.Config.Security.Headers.XXSSProtection,
		ContentTypeNosniff:    m.config.Config.Security.Headers.XContentTypeOptions,
		XFrameOptions:         m.config.Config.Security.Headers.XFrameOptions,
		HSTSMaxAge:            constants.HSTSOneYear,
		HSTSExcludeSubdomains: false,
		ContentSecurityPolicy: m.config.Config.Security.GetCSPDirectives(&m.config.Config.App),
	}))

	// Set security config in context for other middleware
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("security_config", m.config.Config.Security)
			return next(c)
		}
	})

	// Add additional security headers not covered by Echo's Secure middleware
	e.Use(setupAdditionalSecurityHeadersMiddleware())

	if m.config.Config.Security.CSRF.Enabled {
		m.logger.Info("CSRF middleware enabled",
			"enabled", m.config.Config.Security.CSRF.Enabled,
			"token_lookup", m.config.Config.Security.CSRF.TokenLookup)
		csrfMiddleware := setupCSRF(&m.config.Config.Security.CSRF, m.config.Config.App.Env == "development")
		e.Use(csrfMiddleware)
		m.logger.Info("CSRF middleware registered")
	} else {
		m.logger.Info("CSRF middleware disabled")
	}

	// Setup rate limiting using infrastructure config
	if m.config.Config.Security.RateLimit.Enabled {
		e.Use(m.setupRateLimiting())
	}

	// Register session middleware
	m.logger.Info("registering session middleware",
		"app", "goforms",
		"version", versionInfo.Version,
		"environment", m.config.Config.App.Env,
		"build_time", versionInfo.BuildTime,
		"git_commit", versionInfo.GitCommit)
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if m.config.Config.App.IsDevelopment() {
				c.Logger().Debug("Session middleware processing request",
					"path", c.Request().URL.Path,
					"method", c.Request().Method)
			}
			return m.config.SessionManager.Middleware()(next)(c)
		}
	})

	// Register access control middleware
	m.logger.Info("registering access control middleware",
		"app", "goforms",
		"version", versionInfo.Version,
		"environment", m.config.Config.App.Env,
		"build_time", versionInfo.BuildTime,
		"git_commit", versionInfo.GitCommit)
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if m.config.Config.App.IsDevelopment() {
				c.Logger().Debug("Access middleware processing request",
					"path", c.Request().URL.Path,
					"method", c.Request().Method)
			}
			return access.Middleware(m.config.AccessManager, m.logger)(next)(c)
		}
	})

	m.logger.Info("middleware setup completed",
		"app", "goforms",
		"version", versionInfo.Version,
		"environment", m.config.Config.App.Env,
		"build_time", versionInfo.BuildTime,
		"git_commit", versionInfo.GitCommit)
}

// setupCSRF creates and configures CSRF middleware
func setupCSRF(csrfConfig *appconfig.CSRFConfig, isDevelopment bool) echo.MiddlewareFunc {
	sameSite := getSameSite(csrfConfig.CookieSameSite, isDevelopment)
	tokenLength := getTokenLength(csrfConfig.TokenLength)

	return echomw.CSRFWithConfig(echomw.CSRFConfig{
		TokenLength:    uint8(tokenLength), // #nosec G115
		TokenLookup:    csrfConfig.TokenLookup,
		ContextKey:     csrfConfig.ContextKey,
		CookieName:     csrfConfig.CookieName,
		CookiePath:     csrfConfig.CookiePath,
		CookieDomain:   csrfConfig.CookieDomain,
		CookieSecure:   !isDevelopment, // In development, don't require HTTPS
		CookieHTTPOnly: csrfConfig.CookieHTTPOnly,
		CookieSameSite: sameSite,
		CookieMaxAge:   csrfConfig.CookieMaxAge,
		Skipper:        createCSRFSkipper(csrfConfig, isDevelopment),
		ErrorHandler:   createCSRFErrorHandler(csrfConfig, isDevelopment),
	})
}

// getSameSite converts string SameSite to http.SameSite
func getSameSite(cookieSameSite string, isDevelopment bool) http.SameSite {
	switch cookieSameSite {
	case "Lax":
		return http.SameSiteLaxMode
	case "Strict":
		return http.SameSiteStrictMode
	case "None":
		return http.SameSiteNoneMode
	default:
		// In development, default to Lax for cross-origin support
		if isDevelopment {
			return http.SameSiteLaxMode
		}
		return http.SameSiteStrictMode
	}
}

// getTokenLength ensures token length is within bounds for uint8
func getTokenLength(tokenLength int) int {
	if tokenLength <= 0 || tokenLength > 255 {
		return constants.DefaultTokenLength
	}
	return tokenLength
}

// createCSRFSkipper creates the CSRF skipper function
func createCSRFSkipper(csrfConfig *appconfig.CSRFConfig, isDevelopment bool) func(c echo.Context) bool {
	return func(c echo.Context) bool {
		path := c.Request().URL.Path
		method := c.Request().Method

		// Always log when CSRF skipper is called (for debugging)
		if isDevelopment {
			c.Logger().Debug("CSRF skipper called",
				"path", path,
				"method", method)
		}

		// Add debugging in development mode
		if isDevelopment {
			c.Logger().Debug("CSRF middleware processing request",
				"path", path,
				"method", method,
				"token_lookup", csrfConfig.TokenLookup,
				"origin", c.Request().Header.Get("Origin"),
				"csrf_token_present", c.Request().Header.Get("X-CSRF-Token") != "",
				"is_development", isDevelopment,
				"path_signup", constants.PathSignup,
				"path_matches_signup", path == constants.PathSignup,
				"method_post", method == http.MethodPost)
		}

		// Skip CSRF for static files
		if constants.IsStaticFile(path) {
			if isDevelopment {
				c.Logger().Debug("Skipping CSRF for static file", "path", path)
			}
			return true
		}

		// Skip CSRF validation for safe methods, but still generate token
		if method == http.MethodGet || method == http.MethodHead || method == http.MethodOptions {
			if isDevelopment {
				c.Logger().Debug("CSRF middleware will generate token for safe method", "method", method)
			}
			return false
		}

		// Skip CSRF for API endpoints with valid Authorization header or CSRF token
		if strings.HasPrefix(path, "/api/") {
			authHeader := c.Request().Header.Get("Authorization")
			csrfToken := c.Request().Header.Get("X-CSRF-Token")
			if authHeader != "" || csrfToken != "" {
				if isDevelopment {
					c.Logger().Debug("Skipping CSRF for API endpoint with auth/token", "path", path)
				}
				return true
			}
		}

		// Skip CSRF for validation endpoints
		if strings.HasPrefix(path, "/api/validation/") {
			if isDevelopment {
				c.Logger().Debug("Skipping CSRF for validation endpoint", "path", path)
			}
			return true
		}

		// Never skip CSRF for login, signup, or password reset
		if path == constants.PathLogin ||
			path == constants.PathSignup ||
			path == constants.PathResetPassword {
			if isDevelopment {
				c.Logger().Debug("CSRF validation required for auth endpoint", "path", path, "method", method)
			}
			return false
		}

		if isDevelopment {
			c.Logger().Debug("CSRF validation required for endpoint", "path", path, "method", method)
		}
		return false
	}
}

// createCSRFErrorHandler creates the CSRF error handler function
func createCSRFErrorHandler(
	csrfConfig *appconfig.CSRFConfig,
	isDevelopment bool,
) func(err error, c echo.Context) error {
	return func(err error, c echo.Context) error {
		// Add debugging in development mode
		if isDevelopment {
			// Get the actual token from the request
			csrfToken := c.Request().Header.Get("X-CSRF-Token")

			// Get the token from context (if available)
			contextToken := ""
			if token, ok := c.Get(csrfConfig.ContextKey).(string); ok {
				contextToken = token
			}

			// Get cookies for debugging
			cookies := c.Request().Header.Get("Cookie")

			c.Logger().Error("CSRF validation failed",
				"error", err.Error(),
				"path", c.Request().URL.Path,
				"method", c.Request().Method,
				"token_lookup", csrfConfig.TokenLookup,
				"origin", c.Request().Header.Get("Origin"),
				"csrf_token_present", csrfToken != "",
				"csrf_token_length", len(csrfToken),
				"csrf_token_value", csrfToken,
				"context_token_present", contextToken != "",
				"context_token_length", len(contextToken),
				"context_token_value", contextToken,
				"cookies", cookies,
				"content_type", c.Request().Header.Get("Content-Type"),
				"user_agent", c.Request().UserAgent())
		}
		return c.NoContent(http.StatusForbidden)
	}
}

// setupAdditionalSecurityHeadersMiddleware creates and configures additional security headers middleware
func setupAdditionalSecurityHeadersMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get security config from context or use defaults
			securityConfig, ok := c.Get("security_config").(*appconfig.SecurityConfig)
			if !ok {
				// Fallback to default values if config not available
				c.Response().Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
				c.Response().Header().Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
			} else {
				// Use configured values
				c.Response().Header().Set("Referrer-Policy", securityConfig.Headers.ReferrerPolicy)
				c.Response().Header().Set("Strict-Transport-Security", securityConfig.Headers.StrictTransportSecurity)
				c.Response().Header().Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
			}

			return next(c)
		}
	}
}

// EchoLogger implements echo.Logger interface using our custom logger
type EchoLogger struct {
	logger logging.Logger
	config *ManagerConfig
}

func (l *EchoLogger) Print(i ...any) {
	l.logger.Info(fmt.Sprint(i...))
}

func (l *EchoLogger) Printf(format string, args ...any) {
	l.logger.Info(fmt.Sprintf(format, args...))
}

func (l *EchoLogger) Printj(j log.JSON) {
	fields := make([]any, 0, len(j)*constants.FieldPairSize)
	for k, v := range j {
		fields = append(fields, k, fmt.Sprint(v))
	}
	l.logger.Info("", fields...)
}

func (l *EchoLogger) Debug(i ...any) {
	l.logger.Debug(fmt.Sprint(i...))
}

func (l *EchoLogger) Debugf(format string, args ...any) {
	l.logger.Debug(fmt.Sprintf(format, args...))
}

func (l *EchoLogger) Debugj(j log.JSON) {
	fields := make([]any, 0, len(j)*constants.FieldPairSize)
	for k, v := range j {
		fields = append(fields, k, fmt.Sprint(v))
	}
	l.logger.Debug("", fields...)
}

func (l *EchoLogger) Info(i ...any) {
	l.logger.Info(fmt.Sprint(i...))
}

func (l *EchoLogger) Infof(format string, args ...any) {
	l.logger.Info(fmt.Sprintf(format, args...))
}

func (l *EchoLogger) Infoj(j log.JSON) {
	fields := make([]any, 0, len(j)*constants.FieldPairSize)
	for k, v := range j {
		fields = append(fields, k, fmt.Sprint(v))
	}
	l.logger.Info("", fields...)
}

func (l *EchoLogger) Warn(i ...any) {
	l.logger.Warn(fmt.Sprint(i...))
}

func (l *EchoLogger) Warnf(format string, args ...any) {
	l.logger.Warn(fmt.Sprintf(format, args...))
}

func (l *EchoLogger) Warnj(j log.JSON) {
	fields := make([]any, 0, len(j)*constants.FieldPairSize)
	for k, v := range j {
		fields = append(fields, k, fmt.Sprint(v))
	}
	l.logger.Warn("", fields...)
}

func (l *EchoLogger) Error(i ...any) {
	l.logger.Error(fmt.Sprint(i...))
}

func (l *EchoLogger) Errorf(format string, args ...any) {
	l.logger.Error(fmt.Sprintf(format, args...))
}

func (l *EchoLogger) Errorj(j log.JSON) {
	fields := make([]any, 0, len(j)*constants.FieldPairSize)
	for k, v := range j {
		fields = append(fields, k, fmt.Sprint(v))
	}
	l.logger.Error("", fields...)
}

func (l *EchoLogger) Fatal(i ...any) {
	l.logger.Fatal(fmt.Sprint(i...))
}

func (l *EchoLogger) Fatalf(format string, args ...any) {
	l.logger.Fatal(fmt.Sprintf(format, args...))
}

func (l *EchoLogger) Fatalj(j log.JSON) {
	fields := make([]any, 0, len(j)*constants.FieldPairSize)
	for k, v := range j {
		fields = append(fields, k, fmt.Sprint(v))
	}
	l.logger.Fatal("", fields...)
}

func (l *EchoLogger) Panic(i ...any) {
	l.logger.Error(fmt.Sprint(i...))
	panic(fmt.Sprint(i...))
}

func (l *EchoLogger) Panicf(format string, args ...any) {
	l.logger.Error(fmt.Sprintf(format, args...))
	panic(fmt.Sprintf(format, args...))
}

func (l *EchoLogger) Panicj(j log.JSON) {
	fields := make([]any, 0, len(j)*constants.FieldPairSize)
	for k, v := range j {
		fields = append(fields, k, fmt.Sprint(v))
	}
	l.logger.Error("", fields...)
	panic(fmt.Sprintf("%v", j))
}

func (l *EchoLogger) Level() log.Lvl {
	return log.INFO
}

func (l *EchoLogger) SetLevel(level log.Lvl) {
	// No-op as we use our own log level configuration
}

func (l *EchoLogger) SetHeader(h string) {
	// No-op as we use our own log format
}

func (l *EchoLogger) SetPrefix(p string) {
	// No-op as we use our own log format
}

func (l *EchoLogger) Prefix() string {
	return ""
}

func (l *EchoLogger) SetOutput(w io.Writer) {
	// No-op as we use our own output configuration
}

func (l *EchoLogger) Output() io.Writer {
	return os.Stdout
}

// setupRateLimiting creates and configures rate limiting middleware using infrastructure config
func (m *Manager) setupRateLimiting() echo.MiddlewareFunc {
	rateLimitConfig := m.config.Config.Security.RateLimit

	return echomw.RateLimiterWithConfig(echomw.RateLimiterConfig{
		Skipper: func(c echo.Context) bool {
			path := c.Request().URL.Path
			method := c.Request().Method

			// Skip paths from config
			for _, skipPath := range rateLimitConfig.SkipPaths {
				if strings.HasPrefix(path, skipPath) {
					return true
				}
			}

			// Skip methods from config
			for _, skipMethod := range rateLimitConfig.SkipMethods {
				if method == skipMethod {
					return true
				}
			}

			return false
		},
		Store: echomw.NewRateLimiterMemoryStoreWithConfig(
			echomw.RateLimiterMemoryStoreConfig{
				Rate:      rate.Limit(rateLimitConfig.Requests),
				Burst:     rateLimitConfig.Burst,
				ExpiresIn: rateLimitConfig.Window,
			},
		),
		IdentifierExtractor: func(c echo.Context) (string, error) {
			// For login and signup pages, use IP address as identifier
			path := c.Request().URL.Path
			if path == constants.PathLogin || path == constants.PathSignup || path == constants.PathResetPassword {
				return c.RealIP(), nil
			}

			// For form submissions, use form ID and origin
			formID := c.Param("formID")
			origin := c.Request().Header.Get("Origin")
			if formID == "" {
				formID = constants.DefaultUnknown
			}
			if origin == "" {
				origin = constants.DefaultUnknown
			}
			return fmt.Sprintf("%s:%s", formID, origin), nil
		},
		ErrorHandler: func(c echo.Context, err error) error {
			return echo.NewHTTPError(http.StatusTooManyRequests,
				"Rate limit exceeded: too many requests from the same form or origin")
		},
		DenyHandler: func(c echo.Context, identifier string, err error) error {
			return echo.NewHTTPError(http.StatusTooManyRequests,
				"Rate limit exceeded: please try again later")
		},
	})
}
