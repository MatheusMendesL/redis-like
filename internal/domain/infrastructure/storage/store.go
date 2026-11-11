package storage

import (
	"context"
	"redis-project/internal/domain/entity"
	"redis-project/internal/domain/repository"
	"sync"
	"time"
)

type Store struct {
	data        map[string]*entity.Item
	mu          sync.RWMutex
	stopCleanUp chan struct{}
}

func NewStore() repository.KeyValueRepository {
	return &Store{
		data:        make(map[string]*entity.Item),
		stopCleanUp: make(chan struct{}),
	}
}

func (s *Store) Set(ctx context.Context, key, value string) {
	if ctx.Err() != nil {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[key] = &entity.Item{
		Value:     value,
		ExpiresAt: nil,
	}
}

func (s *Store) Get(ctx context.Context, key string) (string, bool) {
	if ctx.Err() != nil {
		return "", false
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	item, exists := s.data[key]
	if !exists || item.IsExpired(time.Now().Unix()) {
		return "", false
	}

	return item.Value, true
}

func (s *Store) Delete(ctx context.Context, key string) int {
	if ctx.Err() != nil {
		return 0
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.data[key]; exists {
		delete(s.data, key)
		return 1
	}

	return 0

}

func (s *Store) Expire(ctx context.Context, key string, seconds int) bool {
	if ctx.Err() != nil {
		return false
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	item, exists := s.data[key]

	if !exists {
		return false
	}

	if seconds <= 0 {
		return false
	}

	expiresAt := time.Now().Unix() + int64(seconds)
	item.ExpiresAt = &expiresAt
	return true

}

func (s *Store) TTL(ctx context.Context, key string) int64 {
	if ctx.Err() != nil {
		return 0
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	item, exists := s.data[key]

	if !exists {
		return -2
	}

	if item.IsExpired(time.Now().Unix()) {
		return -2
	}

	if item.ExpiresAt == nil {
		return -1
	}

	now := time.Now().Unix()
	remaining := *item.ExpiresAt - now

	if remaining <= 0 {
		return -2
	}

	return remaining
}

func (s *Store) Exists(ctx context.Context, key string) bool {
	panic("unimplemented")
}

func (s *Store) Keys(ctx context.Context, pattern string) []string {
	panic("unimplemented")
}

func (s *Store) Persist(ctx context.Context, key string) bool {
	panic("unimplemented")
}

func (s *Store) Size(ctx context.Context) int {
	panic("unimplemented")
}

func (s *Store) StartCleanup(intervalMs int64) {
	panic("unimplemented")
}

func (s *Store) StopCleanup() {
	panic("unimplemented")
}
