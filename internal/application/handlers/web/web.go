package web

import (
	"context"
	"errors"
	"net/http"

	"github.com/goformx/goforms/internal/application/constants"
	"github.com/goformx/goforms/internal/application/middleware/auth"
	mwcontext "github.com/goformx/goforms/internal/application/middleware/context"
	"github.com/goformx/goforms/internal/presentation/templates/pages"
	"github.com/labstack/echo/v4"
)

const (
	// StatusFound is the HTTP status code for redirects
	StatusFound = http.StatusFound // 302
)

// WebHandler handles web page requests
type WebHandler struct {
	*BaseHandler
	AuthMiddleware *auth.Middleware
}

// NewWebHandler creates a new web handler using BaseHandler
func NewWebHandler(base *BaseHandler, authMiddleware *auth.Middleware) (*WebHandler, error) {
	if base == nil {
		return nil, errors.New("base handler cannot be nil")
	}

	if authMiddleware == nil {
		return nil, errors.New("auth middleware cannot be nil")
	}

	return &WebHandler{
		BaseHandler:    base,
		AuthMiddleware: authMiddleware,
	}, nil
}

// Register registers the web routes
func (h *WebHandler) Register(e *echo.Echo) {
	e.GET("/", h.handleHome)
	e.GET("/demo", h.handleDemo)
}

// handleHome handles the home page request
func (h *WebHandler) handleHome(c echo.Context) error {
	data := h.BuildPageData(c, "Home")
	if h.Logger != nil {
		h.Logger.Debug("handleHome: data.User", "user", data.User)
	}

	// Check if user is authenticated and redirect to dashboard
	if mwcontext.IsAuthenticated(c) {
		return c.Redirect(constants.StatusSeeOther, constants.PathDashboard)
	}

	// User is not authenticated, render home page
	if renderErr := h.Renderer.Render(c, pages.Home(data)); renderErr != nil {
		return h.HandleError(c, renderErr, "Failed to render home page")
	}
	return nil
}

// handleDemo handles the demo page request
func (h *WebHandler) handleDemo(c echo.Context) error {
	data := h.BuildPageData(c, "Demo")
	if h.Logger != nil {
		h.Logger.Debug("handleDemo: data.User", "user", data.User)
	}

	// Check if user is authenticated and redirect to dashboard
	if mwcontext.IsAuthenticated(c) {
		return c.Redirect(constants.StatusSeeOther, constants.PathDashboard)
	}

	// User is not authenticated, render demo page
	return h.Renderer.Render(c, pages.Demo(data))
}

// Start initializes the web handler.
// This is called during application startup.
func (h *WebHandler) Start(ctx context.Context) error {
	return nil // No initialization needed
}

// Stop cleans up any resources used by the web handler.
// This is called during application shutdown.
func (h *WebHandler) Stop(ctx context.Context) error {
	return nil // No cleanup needed
}
