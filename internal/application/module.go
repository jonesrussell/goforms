package application

import (
	"go.uber.org/fx"

	"github.com/jonesrussell/goforms/internal/application/handlers"
	"github.com/jonesrussell/goforms/internal/application/logging"
	"github.com/jonesrussell/goforms/internal/application/repositories/database"
	"github.com/jonesrussell/goforms/internal/domain/contact"
	"github.com/jonesrussell/goforms/internal/domain/user"
	"github.com/jonesrussell/goforms/internal/presentation/view"
)

// Module combines all application-level modules and providers.
var Module = fx.Options(
	fx.Provide(
		// View Renderer
		view.NewRenderer,

		NewWebHandler,
		NewAuthHandler,
		user.NewInMemoryTokenRepository,
		NewUserService,
	),
	fx.Provide(func(db *database.DB) user.TokenRepository {
		return user.NewTokenRepository(db)
	}),
)

func NewWebHandler(logger logging.Logger, renderer *view.Renderer, contactService contact.Service) *handlers.WebHandler {
	return handlers.NewWebHandler(logger, handlers.WithRenderer(renderer), handlers.WithContactService(contactService))
}

func NewAuthHandler(logger logging.Logger, userService *user.Service) *handlers.AuthHandler {
	return handlers.NewAuthHandler(logger, userService)
}

func NewUserService(repo user.Repository, tokenRepo user.TokenRepository, logger logging.Logger) *user.Service {
	return user.NewService(repo, tokenRepo, logger)
}

// AsHandler annotates the given constructor to state that
// it provides a handler to the "handlers" group.
func AsHandler(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(handlers.Handler)),
		fx.ResultTags(`group:"handlers"`),
	)
}
