package repository

import (
	"context"
	"gin-starter/entity"
	"github.com/pkg/errors"

	"gorm.io/gorm"
)

type RegencyRepositoryUseCase interface {
	// FindByProvinceID finds regency by province id
	FindByProvinceID(ctx context.Context, id int64) ([]*entity.Regency, error)
}

// RegencyRepository is an struct for Regency repository
type RegencyRepository struct {
	gormDB *gorm.DB
}

// NewRegencyRepository creates a new Regency repository
func NewRegencyRepository(
	db *gorm.DB,
) *RegencyRepository {
	return &RegencyRepository{
		gormDB: db,
	}
}

// FindByProvinceID finds regency by province id
func (repo *RegencyRepository) FindByProvinceID(ctx context.Context, id int64) ([]*entity.Regency, error) {
	models := make([]*entity.Regency, 0)
	if err := repo.gormDB.
		WithContext(ctx).
		Model(&entity.Regency{}).
		Where("province_id = ? ", id).
		Find(&models).
		Error; err != nil {
		return nil, errors.Wrap(err, "[RegencyRepository-FindByProvinceID]")
	}
	return models, nil
}
