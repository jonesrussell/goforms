package user

import (
	"go.uber.org/fx"
)

// ProvideModule returns the fx options for the user module.
func ProvideModule() fx.Option {
	return fx.Options(
		fx.Provide(
			NewRepository,      // Provide the user Repository
			NewTokenRepository, // Provide the TokenRepository
			NewService,         // Provide the user Service
		),
		fx.Provide(
			fx.Annotate(
				NewRepository,
				fx.As(new(Repository)),
				fx.As(new(TokenRepository)), // Provide TokenRepository interface
			),
		),
	)
}
