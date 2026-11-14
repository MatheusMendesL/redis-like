package container

import (
	"redis-project/internal/adapter/handler"
	"redis-project/internal/adapter/protocol"
	"redis-project/internal/domain/repository"
	"redis-project/internal/usecase"
)

type Container struct {
	Store          repository.KeyValueRepository
	Persist        repository.PersistenceRepository
	CommandHandler *usecase.CommandHandler
	TCPHandler     *handler.TCPHandler
	Parser         *protocol.Parser
}

func NewContainer(
	store repository.KeyValueRepository,
	persist repository.PersistenceRepository,
	commandHandler *usecase.CommandHandler,
	tcphandler *handler.TCPHandler,
	parser *protocol.Parser,
) *Container {
	return &Container{
		Store:          store,
		Persist:        persist,
		CommandHandler: commandHandler,
		TCPHandler:     tcphandler,
		Parser:         parser,
	}
}

func (c *Container) Close() error {
	if c.Persist != nil {
		return c.Persist.Close()
	}

	return nil
}
