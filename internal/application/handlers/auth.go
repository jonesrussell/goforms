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
	UserService *user.Service
	Logger      logging.Logger
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(logger logging.Logger, userService *user.Service) *AuthHandler {
	return &AuthHandler{
		Logger:      logger,
		UserService: userService,
	}
}

// Validate validates that required dependencies are set
func (h *AuthHandler) Validate() error {
	if h.UserService == nil {
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
		h.Logger.Error("Failed to bind signup data", logging.Error(err))
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid input"})
	}

	h.Logger.Debug("Received signup data", logging.Any("signup", signup))

	existingUser, err := h.UserService.GetByEmail(signup.Email)
	if err != nil {
		h.Logger.Error("Failed to check existing user", logging.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Internal server error"})
	}

	if existingUser != nil {
		return c.JSON(http.StatusConflict, map[string]string{"message": "Email already in use"})
	}

	createdUser, err := h.UserService.SignUp(c.Request().Context(), &signup)
	if err != nil {
		h.Logger.Error("Failed to sign up user", logging.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to create user"})
	}

	return c.JSON(http.StatusCreated, createdUser)
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

	tokens, err := h.UserService.Login(c.Request().Context(), &login)
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

	if err := h.UserService.Logout(c.Request().Context(), token); err != nil {
		h.Logger.Error("auth: failed to logout user", logging.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to logout")
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Successfully logged out"})
}
