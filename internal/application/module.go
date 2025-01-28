package application

import (
	"go.uber.org/fx"

	"github.com/jonesrussell/goforms/internal/application/handlers"
	"github.com/jonesrussell/goforms/internal/application/logging"
	"github.com/jonesrussell/goforms/internal/domain/contact"
	"github.com/jonesrussell/goforms/internal/domain/user"
	"github.com/jonesrussell/goforms/internal/presentation/view"
)

// Module combines all application-level modules and providers.
var Module = fx.Options(
	// Provide your application-level dependencies here
	fx.Provide(
		// View Renderer
		view.NewRenderer,

		// Handlers
		AsHandler(func(logger logging.Logger, renderer *view.Renderer, contactService contact.Service) *handlers.WebHandler {
			return handlers.NewWebHandler(logger,
				handlers.WithRenderer(renderer),
				handlers.WithContactService(contactService),
			)
		}),
		// Add other handlers as needed
		AsHandler(func(logger logging.Logger, userService user.Service) *handlers.AuthHandler {
			return handlers.NewAuthHandler(logger, userService)
		}),
	),
)

// AsHandler annotates the given constructor to state that
// it provides a handler to the "handlers" group.
func AsHandler(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(handlers.Handler)),
		fx.ResultTags(`group:"handlers"`),
	)
}
