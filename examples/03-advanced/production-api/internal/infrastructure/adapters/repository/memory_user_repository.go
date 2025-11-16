// Package repository provides data persistence implementations for the User Management API.
package repository

import (
	"context"
	"sync"

	"github.com/ocrosby/go-lab/projects/api/internal/domain"
)

type memoryUserRepository struct {
	users map[string]*domain.User
	mutex sync.RWMutex
}

func NewMemoryUserRepository() domain.UserRepository {
	return &memoryUserRepository{
		users: make(map[string]*domain.User),
	}
}

func (r *memoryUserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, domain.ErrUserNotFound
	}

	return copyUser(user), nil
}

func (r *memoryUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, user := range r.users {
		if user.Email == email {
			return copyUser(user), nil
		}
	}

	return nil, domain.ErrUserNotFound
}

func (r *memoryUserRepository) Create(ctx context.Context, user *domain.User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.users[user.ID]; exists {
		return domain.ErrUserAlreadyExists
	}

	r.users[user.ID] = copyUser(user)
	return nil
}

func (r *memoryUserRepository) Update(ctx context.Context, user *domain.User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.users[user.ID]; !exists {
		return domain.ErrUserNotFound
	}

	userCopy := *user
	r.users[user.ID] = &userCopy
	return nil
}

func (r *memoryUserRepository) Delete(ctx context.Context, id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.users[id]; !exists {
		return domain.ErrUserNotFound
	}

	delete(r.users, id)
	return nil
}

func (r *memoryUserRepository) List(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var users []*domain.User
	i := 0

	for _, user := range r.users {
		if i >= offset && len(users) < limit {
			users = append(users, copyUser(user))
		}
		i++
	}

	return users, nil
}
