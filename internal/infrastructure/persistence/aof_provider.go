package persistence

import (
	"redis-project/internal/domain/repository"
)

type AOFProviderOption struct {
	EnableAOF bool
	Filepath  string
}

func NewAOFProvider(opt AOFProviderOption) (repository.PersistenceRepository, error) {
	if !opt.EnableAOF {
		return nil, nil
	}
	return NewAOF(opt.Filepath)
}
