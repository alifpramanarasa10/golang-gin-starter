package service

import (
	"context"
	"gin-starter/common/errors"
	"gin-starter/common/interfaces"
	"gin-starter/config"
	"gin-starter/entity"
	notificationService "gin-starter/modules/notification/v1/service"
	"gin-starter/modules/user/v1/repository"
	"gin-starter/utils"
	"time"

	"github.com/google/uuid"
)

// UserCreator is a struct that contains all the dependencies for the User creator
type UserCreator struct {
	cfg            config.Config
	userRepo       repository.UserRepositoryUseCase
	userRoleRepo   repository.UserRoleRepositoryUseCase
	roleRepo       repository.RoleRepositoryUseCase
	permissionRepo repository.PermissionRepositoryUseCase
	notifCreator   notificationService.NotificationCreatorUseCase
	cloudStorage   interfaces.CloudStorageUseCase
}

// UserCreatorUseCase is a use case for the User creator
type UserCreatorUseCase interface {
	// CreateUser creates a new user
	CreateUser(ctx context.Context, name, email, password, phoneNumber, photo string, dob time.Time) (*entity.User, error)
	// CreateAdmin creates a new admin
	CreateAdmin(ctx context.Context, name, email, password, phoneNumber, photo string, dob time.Time, roleID uuid.UUID) (*entity.User, error)
	// CreatePermission creates a permission
	CreatePermission(ctx context.Context, name, label string) (*entity.Permission, error)
	// CreateRole creates a role
	CreateRole(ctx context.Context, name string, permissionIDs []uuid.UUID, createdBy string) (*entity.Role, error)
}

// NewUserCreator is a constructor for the User creator
func NewUserCreator(
	cfg config.Config,
	userRepo repository.UserRepositoryUseCase,
	userRoleRepo repository.UserRoleRepositoryUseCase,
	roleRepo repository.RoleRepositoryUseCase,
	permissionRepo repository.PermissionRepositoryUseCase,
	notifCreator notificationService.NotificationCreatorUseCase,
	cloudStorage interfaces.CloudStorageUseCase,
) *UserCreator {
	return &UserCreator{
		cfg:            cfg,
		userRepo:       userRepo,
		userRoleRepo:   userRoleRepo,
		roleRepo:       roleRepo,
		permissionRepo: permissionRepo,
		notifCreator:   notifCreator,
		cloudStorage:   cloudStorage,
	}
}

// CreateUser creates a new user
func (uc *UserCreator) CreateUser(ctx context.Context, name, email, password, phoneNumber, photo string, dob time.Time) (*entity.User, error) {
	user := entity.NewUser(
		uuid.New(),
		name,
		email,
		password,
		utils.TimeToNullTime(dob),
		photo,
		phoneNumber,
		"system",
	)

	if err := uc.userRepo.CreateUser(ctx, user); err != nil {
		return nil, errors.ErrInternalServerError.Error()
	}

	return user, nil
}

// CreateAdmin creates a new admin
func (uc *UserCreator) CreateAdmin(ctx context.Context, name, email, password, phoneNumber, photo string, dob time.Time, roleID uuid.UUID) (*entity.User, error) {
	userID := uuid.New()
	user := entity.NewUser(
		userID,
		name,
		email,
		password,
		utils.TimeToNullTime(dob),
		photo,
		phoneNumber,
		"system",
	)

	if err := uc.userRepo.CreateUser(ctx, user); err != nil {
		return nil, errors.ErrInternalServerError.Error()
	}

	userRole := entity.NewUserRole(uuid.New(), userID, roleID, "system")

	if err := uc.userRoleRepo.CreateOrUpdate(ctx, userRole); err != nil {
		return nil, errors.ErrInternalServerError.Error()
	}

	return user, nil
}

// CreatePermission creates a permission
func (uc *UserCreator) CreatePermission(ctx context.Context, name, label string) (*entity.Permission, error) {
	permission := entity.NewPermission(
		uuid.New(),
		name,
		label,
		"system",
	)

	if err := uc.permissionRepo.Create(ctx, permission); err != nil {
		return nil, err
	}

	return permission, nil
}

// CreateRole creates a role
func (uc *UserCreator) CreateRole(ctx context.Context, name string, permissionIDs []uuid.UUID, createdBy string) (*entity.Role, error) {
	role := entity.NewRole(uuid.New(), name, createdBy)
	if err := uc.roleRepo.Create(ctx, role, permissionIDs); err != nil {
		return nil, err
	}

	return role, nil
}
