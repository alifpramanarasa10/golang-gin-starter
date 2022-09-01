package repository

import (
	"context"
	"gin-starter/entity"
	"github.com/pkg/errors"

	"gorm.io/gorm"
)

// ProvinceRepository is an struct for Province repository
type ProvinceRepository struct {
	gormDB *gorm.DB
}

// ProvinceRepositoryUseCase is an interface for Province repository use case
type ProvinceRepositoryUseCase interface {
	// FindAll returns all provinces
	FindAll(ctx context.Context) ([]*entity.Province, error)
}

// NewProvinceRepository creates a new Province repository
func NewProvinceRepository(
	db *gorm.DB,
) *ProvinceRepository {
	return &ProvinceRepository{
		gormDB: db,
	}
}

// FindAll returns all provinces
func (repo *ProvinceRepository) FindAll(ctx context.Context) ([]*entity.Province, error) {
	models := make([]*entity.Province, 0)
	if err := repo.gormDB.
		WithContext(ctx).
		Model(&entity.Province{}).
		Find(&models).
		Error; err != nil {
		return nil, errors.Wrap(err, "[ProvinceRepository-FindAll]")
	}
	return models, nil
}
