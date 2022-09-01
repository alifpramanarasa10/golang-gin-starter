package repository

import (
	"context"
	"gin-starter/entity"
	"github.com/pkg/errors"

	"gorm.io/gorm"
)

type DistrictRepositoryUseCase interface {
	// FindByRegencyID finds district by regency id
	FindByRegencyID(ctx context.Context, id int64) ([]*entity.District, error)
}

// DistrictRepository is an struct for District repository
type DistrictRepository struct {
	gormDB *gorm.DB
}

// NewDistrictRepository creates a new District repository
func NewDistrictRepository(
	db *gorm.DB,
) *DistrictRepository {
	return &DistrictRepository{
		gormDB: db,
	}
}

// FindByRegencyID finds district by regency id
func (repo *DistrictRepository) FindByRegencyID(ctx context.Context, id int64) ([]*entity.District, error) {
	models := make([]*entity.District, 0)
	if err := repo.gormDB.
		WithContext(ctx).
		Model(&entity.District{}).
		Where("regency_id = ? ", id).
		Find(&models).
		Error; err != nil {
		return nil, errors.Wrap(err, "[DistrictRepository-FindByRegencyID]")
	}
	return models, nil
}
