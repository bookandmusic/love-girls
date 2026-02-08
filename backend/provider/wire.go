//go:build wireinject
// +build wireinject

package provider

import (
	"github.com/google/wire"

	"github.com/bookandmusic/love-girl/internal/server"
	"github.com/bookandmusic/love-girl/provider/infra"
)

func InitApp() (*server.App, func(), error) {
	wire.Build(
		// migrate
		infra.InfraSet,

		// storage
		ProvideStorage,

		// repo
		RepoSet,
		// service
		ServiceSet,
		// handler
		HandlerSet,
		// router (includes GinEngine, Router, and App)
		RouterSet,
	)
	return nil, nil, nil
}
