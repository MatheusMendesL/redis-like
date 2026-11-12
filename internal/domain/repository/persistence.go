package repository

import "context"

type PersistenceRepository interface {
	Append(ctx context.Context, command, args []string) error
	Replay(ctx context.Context, store KeyValueRepository) error
	Close() error
}
