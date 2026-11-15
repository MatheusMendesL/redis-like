package usecase

import (
	"context"
	"testing"

	"redis-project/internal/adapter/protocol"
	"redis-project/internal/domain/repository"
	"redis-project/internal/infrastructure/storage"
)

type recordingPersistence struct {
	command []string
	args    []string
}

func (p *recordingPersistence) Append(_ context.Context, command, args []string) error {
	p.command = append([]string(nil), command...)
	p.args = append([]string(nil), args...)
	return nil
}

func (p *recordingPersistence) Replay(context.Context, repository.KeyValueRepository) error {
	return nil
}

func (p *recordingPersistence) Close() error {
	return nil
}

func TestExpirePersistsKeyAndDuration(t *testing.T) {
	store := storage.NewStore()
	persist := &recordingPersistence{}
	handler := NewCommandHandler(store, persist, protocol.NewParser(), NewStats(store))
	ctx := context.Background()

	if response := handler.ExecuteCommand(ctx, &protocol.Command{
		Type: "SET",
		Args: []string{"key", "value"},
	}); response != "OK" {
		t.Fatalf("SET response = %q, want OK", response)
	}

	if response := handler.ExecuteCommand(ctx, &protocol.Command{
		Type: "EXPIRE",
		Args: []string{"key", "60"},
	}); response != "OK" {
		t.Fatalf("EXPIRE response = %q, want OK", response)
	}

	if len(persist.command) != 1 || persist.command[0] != "EXPIRE" {
		t.Fatalf("persisted command = %v, want [EXPIRE]", persist.command)
	}
	if len(persist.args) != 2 || persist.args[0] != "key" || persist.args[1] != "60" {
		t.Fatalf("persisted args = %v, want [key 60]", persist.args)
	}
}
