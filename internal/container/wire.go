//go:build wireinject
// +build wireinject

package container

import (
	"redis-project/internal/adapter/handler"
	"redis-project/internal/adapter/protocol"
	"redis-project/internal/infrastructure/persistence"
	"redis-project/internal/infrastructure/storage"
	"redis-project/internal/usecase"

	"github.com/google/wire"
)

func InitializeContainer(opt persistence.AOFProviderOption) (*Container, func(), error) {
	wire.Build(

		// infrastructure
		storage.NewStore,
		persistence.NewAOFProvider,

		// adapter
		protocol.NewParser,

		// usecase
		usecase.NewStats,
		usecase.NewCommandHandler,

		// handler
		handler.NewTCPHandler,

		// container
		NewContainer,
	)

	return nil, nil, nil
}
