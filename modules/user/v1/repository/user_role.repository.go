package repository

import (
	"context"
	"encoding/json"
	"fmt"
	commonCache "gin-starter/common/cache"
	"gin-starter/common/interfaces"
	"gin-starter/entity"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// UserRoleRepository is a repository for user role
type UserRoleRepository struct {
	db    *gorm.DB
	cache interfaces.Cacheable
}

// UserRoleRepositoryUseCase is a use case for user role
type UserRoleRepositoryUseCase interface {
	// CreateOrUpdate is a method for creating or updating user role
	CreateOrUpdate(ctx context.Context, userRole *entity.UserRole) error
	// FindByUserID is a method for finding user role by user id
	FindByUserID(ctx context.Context, id uuid.UUID) (*entity.UserRole, error)
	// Update is a method for updating user role
	Update(ctx context.Context, userRole *entity.UserRole) error
	// Delete is a method for deleting user role
	Delete(ctx context.Context, id uuid.UUID) error
}

// NewUserRoleRepository is a constructor for UserRoleRepository
func NewUserRoleRepository(db *gorm.DB, cache interfaces.Cacheable) *UserRoleRepository {
	return &UserRoleRepository{db, cache}
}

// CreateOrUpdate is a method for creating or updating user role
func (nc *UserRoleRepository) CreateOrUpdate(ctx context.Context, userRole *entity.UserRole) error {
	var find *entity.UserRole

	findUser := nc.db.
		Where("user_id = ?", userRole.UserID).
		First(&find)

	if err := findUser.Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}

	if findUser.RowsAffected > 0 {
		if err := nc.db.Model(&entity.UserRole{}).
			Where("user_id = ?", userRole.UserID).
			UpdateColumns(map[string]interface{}{
				"role_id": userRole.RoleID,
			}).
			Error; err != nil {
			return err
		}

		return nil
	}

	if err := nc.db.
		WithContext(ctx).
		Model(&entity.UserRole{}).
		Create(userRole).
		Error; err != nil {
		return errors.Wrap(err, "[UserRoleRepository-CreateNews] error while creating user")
	}

	if err := nc.cache.BulkRemove(fmt.Sprintf(commonCache.UserRoleByUserID, "*")); err != nil {
		return err
	}

	return nil
}

// FindByUserID is a method for finding user role by user id
func (nc *UserRoleRepository) FindByUserID(ctx context.Context, id uuid.UUID) (*entity.UserRole, error) {
	category := &entity.UserRole{}

	bytes, _ := nc.cache.Get(fmt.Sprintf(
		commonCache.UserRoleByUserID, id.String()))

	if bytes != nil {
		if err := json.Unmarshal(bytes, &category); err != nil {
			return nil, err
		}
		return category, nil
	}

	if err := nc.db.
		WithContext(ctx).
		Model(&entity.UserRole{}).
		Preload("Role").
		Preload("User").
		Where("user_id = ?", id).
		First(&category).
		Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, errors.Wrap(err, "[NewsRepository-FindByID] error while getting category category")
	}

	if err := nc.cache.Set(fmt.Sprintf(
		commonCache.UserRoleByUserID, id.String()), &category, commonCache.OneMonth); err != nil {
		return nil, err
	}

	return category, nil
}

// Update is a method for updating user role
func (nc *UserRoleRepository) Update(ctx context.Context, userRole *entity.UserRole) error {
	oldTime := userRole.UpdatedAt
	userRole.UpdatedAt = time.Now()
	if err := nc.db.
		WithContext(ctx).
		Transaction(func(tx *gorm.DB) error {
			sourceModel := new(entity.UserRole)
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
				Where("user_id = ?", userRole.UserID).
				Find(&sourceModel).Error; err != nil {
				log.Println("[GamPTKRepository - Update]", err)
				return err
			}
			if err := tx.Model(&entity.UserRole{}).
				Where(`user_id`, userRole.UserID).
				UpdateColumns(sourceModel.MapUpdateFrom(userRole)).Error; err != nil {
				log.Println("[GamPTKRepository - Update]", err)
				return err
			}
			return nil
		}); err != nil {
		userRole.UpdatedAt = oldTime
	}

	if err := nc.cache.BulkRemove(fmt.Sprintf(commonCache.UserRoleByUserID, "*")); err != nil {
		return err
	}

	return nil
}

// Delete is a method for deleting user role
func (nc *UserRoleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := nc.db.WithContext(ctx).
		Model(&entity.UserRole{}).
		Where(`user_id = ?`, id).
		Updates(
			map[string]interface{}{
				"updated_at": time.Now(),
				"deleted_at": time.Now(),
			}).Error; err != nil {
		return errors.Wrap(err, "[UserRepository-DeactivateUser] error when updating user data")
	}

	if err := nc.cache.BulkRemove(fmt.Sprintf(commonCache.UserRoleByUserID, "*")); err != nil {
		return err
	}
	return nil
}
