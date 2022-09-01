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

// PermissionRepository is a repository for permission
type PermissionRepository struct {
	db    *gorm.DB
	cache interfaces.Cacheable
}

// PermissionRepositoryUseCase is an interface for permission repository
type PermissionRepositoryUseCase interface {
	// Create creates a permission
	Create(ctx context.Context, permission *entity.Permission) error
	// FindAll finds all permission
	FindAll(ctx context.Context) ([]*entity.Permission, error)
	// FindByName finds permission by name
	FindByName(ctx context.Context, name string) (*entity.Permission, error)
	// FindByID finds permission by id
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Permission, error)
	// Update updates permission
	Update(ctx context.Context, permission *entity.Permission) error
}

// NewPermissionRepository `NewPermissionRepository` is a function that returns a pointer to a `PermissionRepository` struct
func NewPermissionRepository(db *gorm.DB, cache interfaces.Cacheable) *PermissionRepository {
	return &PermissionRepository{db, cache}
}

// Create creates a permission
func (pr *PermissionRepository) Create(ctx context.Context, permission *entity.Permission) error {
	if err := pr.db.
		WithContext(ctx).
		Model(&entity.Permission{}).
		Create(permission).
		Error; err != nil {
		return errors.Wrap(err, "[PermissionRepository-CreatePermission] error while creating permission")
	}

	if err := pr.cache.BulkRemove(fmt.Sprintf(commonCache.PermissionFindByName, "*")); err != nil {
		return err
	}
	return nil
}

// FindAll finds all permission
func (pr *PermissionRepository) FindAll(ctx context.Context) ([]*entity.Permission, error) {
	permissions := make([]*entity.Permission, 0)
	if err := pr.db.WithContext(ctx).
		Model(&entity.Permission{}).
		Find(&permissions).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, errors.Wrap(err, "[PermissionRepository-FindAll] error while getting permission")
	}

	return permissions, nil
}

// FindByName finds permission by name
func (pr *PermissionRepository) FindByName(ctx context.Context, name string) (*entity.Permission, error) {
	permission := &entity.Permission{}

	bytes, _ := pr.cache.Get(fmt.Sprintf(
		commonCache.PermissionFindByName,
		name,
	))

	if bytes != nil {
		if err := json.Unmarshal(bytes, &permission); err != nil {
			return nil, err
		}
		return permission, nil
	}

	if err := pr.db.WithContext(ctx).
		Model(&entity.Permission{}).
		Where("name = ?", name).
		First(permission).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, errors.Wrap(err, "[PermissionRepository-FindByName] error while getting permission")
	}

	if err := pr.cache.Set(fmt.Sprintf(
		commonCache.PermissionFindByName,
		name), &permission, commonCache.OneMonth); err != nil {
		return nil, err
	}

	return permission, nil
}

// FindByID finds permission by id
func (pr *PermissionRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Permission, error) {
	permission := &entity.Permission{}
	if err := pr.db.WithContext(ctx).
		Model(&entity.Permission{}).
		Where("id = ?", id).
		First(&permission).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, errors.Wrap(err, "[PermissionRepository-FindByID] error while getting permission")
	}

	return permission, nil
}

// Update updates permission
func (pr *PermissionRepository) Update(ctx context.Context, permission *entity.Permission) error {
	oldTime := permission.UpdatedAt
	permission.UpdatedAt = time.Now()
	if err := pr.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		sourceModelNews := new(entity.Permission)
		if err := tx.Clauses(clause.Locking{
			Strength: "UPDATE",
		}).Find(&sourceModelNews, permission.ID).Error; err != nil {
			return errors.Wrap(err, "[PermissionRepository-Update] error while finding permission")
		}

		if err := tx.Model(&entity.Permission{}).
			Where("id = ?", permission.ID).
			UpdateColumns(sourceModelNews.MapUpdateFrom(permission)).
			Error; err != nil {
			return errors.Wrap(err, "[PermissionRepository-Update] error while updating permission")
		}

		return nil
	}); err != nil {
		permission.UpdatedAt = oldTime
	}
	if err := pr.cache.BulkRemove(fmt.Sprintf(commonCache.PermissionFindByName, "*")); err != nil {
		return err
	}
	return nil
}
