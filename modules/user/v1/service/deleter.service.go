package service

import (
	"context"
	"gin-starter/common/errors"
	"gin-starter/config"
	"gin-starter/modules/user/v1/repository"

	"github.com/google/uuid"
)

// UserDeleter is a service for user
type UserDeleter struct {
	cfg      config.Config
	userRepo repository.UserRepositoryUseCase
	roleRepo repository.RoleRepositoryUseCase
}

// UserDeleterUseCase is a use case for user
type UserDeleterUseCase interface {
	// DeleteAdmin deletes admin
	DeleteAdmin(ctx context.Context, id uuid.UUID) error
	// DeleteRole deletes role
	DeleteRole(ctx context.Context, id uuid.UUID, deletedBy string) error
}

// NewUserDeleter creates a new UserDeleter
func NewUserDeleter(
	cfg config.Config,
	userRepo repository.UserRepositoryUseCase,
	roleRepo repository.RoleRepositoryUseCase,
) *UserDeleter {
	return &UserDeleter{
		cfg:      cfg,
		userRepo: userRepo,
		roleRepo: roleRepo,
	}
}

// DeleteAdmin deletes admin
func (ud *UserDeleter) DeleteAdmin(ctx context.Context, id uuid.UUID) error {
	if err := ud.userRepo.DeleteAdmin(ctx, id); err != nil {
		return errors.ErrInternalServerError.Error()
	}

	return nil
}

// DeleteRole deletes role
func (ud *UserDeleter) DeleteRole(ctx context.Context, id uuid.UUID, deletedBy string) error {
	if err := ud.roleRepo.Delete(ctx, id, deletedBy); err != nil {
		return errors.ErrInternalServerError.Error()
	}

	return nil
}
