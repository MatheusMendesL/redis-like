package domainerrors

import "errors"

var (
	ErrKeyNotFound     = errors.New("key not found")
	ErrInvalidArgument = errors.New("invalid argument")
	ErrInvalidCommand  = errors.New("invalid command")
	ErrContextCanceled = errors.New("context canceled")
	ErrExpired         = errors.New("key expired")
	ErrCleanupRunning  = errors.New("cleanup already running")
	ErrCleanupStopped  = errors.New("cleanup already stopped")
)
