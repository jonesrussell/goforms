package handlers

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jonesrussell/goforms/internal/logger"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTemplate_Render(t *testing.T) {
	tmpl := template.Must(template.New("base").Parse(`{{ define "base" }}<h1>{{.Title}}</h1>{{ end }}`))
	renderer := &Template{templates: tmpl}
	e := echo.New()
	rec := httptest.NewRecorder()
	c := e.NewContext(nil, rec)
	data := map[string]interface{}{
		"Title": "Test Title",
	}

	err := renderer.Render(rec, "base", data, c)
	assert.NoError(t, err)
	assert.Equal(t, "<h1>Test Title</h1>", strings.TrimSpace(rec.Body.String()))
}

func setupTestTemplates(t *testing.T) string {
	// Create a temporary directory for test templates
	tmpDir := t.TempDir()
	templatesDir := filepath.Join(tmpDir, "static", "templates")
	err := os.MkdirAll(templatesDir, 0755)
	require.NoError(t, err)

	// Create test templates
	templates := map[string]string{
		"layout.html":  `{{ define "base" }}<!DOCTYPE html><html><head><title>{{.Title}}</title></head><body>{{ template "content" . }}</body></html>{{ end }}`,
		"index.html":   `{{ define "content" }}<h1>Free Form Backend Service</h1>{{ end }}`,
		"contact.html": `{{ define "content" }}<h1>Contact Form Demo</h1>{{ end }}`,
	}

	for name, content := range templates {
		err := os.WriteFile(filepath.Join(templatesDir, name), []byte(content), 0644)
		require.NoError(t, err)
	}

	return tmpDir
}

func setupTestMarketingHandler(t *testing.T) (*MarketingHandler, *echo.Echo, func()) {
	logger := logger.GetLogger()

	// Create test templates inline
	templates := template.Must(template.New("base").Parse(`
		{{ define "base" }}<!DOCTYPE html><html><head><title>{{.Title}}</title></head><body>{{ template "content" . }}</body></html>{{ end }}
		{{ define "content" }}<h1>{{.Title}}</h1>{{ end }}
	`))

	handler := &MarketingHandler{
		logger:    logger,
		templates: templates,
	}
	e := echo.New()
	e.Renderer = &Template{templates: templates}
	return handler, e, func() {
		// Cleanup function can use t for test assertions if needed
		t.Cleanup(func() {})
	}
}

func TestNewMarketingHandler(t *testing.T) {
	// Set up test templates
	tmpDir := setupTestTemplates(t)

	// Change to temp dir for test
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}
	defer func() {
		if err := os.Chdir(oldWd); err != nil {
			t.Errorf("Failed to change back to original directory: %v", err)
		}
	}()
	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	logger := logger.GetLogger()
	handler := NewMarketingHandler(logger)
	assert.NotNil(t, handler)
	assert.NotNil(t, handler.templates)
	assert.NotNil(t, handler.logger)
}

func TestMarketingHandler_HomePage(t *testing.T) {
	handler, e, cleanup := setupTestMarketingHandler(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.HomePage(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "Form Backend")
}

func TestMarketingHandler_ContactPage(t *testing.T) {
	handler, e, cleanup := setupTestMarketingHandler(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/contact", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.ContactPage(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "Contact Form Demo")
}

func TestMarketingHandler_Register(t *testing.T) {
	handler, e, cleanup := setupTestMarketingHandler(t)
	defer cleanup()
	handler.Register(e)

	routes := e.Routes()
	var foundHome, foundContact bool
	for _, route := range routes {
		switch {
		case route.Path == "/" && route.Method == http.MethodGet:
			foundHome = true
		case route.Path == "/contact" && route.Method == http.MethodGet:
			foundContact = true
		}
	}

	require.True(t, foundHome, "Home route should be registered")
	require.True(t, foundContact, "Contact route should be registered")
}
