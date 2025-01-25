package infrastructure

import (
	"context"

	"go.uber.org/fx"

	"github.com/jonesrussell/goforms/internal/application/handler"
	v1 "github.com/jonesrussell/goforms/internal/application/http/v1"
	"github.com/jonesrussell/goforms/internal/domain/contact"
	"github.com/jonesrussell/goforms/internal/domain/subscription"
	"github.com/jonesrussell/goforms/internal/domain/user"
	"github.com/jonesrussell/goforms/internal/infrastructure/config"
	"github.com/jonesrussell/goforms/internal/infrastructure/database"
	"github.com/jonesrussell/goforms/internal/infrastructure/logging"
	"github.com/jonesrussell/goforms/internal/infrastructure/server"
	"github.com/jonesrussell/goforms/internal/infrastructure/store"
	"github.com/jonesrussell/goforms/internal/presentation/view"
)

// Module combines all infrastructure-level modules and providers
//
//nolint:gochecknoglobals // fx modules are designed to be global
var Module = fx.Options(
	// Core infrastructure
	fx.Provide(
		config.New,
		logging.NewLogger,
		store.NewStore,
	),

	// Stores
	fx.Provide(
		store.NewUserStore,
		store.NewContactStore,
		store.NewSubscriptionStore,
	),

	// Services
	fx.Provide(
		user.NewService,
		contact.NewService,
		subscription.NewService,
	),

	// Handlers
	fx.Provide(
		fx.Annotated{
			Group:  "handlers",
			Target: v1.NewHandler,
		},
		fx.Annotated{
			Group:  "handlers",
			Target: v1.NewAuthHandler,
		},
		fx.Annotated{
			Group:  "handlers",
			Target: v1.NewContactAPI,
		},
		fx.Annotated{
			Group:  "handlers",
			Target: v1.NewSubscriptionAPI,
		},
	),

	// Renderer
	fx.Provide(
		view.NewRenderer,
	),

	// Lifecycle hooks
	fx.Invoke(
		registerDatabaseHooks,
	),
)

// HandlerParams contains dependencies for creating handlers
type HandlerParams struct {
	fx.In

	Logger      logging.Logger
	Config      *config.Config
	Server      *server.Server
	VersionInfo handler.VersionInfo `name:"version_info"`
	Renderer    *view.Renderer

	ContactService      contact.Service
	SubscriptionService subscription.Service
	UserService         user.Service
}

// HandlerResult contains all HTTP handlers
type HandlerResult struct {
	fx.Out

	Handlers []handler.Handler `group:"handlers"`
}

// NewHandlers creates all HTTP handlers
func NewHandlers(p HandlerParams) HandlerResult {
	base := handler.Base{Logger: p.Logger}

	return HandlerResult{
		Handlers: []handler.Handler{
			handler.NewVersionHandler(p.VersionInfo, base),
			handler.NewWebHandler(p.Logger, p.ContactService, p.Renderer),
			handler.NewAuthHandler(p.Logger, p.UserService),
			handler.NewContactHandler(p.Logger, p.ContactService),
			handler.NewSubscriptionHandler(p.Logger, p.SubscriptionService),
		},
	}
}

// Stores groups all database store providers
type Stores struct {
	fx.Out

	ContactStore      contact.Store
	SubscriptionStore subscription.Store
	UserStore         user.Store
}

// NewStores creates all database stores
func NewStores(db *database.Database, logger logging.Logger) Stores {
	return Stores{
		ContactStore:      store.NewContactStore(db, logger),
		SubscriptionStore: store.NewSubscriptionStore(db, logger),
		UserStore:         store.NewUserStore(db, logger),
	}
}

func registerDatabaseHooks(lc fx.Lifecycle, db *database.Database, logger logging.Logger) {
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			logger.Info("closing database connection")
			return db.Close()
		},
	})
}
