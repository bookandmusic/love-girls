package infra

import (
	"github.com/google/wire"
)

// InfraSet 包含所有基础设施相关的provider
var InfraSet = wire.NewSet(
	ProvideConfig,
	ProvideLogger,
	ProvideGormLogger,
	ProvideJWT,
	ProvideAuthMiddleware,
	ProvideDB,
	ProvideMigrate,
)
