package repository

import (
	"context"
	"fmt"
	"gin-starter/common/constant"
	"gin-starter/entity"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ActivitiesRepository is a repository for Activities
type ActivitiesRepository struct {
	db *gorm.DB
}

// ActivitiesRepositoryUseCase is a use case for Activities
type ActivitiesRepositoryUseCase interface {
	// GetActivitiesByID is a function to get Activities by id
	GetActivitiesByID(ctx context.Context, id uuid.UUID) (*entity.Activities, error)
	// Update is a function to update Activities
	Update(ctx context.Context, Activities *entity.Activities) error
	// Create is a function to create Activities
	Create(ctx context.Context, Activities *entity.Activities) error
	// GetActivitiess is a function to get Activitiess
	GetActivities(ctx context.Context, query, sort, order string, limit, offset int) ([]*entity.Activities, int64, error)
	// Delete is a function to delete admin Activities
	Delete(ctx context.Context, id uuid.UUID) error
}

// NewActivitiesRepository creates a new ActivitiesRepository
func NewActivitiesRepository(db *gorm.DB) *ActivitiesRepository {
	return &ActivitiesRepository{db}
}

// GetActivitiesByID is a function to get Activities by id
func (ur *ActivitiesRepository) GetActivitiesByID(ctx context.Context, id uuid.UUID) (*entity.Activities, error) {
	result := new(entity.Activities)

	if err := ur.db.
		WithContext(ctx).
		Where("id = ?", id).
		First(result).
		Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, errors.Wrap(err, "[ActivitiesRepository-GetActivitiesByID] Activities not found")
	}

	return result, nil
}

// Update is a function to update Activities
func (ur *ActivitiesRepository) Update(ctx context.Context, Activities *entity.Activities) error {
	oldTime := Activities.UpdatedAt
	Activities.UpdatedAt = time.Now()
	if err := ur.db.
		WithContext(ctx).
		Transaction(func(tx *gorm.DB) error {
			sourceModel := new(entity.Activities)
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Find(&sourceModel, Activities.ID).Error; err != nil {
				log.Println("[ActivitiesRepository - Update]", err)
				return err
			}
			if err := tx.Model(&entity.Activities{}).
				Where(`id`, Activities.ID).
				UpdateColumns(sourceModel.MapUpdateFrom(Activities)).Error; err != nil {
				log.Println("[ActivitiesRepository - Update]", err)
				return err
			}
			return nil
		}); err != nil {
		Activities.UpdatedAt = oldTime
	}
	return nil
}

// Create is a function to create Activities
func (ur *ActivitiesRepository) Create(ctx context.Context, Activities *entity.Activities) error {
	if err := ur.db.
		WithContext(ctx).
		Model(&entity.Activities{}).
		Create(Activities).
		Error; err != nil {
		return errors.Wrap(err, "[ActivitiesRepository-Create] error while creating Activities")
	}

	return nil
}

// GetActivitiess is a function to get all Activitiess
func (ur *ActivitiesRepository) GetActivities(ctx context.Context, query, sort, order string, limit, offset int) ([]*entity.Activities, int64, error) {
	var Activities []*entity.Activities
	var total int64
	var gormDB = ur.db.
		WithContext(ctx).
		Model(&entity.Activities{})

	gormDB.Count(&total)

	gormDB = gormDB.Limit(limit).
		Offset(offset)

	if query != "" {
		gormDB = gormDB.
			Where("name ILIKE ?", "%"+query+"%")
	}

	if order != constant.Ascending && order != constant.Descending {
		order = constant.Descending
	}

	if sort == "" {
		sort = "created_at"
	}

	gormDB = gormDB.Order(fmt.Sprintf("%s %s", sort, order))

	if err := gormDB.Find(&Activities).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, 0, nil
		}
		return nil, 0, errors.Wrap(err, "[ActivitiesRepository-GetActivitiess] error when looking up all Activities")
	}

	return Activities, total, nil
}

// Delete is a function to delete Activities
func (ur *ActivitiesRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := ur.db.WithContext(ctx).
		Model(&entity.Activities{}).
		Where(`id = ?`, id).
		Delete(&entity.Activities{}, "id = ?", id).Error; err != nil {
		return errors.Wrap(err, "[ActivitiesRepository-Delete] error when deleting Activities data")
	}

	return nil
}
