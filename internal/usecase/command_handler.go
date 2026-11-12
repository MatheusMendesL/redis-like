package usecase

import (
	"context"
	"redis-project/internal/adapter/protocol"
	"redis-project/internal/domain/command"
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

func (h *CommandHandler) AllFuncs() []func(context.Context, []string) string {
	return []func(context.Context, []string) string{
		h.handleSet,
		h.handleGet,
		h.handleDel,
		h.handleExpire,
		h.handleTTL,
		h.handlePersist,
		nil, 
		h.handleKeys,
		h.handleExists,
		h.handlePing,
		h.handleInfo,
	}
}

func (h *CommandHandler) ExecuteCommand(ctx context.Context, cmd *protocol.Command) string {
	h.stats.IncrementCommands()

	types := command.All()
	funcs := h.AllFuncs()

	for i, t := range types {
		if cmd.Type == t {
			if funcs[i] != nil {
				return funcs[i](ctx, cmd.Args)
			}
			break
		}
	}

	return "ERR: this method doesn´t exist"
}

func (h *CommandHandler) handleSet(ctx context.Context, args []string) string {
	return ""
}

func (h *CommandHandler) handleGet(ctx context.Context, args []string) string {
	return ""
}

func (h *CommandHandler) handleDel(ctx context.Context, args []string) string {
	return ""
}

func (h *CommandHandler) handleExists(ctx context.Context, args []string) string {
	return ""
}

func (h *CommandHandler) handleExpire(ctx context.Context, args []string) string {
	return ""
}

func (h *CommandHandler) handleExpireAt(ctx context.Context, args []string) string {
	return ""
}

func (h *CommandHandler) handlePersist(ctx context.Context, args []string) string {
	return ""
}

func (h *CommandHandler) handleTTL(ctx context.Context, args []string) string {
	return ""
}

func (h *CommandHandler) handleKeys(ctx context.Context, args []string) string {
	return ""
}

func (h *CommandHandler) handlePing(ctx context.Context, args []string) string {
	return ""
}

func (h *CommandHandler) handleInfo(ctx context.Context, args []string) string {
	return ""
}
