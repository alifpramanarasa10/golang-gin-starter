package service

import (
	"context"
	"gin-starter/common/errors"
	"gin-starter/config"
	"gin-starter/entity"
	"gin-starter/modules/user/v1/repository"

	"github.com/google/uuid"
)

// UserFinder is a service for user
type UserFinder struct {
	ufg            config.Config
	userRepo       repository.UserRepositoryUseCase
	userRoleRepo   repository.UserRoleRepositoryUseCase
	roleRepo       repository.RoleRepositoryUseCase
	permissionRepo repository.PermissionRepositoryUseCase
}

// UserFinderUseCase is a usecase for user
type UserFinderUseCase interface {
	// GetUsers gets all users
	GetUsers(ctx context.Context, query, order, sort string, limit, offset int) ([]*entity.User, int64, error)
	// GetUserByID gets a user by ID
	GetUserByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	// GetAdminUsers gets all admin users
	GetAdminUsers(ctx context.Context, query, order, sort string, limit, offset int) ([]*entity.User, int64, error)
	// GetAdminUserByID gets a admin user by ID
	GetAdminUserByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	// GetUserByEmail gets user by email
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	// GetUserByForgotPasswordToken gets user by forgot password token
	GetUserByForgotPasswordToken(ctx context.Context, token string) (*entity.User, error)
	// GetRoles gets all roles
	GetRoles(ctx context.Context, query, sort, order string, limit, offset int) ([]*entity.Role, error)
	// GetPermissions gets all permissions
	GetPermissions(ctx context.Context) ([]*entity.Permission, error)
	// GetUserPermissions gets all user permissions
	GetUserPermissions(ctx context.Context, userID uuid.UUID) ([]*entity.Permission, error)
}

// NewUserFinder creates a new UserFinder
func NewUserFinder(
	ufg config.Config,
	userRepo repository.UserRepositoryUseCase,
	userRoleRepo repository.UserRoleRepositoryUseCase,
	roleRepo repository.RoleRepositoryUseCase,
	permissionRepo repository.PermissionRepositoryUseCase,
) *UserFinder {
	return &UserFinder{
		ufg:            ufg,
		userRepo:       userRepo,
		userRoleRepo:   userRoleRepo,
		roleRepo:       roleRepo,
		permissionRepo: permissionRepo,
	}
}

// GetUsers gets all users
func (uf *UserFinder) GetUsers(ctx context.Context, query, sort, order string, limit, offset int) ([]*entity.User, int64, error) {
	users, total, err := uf.userRepo.GetUsers(ctx, query, sort, order, limit, offset)

	if err != nil {
		return nil, 0, errors.ErrInternalServerError.Error()
	}

	return users, total, nil
}

// GetAdminUsers gets all admin users
func (uf *UserFinder) GetAdminUsers(ctx context.Context, query, sort, order string, limit, offset int) ([]*entity.User, int64, error) {
	users, total, err := uf.userRepo.GetAdminUsers(ctx, query, sort, order, limit, offset)

	if err != nil {
		return nil, 0, errors.ErrInternalServerError.Error()
	}

	return users, total, nil
}

// GetAdminUserByID gets a admin user by ID
func (uf *UserFinder) GetAdminUserByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	user, err := uf.userRepo.GetUserByID(ctx, id)

	if err != nil {
		return nil, errors.ErrInternalServerError.Error()
	}

	if user == nil {
		return nil, errors.ErrRecordNotFound.Error()
	}

	return user, nil
}

// GetUserByID gets user by id
func (uf *UserFinder) GetUserByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	user, err := uf.userRepo.GetUserByID(ctx, id)

	if err != nil {
		return user, errors.ErrInternalServerError.Error()
	}

	if user == nil {
		return nil, errors.ErrRecordNotFound.Error()
	}

	return user, nil
}

// GetUserByEmail gets user by email
func (uf *UserFinder) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	user, err := uf.userRepo.GetUserByEmail(ctx, email)

	if err != nil {
		return user, errors.ErrInternalServerError.Error()
	}

	return user, nil
}

// GetUserByForgotPasswordToken gets user by forgot password token
func (uf *UserFinder) GetUserByForgotPasswordToken(ctx context.Context, token string) (*entity.User, error) {
	user, err := uf.userRepo.GetUserByForgotPasswordToken(ctx, token)

	if err != nil {
		return user, errors.ErrInternalServerError.Error()
	}

	return user, nil
}

// GetRoles gets all roles
func (uf *UserFinder) GetRoles(ctx context.Context, query, sort, order string, limit, offset int) ([]*entity.Role, error) {
	roles, err := uf.roleRepo.FindAll(ctx, query, sort, order, limit, offset)

	if err != nil {
		return nil, errors.ErrInternalServerError.Error()
	}

	return roles, nil
}

// GetPermissions get permissions
func (uf *UserFinder) GetPermissions(ctx context.Context) ([]*entity.Permission, error) {
	permissions, err := uf.permissionRepo.FindAll(ctx)
	if err != nil {
		return nil, errors.ErrInternalServerError.Error()
	}

	if permissions == nil {
		return nil, errors.ErrRecordNotFound.Error()
	}

	return permissions, err
}

// GetUserPermissions get list permission of user
func (uf *UserFinder) GetUserPermissions(ctx context.Context, userID uuid.UUID) ([]*entity.Permission, error) {
	userRole, err := uf.userRoleRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, errors.ErrInternalServerError.Error()
	}

	if userRole == nil {
		return nil, errors.ErrRecordNotFound.Error()
	}

	role, err := uf.roleRepo.FindByID(ctx, userRole.RoleID)
	if err != nil {
		return nil, errors.ErrInternalServerError.Error()
	}

	if role == nil {
		return nil, errors.ErrRecordNotFound.Error()
	}

	permissions := make([]*entity.Permission, 0)
	for _, p := range role.RolePermissions {
		permissions = append(permissions, p.Permission)
	}

	return permissions, nil
}
