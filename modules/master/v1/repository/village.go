package repository

import (
	"context"
	"gin-starter/entity"
	"github.com/pkg/errors"

	"gorm.io/gorm"
)

type VillageRepositoryUseCase interface {
	// FindByDistrictID finds village by district id
	FindByDistrictID(ctx context.Context, id int64) ([]*entity.Village, error)
}

// VillageRepository is an struct for Village repository
type VillageRepository struct {
	gormDB *gorm.DB
}

// NewVillageRepository creates a new Village repository
func NewVillageRepository(
	db *gorm.DB,
) *VillageRepository {
	return &VillageRepository{
		gormDB: db,
	}
}

// FindByDistrictID finds village by district id
func (repo *VillageRepository) FindByDistrictID(ctx context.Context, id int64) ([]*entity.Village, error) {
	models := make([]*entity.Village, 0)
	if err := repo.gormDB.
		WithContext(ctx).
		Model(&entity.Village{}).
		Where("district_id = ? ", id).
		Find(&models).
		Error; err != nil {
		return nil, errors.Wrap(err, "[VillageRepository-FindByDistrictID]")
	}
	return models, nil
}
