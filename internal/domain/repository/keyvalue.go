package repository

import "context"

type KeyValueRepository interface {
	Set(ctx context.Context, key string, value string) error
	Get(ctx context.Context, key string) (string, bool, error)
	Delete(ctx context.Context, key string) int
	Expire(ctx context.Context, key string, seconds int) (bool, error)
	TTL(ctx context.Context, key string) (int64, error)
	Persist(ctx context.Context, key string) (bool, error)
	Keys(ctx context.Context, pattern string) ([]string, error)
	Exists(ctx context.Context, key string) (bool, error)
	Size(ctx context.Context) (int, error)

	StartCleanup(intervalMs int64) error
	StopCleanup() error
}
