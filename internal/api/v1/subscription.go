package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/jonesrussell/goforms/internal/core/subscription"
	"github.com/jonesrussell/goforms/internal/logger"
	"github.com/jonesrussell/goforms/internal/response"
)

// SubscriptionAPI handles subscription-related API endpoints
type SubscriptionAPI struct {
	service *subscription.Service
	logger  logger.Logger
}

// NewSubscriptionAPI creates a new subscription API handler
func NewSubscriptionAPI(service *subscription.Service, log logger.Logger) *SubscriptionAPI {
	return &SubscriptionAPI{
		service: service,
		logger:  log,
	}
}

// Register registers the subscription API routes
func (api *SubscriptionAPI) Register(e *echo.Echo) {
	v1 := e.Group("/api/v1")
	subscriptions := v1.Group("/subscriptions")

	subscriptions.POST("", api.CreateSubscription)
	subscriptions.GET("", api.ListSubscriptions)
	subscriptions.GET("/:id", api.GetSubscription)
	subscriptions.PUT("/:id/status", api.UpdateSubscriptionStatus)
	subscriptions.DELETE("/:id", api.DeleteSubscription)
}

// CreateSubscription handles subscription creation
func (api *SubscriptionAPI) CreateSubscription(c echo.Context) error {
	var sub subscription.Subscription
	if err := c.Bind(&sub); err != nil {
		api.logger.Error("failed to bind subscription", logger.Error(err))
		return response.Error(c, http.StatusBadRequest, "invalid request")
	}

	if err := api.service.CreateSubscription(c.Request().Context(), &sub); err != nil {
		api.logger.Error("failed to create subscription", logger.Error(err))
		return response.Error(c, http.StatusInternalServerError, "failed to create subscription")
	}

	return response.Success(c, http.StatusCreated, sub)
}

// ListSubscriptions handles listing subscriptions
func (api *SubscriptionAPI) ListSubscriptions(c echo.Context) error {
	subs, err := api.service.ListSubscriptions(c.Request().Context())
	if err != nil {
		api.logger.Error("failed to list subscriptions", logger.Error(err))
		return response.Error(c, http.StatusInternalServerError, "failed to list subscriptions")
	}

	return response.Success(c, http.StatusOK, subs)
}

// GetSubscription handles retrieving a single subscription
func (api *SubscriptionAPI) GetSubscription(c echo.Context) error {
	id, err := response.ParseInt64Param(c, "id")
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid subscription id")
	}

	sub, err := api.service.GetSubscription(c.Request().Context(), id)
	if err != nil {
		api.logger.Error("failed to get subscription", logger.Error(err))
		return response.Error(c, http.StatusInternalServerError, "failed to get subscription")
	}

	if sub == nil {
		return response.Error(c, http.StatusNotFound, "subscription not found")
	}

	return response.Success(c, http.StatusOK, sub)
}

// UpdateSubscriptionStatus handles updating a subscription's status
func (api *SubscriptionAPI) UpdateSubscriptionStatus(c echo.Context) error {
	id, err := response.ParseInt64Param(c, "id")
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid subscription id")
	}

	var req struct {
		Status subscription.Status `json:"status"`
	}
	if err := c.Bind(&req); err != nil {
		api.logger.Error("failed to bind status update request", logger.Error(err))
		return response.Error(c, http.StatusBadRequest, "invalid request")
	}

	if err := api.service.UpdateSubscriptionStatus(c.Request().Context(), id, req.Status); err != nil {
		api.logger.Error("failed to update subscription status", logger.Error(err))
		return response.Error(c, http.StatusInternalServerError, "failed to update subscription status")
	}

	return response.Success(c, http.StatusOK, nil)
}

// DeleteSubscription handles subscription deletion
func (api *SubscriptionAPI) DeleteSubscription(c echo.Context) error {
	id, err := response.ParseInt64Param(c, "id")
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid subscription id")
	}

	if err := api.service.DeleteSubscription(c.Request().Context(), id); err != nil {
		api.logger.Error("failed to delete subscription", logger.Error(err))
		return response.Error(c, http.StatusInternalServerError, "failed to delete subscription")
	}

	return response.Success(c, http.StatusOK, nil)
}
