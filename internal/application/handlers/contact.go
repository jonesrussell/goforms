package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/jonesrussell/goforms/internal/application/logging"
	"github.com/jonesrussell/goforms/internal/domain/common"
	"github.com/jonesrussell/goforms/internal/domain/contact"
)

// ContactHandler handles contact form submissions. It implements the Handler interface
// and follows the standard handler patterns:
//
// 1. Embeds Base for common functionality
// 2. Uses options pattern for dependency injection
// 3. Validates all required dependencies
// 4. Provides clear API documentation
//
// Example Usage:
//
//	handler := NewContactHandler(logger,
//	    WithContactServiceOpt(contactService),
//	)
//	handler.Register(e) // Register routes with Echo
type ContactHandler struct {
	Base
	contactService contact.Service
}

// ContactHandlerOption configures a ContactHandler.
type ContactHandlerOption func(*ContactHandler)

// WithContactServiceOpt sets the contact service for the handler.
// This is a required dependency for handling contact form operations.
func WithContactServiceOpt(svc contact.Service) ContactHandlerOption {
	return func(h *ContactHandler) {
		h.contactService = svc
	}
}

// NewContactHandler creates a new ContactHandler with the provided options.
// The logger is required and must be provided. The contact service must be
// provided using WithContactServiceOpt.
func NewContactHandler(logger logging.Logger, opts ...ContactHandlerOption) *ContactHandler {
	h := &ContactHandler{
		Base: NewBase(WithLogger(logger)), // Use WithLogger from base.go
	}

	for _, opt := range opts {
		opt(h)
	}

	return h
}

// Validate ensures all required dependencies are properly set.
// It checks both the base handler requirements and ContactHandler-specific
// dependencies.
func (h *ContactHandler) Validate() error {
	if err := h.Base.Validate(); err != nil {
		return err
	}
	if h.contactService == nil {
		return fmt.Errorf("missing required dependency: contact service")
	}
	return nil
}

// Register registers the contact routes with the Echo instance.
// It validates dependencies before registering routes to fail fast
// if configuration is incomplete.
func (h *ContactHandler) Register(e *echo.Echo) {
	if err := h.Validate(); err != nil {
		h.LogError("failed to validate handler", err)
		return
	}

	g := e.Group("/api/v1/contact")
	g.POST("", h.handleSubmit)
	g.GET("", h.handleList)
	g.GET("/:id", h.handleGet)
	g.PUT("/:id", h.handleUpdate)
}

// handleSubmit handles contact form submissions
// @Summary Submit contact form
// @Description Submit a new contact form
// @Tags contact
// @Accept json
// @Produce json
// @Param submission body common.Submission true "Contact form submission"
// @Success 201 {object} common.Submission
// @Failure 400 {object} echo.HTTPError
// @Router /api/v1/contact [post]
func (h *ContactHandler) handleSubmit(c echo.Context) error {
	var submission common.Submission
	if err := c.Bind(&submission); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request format")
	}

	if err := c.Validate(submission); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := h.contactService.Submit(c.Request().Context(), &submission); err != nil {
		h.LogError("failed to submit contact form", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to submit contact form")
	}

	return c.JSON(http.StatusCreated, submission)
}

// handleList handles listing contact form submissions
// @Summary List contact form submissions
// @Description Get a list of all contact form submissions
// @Tags contact
// @Produce json
// @Success 200 {array} common.Submission
// @Failure 500 {object} echo.HTTPError
// @Router /api/v1/contact [get]
func (h *ContactHandler) handleList(c echo.Context) error {
	submissions, err := h.contactService.ListSubmissions(c.Request().Context())
	if err != nil {
		h.LogError("failed to list submissions", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to list submissions")
	}

	return c.JSON(http.StatusOK, submissions)
}

// handleGet handles getting a single contact form submission
// @Summary Get contact form submission
// @Description Get a specific contact form submission by ID
// @Tags contact
// @Produce json
// @Param id path int true "Submission ID"
// @Success 200 {object} common.Submission
// @Failure 400 {object} echo.HTTPError
// @Failure 404 {object} echo.HTTPError
// @Router /api/v1/contact/{id} [get]
func (h *ContactHandler) handleGet(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	submission, err := h.contactService.GetSubmission(c.Request().Context(), id)
	if err != nil {
		h.LogError("failed to get submission", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get submission")
	}

	if submission == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Submission not found")
	}

	return c.JSON(http.StatusOK, submission)
}

// handleUpdate handles updating a submission's status
// @Summary Update submission status
// @Description Update the status of a contact form submission
// @Tags contact
// @Accept json
// @Produce json
// @Param id path int true "Submission ID"
// @Param status body common.Status true "New status"
// @Success 200 {object} common.Submission
// @Failure 400 {object} echo.HTTPError
// @Failure 404 {object} echo.HTTPError
// @Router /api/v1/contact/{id} [put]
func (h *ContactHandler) handleUpdate(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	var status common.Status
	if err := c.Bind(&status); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid status format")
	}

	if err := h.contactService.UpdateSubmissionStatus(c.Request().Context(), id, string(status)); err != nil {
		h.LogError("failed to update submission status", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update status")
	}

	return c.NoContent(http.StatusOK)
}

// UpdateStatus handles updating the status of a contact submission
func (h *ContactHandler) UpdateStatus(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}
	status := c.QueryParam("status") // Assuming status is passed as a query parameter

	// Ensure status is of type string
	if err := h.contactService.UpdateSubmissionStatus(c.Request().Context(), id, status); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update status"})
	}

	return c.NoContent(http.StatusNoContent)
}
