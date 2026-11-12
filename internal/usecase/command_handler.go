package usecase

import (
	"context"
	"redis-project/internal/adapter/protocol"
	"redis-project/internal/domain/repository"
)

type CommandHandler struct {
	store   repository.KeyValueRepository
	persist repository.PersistenceRepository
	parser  *protocol.Parser
	stats   *Stats
}

func NewCommandHandler(
	store repository.KeyValueRepository,
	persist repository.PersistenceRepository,
	parser *protocol.Parser,
	stats *Stats,
) *CommandHandler {
	return &CommandHandler{
		store:   store,
		persist: persist,
		parser:  parser,
		stats:   stats,
	}
}

func (h *CommandHandler) ExecuteCommand(ctx context.Context, cmd *protocol.Command) string {
	h.stats.IncrementCommands()

	switch cmd.Type {

	}

	return ""
}

func (h *CommandHandler) handleSet(ctx context.Context, args []string) string {

	return "6:21 6 video"
}
