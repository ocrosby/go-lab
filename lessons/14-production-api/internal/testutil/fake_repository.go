package testutil

import (
	"context"
	"sync"

	"github.com/ocrosby/go-lab/lessons/14-production-api/internal/domain"
)

// FakeUserRepository is an in-memory implementation of domain.UserRepository
// intended for unit tests. It behaves like the production memory repository —
// preserving Get/Create/Update/Delete semantics and returning the same domain
// errors — with two additions useful in tests:
//
//  1. Fault injection: set FailNextX to a non-nil error to make the next call
//     to method X return that error and clear the field. Use this only where
//     the test's behavior depends on the repository failing (e.g. testing
//     the service's error-path response); prefer the repo's real semantics
//     everywhere else.
//  2. State inspection: Users() returns a snapshot of the current contents
//     for assertions that need to look at the persisted state directly.
//
// Prefer this over a gomock-generated mock: it lets tests assert on
// observable outcomes ("the user is now persisted") rather than on call
// traces ("Create was called"), so tests survive refactors of the
// service's implementation.
type FakeUserRepository struct {
	mu    sync.RWMutex
	users map[string]*domain.User

	FailNextGetByID    error
	FailNextGetByEmail error
	FailNextCreate     error
	FailNextUpdate     error
	FailNextDelete     error
	FailNextList       error
}

// NewFakeUserRepository returns an empty FakeUserRepository. Pre-populate
// via Seed if the test needs a starting state.
func NewFakeUserRepository() *FakeUserRepository {
	return &FakeUserRepository{users: make(map[string]*domain.User)}
}

// Seed inserts users directly, bypassing Create. Use to prepare a starting
// state without cluttering the arrange section with repeated Create calls.
func (r *FakeUserRepository) Seed(users ...*domain.User) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, u := range users {
		copy := *u
		r.users[u.ID] = &copy
	}
}

// Users returns a snapshot of the current users, keyed by ID.
func (r *FakeUserRepository) Users() map[string]domain.User {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make(map[string]domain.User, len(r.users))
	for id, u := range r.users {
		out[id] = *u
	}
	return out
}

func (r *FakeUserRepository) GetByID(_ context.Context, id string) (*domain.User, error) {
	if err := r.consume(&r.FailNextGetByID); err != nil {
		return nil, err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, ok := r.users[id]
	if !ok {
		return nil, domain.ErrUserNotFound
	}
	copy := *u
	return &copy, nil
}

func (r *FakeUserRepository) GetByEmail(_ context.Context, email string) (*domain.User, error) {
	if err := r.consume(&r.FailNextGetByEmail); err != nil {
		return nil, err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, u := range r.users {
		if u.Email == email {
			copy := *u
			return &copy, nil
		}
	}
	return nil, domain.ErrUserNotFound
}

func (r *FakeUserRepository) Create(_ context.Context, user *domain.User) error {
	if err := r.consume(&r.FailNextCreate); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.users[user.ID]; exists {
		return domain.ErrUserAlreadyExists
	}
	copy := *user
	r.users[user.ID] = &copy
	return nil
}

func (r *FakeUserRepository) Update(_ context.Context, user *domain.User) error {
	if err := r.consume(&r.FailNextUpdate); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.users[user.ID]; !exists {
		return domain.ErrUserNotFound
	}
	copy := *user
	r.users[user.ID] = &copy
	return nil
}

func (r *FakeUserRepository) Delete(_ context.Context, id string) error {
	if err := r.consume(&r.FailNextDelete); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.users[id]; !exists {
		return domain.ErrUserNotFound
	}
	delete(r.users, id)
	return nil
}

func (r *FakeUserRepository) List(_ context.Context, limit, offset int) ([]*domain.User, error) {
	if err := r.consume(&r.FailNextList); err != nil {
		return nil, err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*domain.User, 0, len(r.users))
	i := 0
	for _, u := range r.users {
		if i >= offset && len(out) < limit {
			copy := *u
			out = append(out, &copy)
		}
		i++
	}
	return out, nil
}

// consume atomically reads and clears a FailNext* field. Returns the error
// that was set, or nil if none.
func (r *FakeUserRepository) consume(slot *error) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	err := *slot
	*slot = nil
	return err
}
