package handlers

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"

	"github.com/jonesrussell/goforms/internal/application/logging"
	"github.com/jonesrussell/goforms/internal/domain/user"
)

// AuthHandlerOption defines an auth handler option
type AuthHandlerOption func(*AuthHandler)

// AuthHandler handles authentication related requests
type AuthHandler struct {
	UserService user.Service
	Logger      logging.Logger
	validate    *validator.Validate
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(logger logging.Logger, userService user.Service) *AuthHandler {
	v := validator.New()
	return &AuthHandler{
		Logger:      logger,
		UserService: userService,
		validate:    v,
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
	var signupRequest user.Signup

	// Log the raw request body for debugging
	bodyBytes, err := io.ReadAll(c.Request().Body)
	if err != nil {
		h.Logger.Error("Failed to read request body", logging.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Internal server error"})
	}
	h.Logger.Debug("Raw request body", logging.Any("body", string(bodyBytes)))

	// Reset the request body so it can be read again
	c.Request().Body = io.NopCloser(bytes.NewReader(bodyBytes))

	if err := c.Bind(&signupRequest); err != nil {
		h.Logger.Error("Failed to bind signup data", logging.Error(err))
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid input"})
	}

	// Log the received signup data
	h.Logger.Debug("Received signup data", logging.Any("signup", signupRequest))

	// Hash the password before saving
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(signupRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		h.Logger.Error("Failed to hash password", logging.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Internal server error"})
	}
	signupRequest.Password = string(hashedPassword) // Replace the plain password with the hashed one

	// Check if the email already exists
	existingUser, err := h.UserService.GetByEmail(signupRequest.Email)
	if err != nil {
		h.Logger.Error("Failed to check existing user", logging.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Internal server error"})
	}

	if existingUser != nil {
		return c.JSON(http.StatusConflict, map[string]string{"message": "Email already in use"})
	}

	// Pass the Signup struct to the SignUp method
	createdUser, err := h.UserService.SignUp(&signupRequest)
	if err != nil {
		h.Logger.Error("Failed to sign up user", logging.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to create user"})
	}

	h.Logger.Debug("User signed up successfully", logging.Any("createdUser", createdUser))
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"success": true,
		"message": "Signup successful!",
		"user":    createdUser,
	})
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
	var loginRequest user.Login

	if err := c.Bind(&loginRequest); err != nil {
		h.Logger.Error("Failed to bind login data", logging.Error(err))
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid input"})
	}

	h.Logger.Debug("Received login data", logging.Any("login", loginRequest))

	// Validate the loginRequest here
	if err := h.validate.Struct(loginRequest); err != nil {
		h.Logger.Error("Validation failed", logging.Error(err))
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid input"})
	}

	tokens, err := h.UserService.Login(c.Request().Context(), &loginRequest)
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
