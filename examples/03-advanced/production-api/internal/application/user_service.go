// Package application contains the business logic and use cases for the User Management API.
package application

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/ocrosby/go-lab/projects/api/internal/config"
	"github.com/ocrosby/go-lab/projects/api/internal/domain"
	"github.com/ocrosby/go-lab/projects/api/internal/validation"
)

type userService struct {
	userRepo  domain.UserRepository
	logger    *zap.Logger
	validator *validation.Validator
}

func NewUserService(userRepo domain.UserRepository, logger *zap.Logger) domain.UserService {
	return &userService{
		userRepo:  userRepo,
		logger:    logger,
		validator: validation.NewValidator(),
	}
}

func (s *userService) CreateUser(ctx context.Context, email, name string) (*domain.User, error) {
	if err := s.validator.ValidateUserCreation(email, name); err != nil {
		return nil, s.validator.WrapValidationError(err)
	}

	existingUser, err := s.userRepo.GetByEmail(ctx, email)
	if err == nil && existingUser != nil {
		return nil, domain.ErrUserAlreadyExists
	}

	user := &domain.User{
		ID:        s.generateID(),
		Email:     email,
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		s.logger.Error("failed to create user", zap.Error(err), zap.String("email", email))
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	s.logger.Info("user created successfully", zap.String("user_id", user.ID))
	return user, nil
}

func (s *userService) GetUser(ctx context.Context, id string) (*domain.User, error) {
	if err := s.validator.ValidateUserID(id); err != nil {
		return nil, s.validator.WrapValidationError(err)
	}

	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		if err == domain.ErrUserNotFound {
			return nil, domain.ErrUserNotFound
		}
		s.logger.Error("failed to get user", zap.Error(err), zap.String("id", id))
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (s *userService) UpdateUser(ctx context.Context, id, name string) (*domain.User, error) {
	if err := s.validator.ValidateUserUpdate(id, name); err != nil {
		return nil, s.validator.WrapValidationError(err)
	}

	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		if err == domain.ErrUserNotFound {
			return nil, domain.ErrUserNotFound
		}
		s.logger.Error("failed to get user for update", zap.Error(err), zap.String("id", id))
		return nil, fmt.Errorf("failed to get user for update: %w", err)
	}

	user.Name = name
	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(ctx, user); err != nil {
		s.logger.Error("failed to update user", zap.Error(err), zap.String("id", id))
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	s.logger.Info("user updated successfully", zap.String("user_id", user.ID))
	return user, nil
}

func (s *userService) DeleteUser(ctx context.Context, id string) error {
	if err := s.validator.ValidateUserID(id); err != nil {
		return s.validator.WrapValidationError(err)
	}

	if err := s.userRepo.Delete(ctx, id); err != nil {
		if err == domain.ErrUserNotFound {
			return domain.ErrUserNotFound
		}
		s.logger.Error("failed to delete user", zap.Error(err), zap.String("id", id))
		return fmt.Errorf("failed to delete user: %w", err)
	}

	s.logger.Info("user deleted successfully", zap.String("user_id", id))
	return nil
}

func (s *userService) ListUsers(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	if limit < 1 {
		limit = config.DefaultPaginationLimit
	}
	if offset < 0 {
		offset = config.DefaultPaginationOffset
	}

	users, err := s.userRepo.List(ctx, limit, offset)
	if err != nil {
		s.logger.Error("failed to list users", zap.Error(err), zap.Int("limit", limit), zap.Int("offset", offset))
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	return users, nil
}

func (s *userService) generateID() string {
	return fmt.Sprintf("user_%d", time.Now().UnixNano())
}
