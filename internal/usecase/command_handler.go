package usecase

import (
	"context"
	"redis-project/internal/adapter/protocol"
	"redis-project/internal/domain/command"
	"redis-project/internal/domain/repository"
	"strconv"
	"strings"
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

	types := command.All(false)
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
	if len(args) != 2 {
		return h.parser.FormatError("SET requires at least 2 arguments")
	}

	key := args[0]
	value := strings.Join(args[1:], " ")
	if err := h.store.Set(ctx, key, value); err != nil {
		return h.parser.FormatError(err.Error())
	}

	if h.persist != nil {
		if err := h.persist.Append(ctx, []string{command.SET.String()}, args); err != nil {
			return h.parser.FormatError(err.Error())
		}
	}
	return h.parser.FormatOK()
}

func (h *CommandHandler) handleGet(ctx context.Context, args []string) string {
	if len(args) != 1 {
		return h.parser.FormatError("GET requires at least 1 arguments")
	}

	value, found, err := h.store.Get(ctx, args[0])
	if err != nil {
		return h.parser.FormatError(err.Error())
	}

	if found {
		return value
	}
	return h.parser.FormatOK()
}

func (h *CommandHandler) handleDel(ctx context.Context, args []string) string {
	if len(args) != 1 {
		return h.parser.FormatError("DEL requires at least 1 arguments")
	}

	count := 0
	for _, keys := range args {
		count += h.store.Delete(ctx, keys)
	}

	if h.persist != nil && count > 0 {
		for _, key := range args {
			if err := h.persist.Append(ctx, []string{command.DEL.String()}, []string{key}); err != nil {
				return h.parser.FormatError(err.Error())
			}
		}
	}

	return h.parser.FormatResponse(count)
}

func (h *CommandHandler) handleExists(ctx context.Context, args []string) string {
	if len(args) < 1 {
		return h.parser.FormatError("EXISTS requires at least 1 arguments")
	}

	count := 0
	for _, keys := range args {
		exists, err := h.store.Exists(ctx, keys)
		if err != nil {
			return h.parser.FormatError(err.Error())
		}
		if exists {
			count++
		}
	}

	return h.parser.FormatResponse(count)
}

func (h *CommandHandler) handleExpire(ctx context.Context, args []string) string {
	if len(args) != 2 {
		return h.parser.FormatError("EXPIRE requires 2 arguments")
	}

	key := args[0]
	seconds, err := strconv.Atoi(args[1])

	if err != nil {
		return h.parser.FormatError("Invalid seconds value")
	}

	success, err := h.store.Expire(ctx, key, seconds)

	if err != nil {
		return h.parser.FormatError(err.Error())
	}

	if success {
		if h.persist != nil {
			if err := h.persist.Append(ctx, []string{command.EXPIRE.String()}, []string{key}); err != nil {
				return h.parser.FormatError(err.Error())
			}
		}
		return h.parser.FormatOK()
	}

	return h.parser.FormatResponse(success)
}

func (h *CommandHandler) handlePersist(ctx context.Context, args []string) string {
	if len(args) != 1 {
		return h.parser.FormatError("PERSIST requires at least 1 arguments")
	}

	success, err := h.store.Persist(ctx, args[0])

	if err != nil {
		return h.parser.FormatError(err.Error())
	}

	if success {
		if h.persist != nil {
			if err := h.persist.Append(ctx, []string{command.PERSIST.String()}, []string{args[0]}); err != nil {
				return h.parser.FormatError(err.Error())
			}
		}
		return h.parser.FormatOK()
	}
	return h.parser.FormatResponse(success)
}

func (h *CommandHandler) handleTTL(ctx context.Context, args []string) string {
	if len(args) != 1 {
		return h.parser.FormatError("TTL requires at least 1 arguments")
	}

	ttl, err := h.store.TTL(ctx, args[0])

	if err != nil {
		return h.parser.FormatError(err.Error())
	}

	return h.parser.FormatResponse(ttl)
}

func (h *CommandHandler) handleKeys(ctx context.Context, args []string) string {
	pattern := "*"

	if len(args) > 0 {
		pattern = args[0]
	}

	keys, err := h.store.Keys(ctx, pattern)

	if err != nil {
		return h.parser.FormatError(err.Error())
	}

	if len(keys) == 0 {
		return ""
	}

	return strings.Join(keys, " ")
}

func (h *CommandHandler) handlePing(ctx context.Context, args []string) string {
	message := "PING"

	if len(args) > 0 {
		message = strings.Join(args, " ")
	}

	return message
}

func (h *CommandHandler) handleInfo(ctx context.Context, args []string) string {
	stats, err := h.stats.GetInfo(ctx)

	if err != nil {
		return h.parser.FormatError(err.Error())
	}

	return stats
}
