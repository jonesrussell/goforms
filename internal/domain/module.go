package domain

import (
	"go.uber.org/fx"

	"github.com/jonesrussell/goforms/internal/domain/contact"
	"github.com/jonesrussell/goforms/internal/domain/user"
)

// ProvideModule aggregates all domain modules.
func ProvideModule() fx.Option {
	return fx.Options(
		user.ProvideModule(),    // Provide the user module
		contact.ProvideModule(), // Provide the contact module
	)
}
