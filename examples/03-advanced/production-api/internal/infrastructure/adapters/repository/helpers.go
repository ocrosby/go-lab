// Package repository provides data persistence implementations for the User Management API.
package repository

import "github.com/ocrosby/go-lab/projects/api/internal/domain"

// copyUser creates a deep copy of a user to ensure data isolation
func copyUser(user *domain.User) *domain.User {
	if user == nil {
		return nil
	}

	userCopy := *user
	return &userCopy
}

// copyUsers creates deep copies of a slice of users
func copyUsers(users []*domain.User) []*domain.User {
	if users == nil {
		return nil
	}

	copies := make([]*domain.User, len(users))
	for i, user := range users {
		copies[i] = copyUser(user)
	}

	return copies
}
