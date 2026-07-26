// Package repository provides data persistence implementations for the User Management API.
package repository

import "github.com/ocrosby/go-lab/lessons/14-production-api/internal/domain"

// copyUser creates a deep copy of a user to ensure data isolation
func copyUser(user *domain.User) *domain.User {
	if user == nil {
		return nil
	}

	userCopy := *user
	return &userCopy
}
