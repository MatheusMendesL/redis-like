package storage

import (
	"context"
	"path/filepath"
	"redis-project/internal/domain/entity"
	domainerrors "redis-project/internal/domain/errors"
	"redis-project/internal/domain/repository"
	"sync"
	"time"
)

var _ repository.KeyValueRepository = (*Store)(nil)

type Store struct {
	data           map[string]*entity.Item
	mu             sync.RWMutex
	stopCleanUp    chan struct{}
	cleanupOnce    sync.Once
	cleanupStarted bool
}

func NewStore() repository.KeyValueRepository {
	return &Store{
		data:        make(map[string]*entity.Item),
		stopCleanUp: make(chan struct{}),
	}
}

func (s *Store) Set(ctx context.Context, key, value string) error {
	if err := ctx.Err(); err != nil {
		return domainerrors.ErrContextCanceled
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[key] = &entity.Item{
		Value:     value,
		ExpiresAt: nil,
	}

	return nil
}

func (s *Store) Get(ctx context.Context, key string) (string, bool, error) {
	if err := ctx.Err(); err != nil {
		return "", false, domainerrors.ErrContextCanceled
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	item, exists := s.data[key]

	if !exists || item.IsExpired(time.Now().Unix()) {
		return "", false, domainerrors.ErrKeyNotFound
	}

	return item.Value, true, nil
}

func (s *Store) Delete(ctx context.Context, key string) int {
	if err := ctx.Err(); err != nil {
		return 0
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	item, exists := s.data[key]

	if !exists || item.IsExpired(time.Now().Unix()) {
		return 0
	}

	delete(s.data, key)
	return 1
}

func (s *Store) Expire(ctx context.Context, key string, seconds int) (bool, error) {
	if err := ctx.Err(); err != nil {
		return false, domainerrors.ErrContextCanceled
	}

	if seconds <= 0 {
		return false, domainerrors.ErrInvalidArgument
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	item, exists := s.data[key]

	if !exists || item.IsExpired(time.Now().Unix()) {
		return false, nil
	}

	expiresAt := time.Now().Unix() + int64(seconds)
	item.ExpiresAt = &expiresAt

	return true, nil
}

func (s *Store) TTL(ctx context.Context, key string) (int64, error) {
	if err := ctx.Err(); err != nil {
		return 0, domainerrors.ErrContextCanceled
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	item, exists := s.data[key]
	now := time.Now().Unix()

	if !exists || item.IsExpired(now) {
		return -2, nil
	}

	if item.ExpiresAt == nil {
		return -1, nil
	}

	return *item.ExpiresAt - now, nil
}

func (s *Store) Persist(ctx context.Context, key string) (bool, error) {
	if err := ctx.Err(); err != nil {
		return false, domainerrors.ErrContextCanceled
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	item, exists := s.data[key]

	if !exists || item.IsExpired(time.Now().Unix()) {
		return false, nil
	}

	if item.ExpiresAt == nil {
		return false, nil
	}

	item.ExpiresAt = nil

	return true, nil
}

func (s *Store) Keys(ctx context.Context, pattern string) ([]string, error) {
	if err := ctx.Err(); err != nil {
		return nil, domainerrors.ErrContextCanceled
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	now := time.Now().Unix()
	matches := make([]string, 0)

	for key, item := range s.data {
		if item.IsExpired(now) {
			continue
		}

		if matchPattern(key, pattern) {
			matches = append(matches, key)
		}
	}

	return matches, nil
}

func (s *Store) Exists(ctx context.Context, key string) (bool, error) {
	if err := ctx.Err(); err != nil {
		return false, domainerrors.ErrContextCanceled
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	item, exists := s.data[key]

	if !exists {
		return false, nil
	}

	return !item.IsExpired(time.Now().Unix()), nil
}

func (s *Store) Size(ctx context.Context) (int, error) {
	if err := ctx.Err(); err != nil {
		return 0, domainerrors.ErrContextCanceled
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	now := time.Now().Unix()
	size := 0

	for _, item := range s.data {
		if !item.IsExpired(now) {
			size++
		}
	}

	return size, nil
}

func (s *Store) StartCleanup(intervalMs int64) error {
	if intervalMs <= 0 {
		return domainerrors.ErrInvalidArgument
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.cleanupStarted {
		return domainerrors.ErrCleanupRunning
	}

	s.cleanupStarted = true

	interval := time.Duration(intervalMs) * time.Millisecond

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				s.CleanupExpired()
			case <-s.stopCleanUp:
				return
			}
		}
	}()

	return nil
}

func (s *Store) StopCleanup() error {
	if !s.cleanupStarted {
		return domainerrors.ErrCleanupStopped
	}

	s.cleanupOnce.Do(func() {
		close(s.stopCleanUp)
	})

	return nil
}

func (s *Store) CleanupExpired() {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().Unix()

	for key, item := range s.data {
		if item.IsExpired(now) {
			delete(s.data, key)
		}
	}
}

func matchPattern(key, pattern string) bool {
	if pattern == "*" {
		return true
	}

	matched, err := filepath.Match(pattern, key)

	if err != nil {
		return key == pattern
	}

	return matched
}
