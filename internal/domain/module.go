package domain

import (
	"go.uber.org/fx"

	"github.com/jonesrussell/goforms/internal/domain/contact"
	"github.com/jonesrussell/goforms/internal/domain/user"
)

// Module combines all domain services
var Module = fx.Options(
	user.Module,
	contact.Module,
)
