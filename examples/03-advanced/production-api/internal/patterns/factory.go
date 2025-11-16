// Package patterns demonstrates implementation of Gang of Four design patterns.
package patterns

import (
	"github.com/ocrosby/go-lab/projects/api/internal/domain"
	"github.com/ocrosby/go-lab/projects/api/internal/infrastructure/adapters/repository"
)

type RepositoryType string

const (
	MemoryRepositoryType RepositoryType = "memory"
)

type UserRepositoryFactory interface {
	CreateUserRepository(repoType RepositoryType) domain.UserRepository
}

type userRepositoryFactory struct{}

func NewUserRepositoryFactory() UserRepositoryFactory {
	return &userRepositoryFactory{}
}

func (f *userRepositoryFactory) CreateUserRepository(repoType RepositoryType) domain.UserRepository {
	switch repoType {
	case MemoryRepositoryType:
		return repository.NewMemoryUserRepository()
	default:
		return repository.NewMemoryUserRepository()
	}
}
