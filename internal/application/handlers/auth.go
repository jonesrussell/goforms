package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/jonesrussell/goforms/internal/application/logging"
	"github.com/jonesrussell/goforms/internal/domain/user"
)

// AuthHandlerOption defines an auth handler option
type AuthHandlerOption func(*AuthHandler)

// AuthHandler handles authentication related requests
type AuthHandler struct {
	userService user.Service
	Logger      logging.Logger
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(logger logging.Logger, userService user.Service) *AuthHandler {
	return &AuthHandler{
		Logger:      logger,
		userService: userService,
	}
}

// Validate validates that required dependencies are set
func (h *AuthHandler) Validate() error {
	if h.userService == nil {
		return fmt.Errorf("user service is required")
	}
	if h.Logger == nil {
		return fmt.Errorf("logger is required")
	}
	return nil
}

// Register registers the auth routes
func (h *AuthHandler) Register(e *echo.Echo) {
	if err := h.Validate(); err != nil {
		h.Logger.Error("failed to validate handler", logging.Error(err))
		return
	}

	g := e.Group("/api/v1/auth")
	g.POST("/signup", h.handleSignup)
	g.POST("/login", h.handleLogin)
	g.POST("/logout", h.handleLogout)
}

// handleSignup handles user registration
func (h *AuthHandler) handleSignup(c echo.Context) error {
	var signup user.Signup
	if err := c.Bind(&signup); err != nil {
		h.Logger.LogWithPrefix("error", "auth", "failed to bind request", logging.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request format")
	}

	if err := c.Validate(signup); err != nil {
		h.Logger.LogWithPrefix("error", "auth", "validation failed", logging.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	newUser := user.ConvertSignupToUser(&signup)
	newUser, err := h.userService.SignUp(c.Request().Context(), newUser)
	if err != nil {
		h.Logger.LogWithPrefix("error", "auth", "failed to create user", logging.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	h.Logger.LogWithPrefix("info", "auth", "User created successfully", logging.String("email", signup.Email))
	return c.JSON(http.StatusCreated, newUser)
}

// handleLogin handles user authentication
// @Summary Authenticate user
// @Description Authenticate user and return JWT tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param login body user.Login true "User login credentials"
// @Success 200 {object} user.TokenPair
// @Failure 400 {object} echo.HTTPError
// @Failure 401 {object} echo.HTTPError
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) handleLogin(c echo.Context) error {
	var login user.Login
	if err := c.Bind(&login); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request format")
	}

	if err := c.Validate(login); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	tokens, err := h.userService.Login(c.Request().Context(), &login)
	if err != nil {
		h.Logger.Error("auth: failed to authenticate user", logging.Error(err))
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid credentials")
	}

	return c.JSON(http.StatusOK, tokens)
}

// handleLogout handles user logout
// @Summary Logout user
// @Description Invalidate user's tokens
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]string
// @Failure 401 {object} echo.HTTPError
// @Router /api/v1/auth/logout [post]
func (h *AuthHandler) handleLogout(c echo.Context) error {
	token := c.Request().Header.Get("Authorization")
	if token == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "Missing authorization token")
	}

	if err := h.userService.Logout(c.Request().Context(), token); err != nil {
		h.Logger.Error("auth: failed to logout user", logging.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to logout")
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Successfully logged out"})
}
