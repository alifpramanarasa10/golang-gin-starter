package repository

import (
	"context"
	"encoding/json"
	"fmt"
	commonCache "gin-starter/common/cache"
	"gin-starter/common/interfaces"
	"gin-starter/entity"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// RoleRepository is a repository for role
type RoleRepository struct {
	db    *gorm.DB
	cache interfaces.Cacheable
}

// RoleRepositoryUseCase is a use case for role
type RoleRepositoryUseCase interface {
	// Create creates a role
	Create(ctx context.Context, role *entity.Role, permissionIDs []uuid.UUID) error
	// FindAll finds all roles
	FindAll(ctx context.Context, query, sort, order string, limit, offset int) ([]*entity.Role, error)
	// FindByID finds a role by id
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Role, error)
	// Delete deletes a role
	Delete(ctx context.Context, id uuid.UUID, deletedBy string) error
	// FindByName finds a role by name
	FindByName(ctx context.Context, slug string) (*entity.Role, error)
	// Update update a role
	Update(ctx context.Context, role *entity.Role, rolePermissions []*entity.RolePermission) error
}

// NewRoleRepository creates a new role repository
func NewRoleRepository(db *gorm.DB, cache interfaces.Cacheable) *RoleRepository {
	return &RoleRepository{db, cache}
}

// Create creates a role
func (nc *RoleRepository) Create(ctx context.Context, role *entity.Role, permissionIDs []uuid.UUID) error {
	if err := nc.db.
		WithContext(ctx).
		Transaction(func(tx *gorm.DB) error {
			if err := tx.Model(&entity.Role{}).Create(role).Error; err != nil {
				return errors.Wrap(err, "[RoleRepository-CreateRole] error while creating role")
			}

			if len(permissionIDs) > 0 {
				for _, pid := range permissionIDs {
					rolePermission := entity.NewRolePermission(uuid.New(), role.ID, pid, role.CreatedBy.String)
					rolePermission.UpdatedBy = role.UpdatedBy
					if err := tx.Model(&entity.RolePermission{}).Create(rolePermission).Error; err != nil {
						return errors.Wrap(err, "[RoleRepository-CreateRolePermission] error while creating role permission")
					}
				}
			}
			return nil
		}); err != nil {
		return errors.Wrap(err, "[RoleRepository-CreateRole] error while creating role")
	}

	if err := nc.cache.BulkRemove(fmt.Sprintf(commonCache.RolePermissionFindByRoleIDAndPermissionID, "*", "*")); err != nil {
		return err
	}

	if err := nc.cache.BulkRemove(fmt.Sprintf(commonCache.RoleFindByID, "*")); err != nil {
		return err
	}

	return nil
}

// FindByID finds a role by id
func (nc *RoleRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Role, error) {
	role := &entity.Role{}

	bytes, _ := nc.cache.Get(fmt.Sprintf(
		commonCache.RoleFindByID,
		id,
	))

	if bytes != nil {
		if err := json.Unmarshal(bytes, &role); err != nil {
			return nil, err
		}
		return role, nil
	}

	if err := nc.db.
		WithContext(ctx).
		Model(&entity.Role{}).
		Preload("RolePermissions").
		Preload("RolePermissions.Permission").
		Where("id = ?", id).
		First(&role).
		Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, errors.Wrap(err, "[RoleRepository-FindByID] error while getting role")
	}

	if err := nc.cache.Set(fmt.Sprintf(
		commonCache.RoleFindByID,
		id,
	), role, commonCache.OneMonth); err != nil {
		return nil, err
	}

	return role, nil
}

// FindAll finds all roles
func (nc *RoleRepository) FindAll(ctx context.Context, query, sort, order string, limit, offset int) ([]*entity.Role, error) {
	role := make([]*entity.Role, 0)
	var gormDB = nc.db.
		WithContext(ctx).
		Model(&entity.Role{}).
		Preload("RolePermissions").
		Preload("RolePermissions.Permission")

	if query != "" {
		gormDB.Where("name ILIKE ?", "%"+query+"%")
	}

	if sort != "" {
		gormDB.Order(fmt.Sprintf("%s %s", sort, order))
	}

	if limit > 0 {
		gormDB.Limit(limit)
	}

	if offset > 0 {
		gormDB.Offset(offset)
	}

	if err := gormDB.
		Find(&role).
		Error; err != nil {
		return nil, errors.Wrap(err, "[RoleRepository-GetNewsCategories] error while getting news category")
	}

	return role, nil
}

// Delete deletes a role
func (nc *RoleRepository) Delete(ctx context.Context, id uuid.UUID, deletedBy string) error {
	if err := nc.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// soft delete role
		if err := nc.db.WithContext(ctx).
			Model(&entity.Role{}).
			Where(`id = ?`, id).
			Updates(
				map[string]interface{}{
					"deleted_by": deletedBy,
					"updated_at": time.Now(),
					"deleted_at": time.Now(),
				}).Error; err != nil {
			return errors.Wrap(err, "[RoleRepository-DeactivateRole] error when updating role data")
		}

		if err := nc.db.WithContext(ctx).
			Model(&entity.RolePermission{}).
			Where(`role_id = ?`, id).
			Updates(
				map[string]interface{}{
					"deleted_by": deletedBy,
					"updated_at": time.Now(),
					"deleted_at": time.Now(),
				}).Error; err != nil {
			return errors.Wrap(err, "[RoleRepository-DeactivateRole] error when updating role data")
		}

		return nil
	}); err != nil {
		return errors.Wrap(err, "[RoleRepository-Delete] error while deleting role")
	}

	if err := nc.cache.BulkRemove(fmt.Sprintf(commonCache.RolePermissionFindByRoleIDAndPermissionID, "*", "*")); err != nil {
		return err
	}

	if err := nc.cache.BulkRemove(fmt.Sprintf(commonCache.RoleFindByID, "*")); err != nil {
		return err
	}

	return nil
}

// FindByName finds a role by name
func (nc *RoleRepository) FindByName(ctx context.Context, name string) (*entity.Role, error) {
	role := &entity.Role{}

	if err := nc.db.
		WithContext(ctx).
		Model(&entity.Role{}).
		Preload("RolePermissions").
		Preload("RolePermissions.Permission").
		Where("name = ?", name).
		First(&role).
		Error; err != nil {
		return nil, errors.Wrap(err, "[RoleRepository-FindByName] error while getting role")
	}

	return role, nil
}

// Update update a role
func (nc *RoleRepository) Update(ctx context.Context, role *entity.Role, rolePermissions []*entity.RolePermission) error {
	oldTime := role.UpdatedAt
	role.UpdatedAt = time.Now()
	if err := nc.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// update role
		sourceModelNews := new(entity.Role)
		if err := tx.Clauses(clause.Locking{
			Strength: "UPDATE",
		}).Find(&sourceModelNews, role.ID).Error; err != nil {
			return errors.Wrap(err, "[RoleRepository-Update] error while finding role")
		}

		if err := tx.Model(&entity.Role{}).
			Where("id = ?", role.ID).
			UpdateColumns(sourceModelNews.MapUpdateFrom(role)).
			Error; err != nil {
			return errors.Wrap(err, "[RoleRepository-Update] error while updating role")
		}

		// delete role permission
		if err := tx.Model(&entity.RolePermission{}).
			Where("role_id = ?", role.ID).
			Unscoped().
			Delete(&entity.RolePermission{}).
			Error; err != nil {
			return errors.Wrap(err, "[RoleRepository-Update] error while deleting role permission")
		}

		// update role permission
		if len(rolePermissions) > 0 {
			for _, rp := range rolePermissions {
				if err := tx.WithContext(ctx).Model(&entity.RolePermission{}).Create(rp).Error; err != nil {
					return errors.Wrap(err, "[RoleRepository-Create] error while creating role permission")
				}
			}
		}

		return nil
	}); err != nil {
		role.UpdatedAt = oldTime
	}

	if err := nc.cache.BulkRemove(fmt.Sprintf(commonCache.RolePermissionFindByRoleIDAndPermissionID, "*", "*")); err != nil {
		return err
	}

	if err := nc.cache.BulkRemove(fmt.Sprintf(commonCache.RoleFindByID, "*")); err != nil {
		return err
	}

	return nil
}
