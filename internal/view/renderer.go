package view

import (
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"

	"github.com/jonesrussell/goforms/internal/logger"
)

type Renderer struct {
	logger logger.Logger
}

func NewRenderer(logger logger.Logger) *Renderer {
	return &Renderer{
		logger: logger,
	}
}

func (r *Renderer) Render(c echo.Context, t templ.Component) error {
	if err := t.Render(c.Request().Context(), c.Response().Writer); err != nil {
		r.logger.Error("failed to render template",
			logger.Error(err),
			logger.String("template", fmt.Sprintf("%T", t)),
		)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to render page")
	}
	return nil
}
