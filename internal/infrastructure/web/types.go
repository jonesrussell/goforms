// Package web provides utilities for handling web assets in the application.
// It supports both development mode (using Vite dev server) and production mode
// (using built assets from the Vite manifest).
package web

import (
	"context"
	"errors"

	"github.com/labstack/echo/v4"
)

// AssetType represents the type of asset
type AssetType string

const (
	// AssetTypeJS represents JavaScript files
	AssetTypeJS AssetType = "js"
	// AssetTypeCSS represents CSS files
	AssetTypeCSS AssetType = "css"
	// AssetTypeImage represents image files
	AssetTypeImage AssetType = "image"
	// AssetTypeFont represents font files
	AssetTypeFont AssetType = "font"
	MaxPathLength           = 100
)

// Asset-related errors
var (
	ErrAssetNotFound    = errors.New("asset not found")
	ErrInvalidManifest  = errors.New("invalid manifest")
	ErrInvalidPath      = errors.New("invalid asset path")
	ErrManifestNotFound = errors.New("manifest not found")
)

// ManifestEntry represents an entry in the Vite manifest file
type ManifestEntry struct {
	File    string   `json:"file"`
	Name    string   `json:"name"`
	Src     string   `json:"src"`
	IsEntry bool     `json:"isEntry"`
	CSS     []string `json:"css"`
}

// Manifest represents the Vite manifest file
type Manifest map[string]ManifestEntry

// AssetResolver interface separates resolution logic from management
type AssetResolver interface {
	ResolveAssetPath(ctx context.Context, path string) (string, error)
}

// AssetServer defines the interface for serving assets
type AssetServer interface {
	// RegisterRoutes registers the necessary routes for serving assets
	RegisterRoutes(e *echo.Echo) error
}

// WebModule encapsulates the asset manager and server to eliminate global state
type WebModule struct {
	AssetManager *AssetManager
	AssetServer  AssetServer
}
