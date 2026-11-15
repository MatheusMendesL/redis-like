package persistence

import (
	"context"
	"path/filepath"
	"testing"

	"redis-project/internal/infrastructure/storage"
)

func TestReplayRestoresExpireDuration(t *testing.T) {
	path := filepath.Join(t.TempDir(), "data.aof")
	aof, err := NewAOF(path)
	if err != nil {
		t.Fatal(err)
	}
	defer aof.Close()

	ctx := context.Background()
	if err := aof.Append(ctx, []string{"SET"}, []string{"key", "value"}); err != nil {
		t.Fatal(err)
	}
	if err := aof.Append(ctx, []string{"EXPIRE"}, []string{"key", "60"}); err != nil {
		t.Fatal(err)
	}

	store := storage.NewStore()
	if err := aof.Replay(ctx, store); err != nil {
		t.Fatal(err)
	}

	ttl, err := store.TTL(ctx, "key")
	if err != nil {
		t.Fatal(err)
	}
	if ttl < 59 || ttl > 60 {
		t.Fatalf("expected TTL close to 60 seconds, got %d", ttl)
	}
}
