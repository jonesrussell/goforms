package contact

import (
	"go.uber.org/fx"
)

// ProvideModule returns the fx options for the contact module.
func ProvideModule() fx.Option {
	return fx.Options(
		fx.Provide(
			NewRepository, // Provide the contact repository
			NewService,    // Provide the contact service
		),
		fx.Provide(
			fx.Annotate(
				NewService,
				fx.As(new(Service)), // Provide Service interface
			),
		),
	)
}
