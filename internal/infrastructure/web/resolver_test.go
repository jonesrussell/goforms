package web_test

import (
	"embed"
	"testing"

	"github.com/goformx/goforms/internal/infrastructure/config"
	"github.com/goformx/goforms/internal/infrastructure/web"
	mocklogging "github.com/goformx/goforms/test/mocks/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestDevelopmentAssetResolver(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		App: config.AppConfig{
			Scheme:      "http",
			ViteDevHost: "localhost",
			ViteDevPort: "3000",
		},
	}
	mockLogger := mocklogging.NewMockLogger(ctrl)

	mockLogger.EXPECT().Debug(
		"resolving development asset path", "path", gomock.Any(), "host_port", gomock.Any(),
	).AnyTimes()
	mockLogger.EXPECT().Debug(
		"development asset resolved", "input", gomock.Any(), "output", gomock.Any(),
	).AnyTimes()

	resolver := web.NewDevelopmentAssetResolver(cfg, mockLogger)

	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "src file",
			path:     "src/js/pages/main.ts",
			expected: "http://localhost:3000/src/js/pages/main.ts",
		},
		{
			name:     "css file",
			path:     "main.css",
			expected: "http://localhost:3000/src/css/main.css",
		},
		{
			name:     "js file",
			path:     "main.js",
			expected: "http://localhost:3000/src/js/pages/main.ts",
		},
		{
			name:     "ts file",
			path:     "main.ts",
			expected: "http://localhost:3000/src/js/pages/main.ts",
		},
		{
			name:     "other file",
			path:     "image.png",
			expected: "http://localhost:3000/image.png",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := resolver.ResolveAssetPath(t.Context(), tt.path)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestProductionAssetResolver(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	manifest := web.Manifest{
		"main.js": {
			File: "assets/main.abc123.js",
		},
		"style.css": {
			File: "assets/style.def456.css",
		},
	}
	mockLogger := mocklogging.NewMockLogger(ctrl)

	mockLogger.EXPECT().Debug(
		"resolving production asset path", "path", gomock.Any(), "manifest_entries", gomock.Any(),
	).AnyTimes()
	mockLogger.EXPECT().Debug(
		"production asset resolved", "input", gomock.Any(), "output", gomock.Any(),
	).AnyTimes()
	mockLogger.EXPECT().Debug(
		"asset not found in manifest", "path", gomock.Any(), "available_keys", gomock.Any(),
	).AnyTimes()

	resolver := web.NewProductionAssetResolver(manifest, mockLogger)

	tests := []struct {
		name        string
		path        string
		expected    string
		expectError bool
	}{
		{
			name:     "existing js file",
			path:     "main.js",
			expected: "/assets/main.abc123.js",
		},
		{
			name:     "existing css file",
			path:     "style.css",
			expected: "/assets/style.def456.css",
		},
		{
			name:        "non-existent file",
			path:        "missing.js",
			expectError: true,
		},
		{
			name:        "empty path",
			path:        "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := resolver.ResolveAssetPath(t.Context(), tt.path)
			if tt.expectError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAssetManager(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		App: config.AppConfig{
			Env:         "development",
			Scheme:      "http",
			ViteDevHost: "localhost",
			ViteDevPort: "3000",
		},
	}
	mockLogger := mocklogging.NewMockLogger(ctrl)

	mockLogger.EXPECT().Info(
		"asset manager initialized in development mode",
	).Times(1)
	mockLogger.EXPECT().Debug(
		"asset manager resolving path", "input_path", gomock.Any(),
	).AnyTimes()
	mockLogger.EXPECT().Debug(
		"asset path found in cache", "input_path", gomock.Any(), "cached_path", gomock.Any(),
	).AnyTimes()
	mockLogger.EXPECT().Debug(
		"asset path resolved", "input_path", gomock.Any(), "resolved_path", gomock.Any(), "environment", gomock.Any(),
	).AnyTimes()
	mockLogger.EXPECT().Debug(
		"asset path cache cleared",
	).AnyTimes()
	mockLogger.EXPECT().Debug(
		"resolving development asset path", "path", gomock.Any(), "host_port", gomock.Any(),
	).AnyTimes()
	mockLogger.EXPECT().Debug(
		"development asset resolved", "input", gomock.Any(), "output", gomock.Any(),
	).AnyTimes()

	var distFS embed.FS

	manager, err := web.NewAssetManager(cfg, mockLogger, distFS)
	require.NoError(t, err)
	assert.NotNil(t, manager)

	path := manager.AssetPath("main.js")
	assert.Equal(t, "http://localhost:3000/src/js/pages/main.ts", path)

	path2 := manager.AssetPath("main.js")
	assert.Equal(t, path, path2)

	assetType := manager.GetAssetType("main.js")
	assert.Equal(t, web.AssetTypeJS, assetType)

	assetType = manager.GetAssetType("style.css")
	assert.Equal(t, web.AssetTypeCSS, assetType)

	manager.ClearCache()
}

func TestAssetManager_ProductionMode(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		App: config.AppConfig{
			Env: "production",
		},
	}
	mockLogger := mocklogging.NewMockLogger(ctrl)

	mockLogger.EXPECT().Debug(
		"loading manifest from embedded filesystem", "path", "dist/.vite/manifest.json",
	).AnyTimes()

	var distFS embed.FS

	manager, err := web.NewAssetManager(cfg, mockLogger, distFS)
	require.Error(t, err)
	require.Nil(t, manager)
}

func TestAssetManager_ErrorHandling(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocklogging.NewMockLogger(ctrl)

	manager, err := web.NewAssetManager(nil, mockLogger, embed.FS{})
	require.Error(t, err)
	require.Nil(t, manager)
	require.Contains(t, err.Error(), "config is required")

	cfg := &config.Config{
		App: config.AppConfig{
			Env: "development",
		},
	}
	manager, err = web.NewAssetManager(cfg, nil, embed.FS{})
	require.Error(t, err)
	require.Nil(t, manager)
	require.Contains(t, err.Error(), "logger is required")
}
