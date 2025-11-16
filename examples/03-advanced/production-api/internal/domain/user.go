// Package domain contains the core business entities and interfaces for the User Management API.
package domain

import (
	"context"
	"time"
)

type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserRepository interface {
	GetByID(ctx context.Context, id string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]*User, error)
}

type UserService interface {
	CreateUser(ctx context.Context, email, name string) (*User, error)
	GetUser(ctx context.Context, id string) (*User, error)
	UpdateUser(ctx context.Context, id, name string) (*User, error)
	DeleteUser(ctx context.Context, id string) error
	ListUsers(ctx context.Context, limit, offset int) ([]*User, error)
}
