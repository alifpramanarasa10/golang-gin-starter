package repository

import (
	"context"
	"encoding/json"
	"fmt"
	commonCache "gin-starter/common/cache"
	"gin-starter/common/errors"
	"gin-starter/common/interfaces"
	"gin-starter/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RolePermissionRepository is a repository for role permission
type RolePermissionRepository struct {
	db    *gorm.DB
	cache interfaces.Cacheable
}

// RolePermissionRepositoryUseCase is a use case for role permission
type RolePermissionRepositoryUseCase interface {
	// FindByRoleIDAndPermissionID is a method for finding role permission by role id and permission id
	FindByRoleIDAndPermissionID(ctx context.Context, roleID uuid.UUID, permissionID uuid.UUID) (rolePermission *entity.RolePermission, err error)
}

// NewRolePermissionRepository is a constructor for RolePermissionRepository
func NewRolePermissionRepository(db *gorm.DB, cache interfaces.Cacheable) *RolePermissionRepository {
	return &RolePermissionRepository{db, cache}
}

// FindByRoleIDAndPermissionID is a use case for finding role permission by role id and permission id
func (rpr *RolePermissionRepository) FindByRoleIDAndPermissionID(ctx context.Context, roleID uuid.UUID, permissionID uuid.UUID) (rolePermission *entity.RolePermission, err error) {
	rolePermission = &entity.RolePermission{}
	bytes, _ := rpr.cache.Get(fmt.Sprintf(
		commonCache.RolePermissionFindByRoleIDAndPermissionID,
		roleID,
		permissionID,
	))

	if bytes != nil {
		if err := json.Unmarshal(bytes, &rolePermission); err != nil {
			return nil, err
		}
		return rolePermission, nil
	}

	if err := rpr.db.WithContext(ctx).
		Model(&entity.RolePermission{}).
		Where("role_id = ? AND permission_id = ?", roleID, permissionID).
		First(&rolePermission).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, errors.ErrRecordNotFound.Error()
	}

	if err := rpr.cache.Set(fmt.Sprintf(
		commonCache.RolePermissionFindByRoleIDAndPermissionID,
		roleID,
		permissionID,
	), &rolePermission, commonCache.OneMonth); err != nil {
		return nil, err
	}

	return rolePermission, nil
}
