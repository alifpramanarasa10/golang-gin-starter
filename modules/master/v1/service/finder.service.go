package service

import (
	"context"
	"gin-starter/common/errors"
	"gin-starter/config"
	"gin-starter/entity"
	"gin-starter/modules/master/v1/repository"
)

// MasterFinder is a service for master
type MasterFinder struct {
	cfg          config.Config
	provinceRepo repository.ProvinceRepositoryUseCase
	regencyRepo  repository.RegencyRepositoryUseCase
	districtRepo repository.DistrictRepositoryUseCase
	villageRepo  repository.VillageRepositoryUseCase
}

// MasterFinderUseCase is a usecase for master
type MasterFinderUseCase interface {
	// GetProvinces returns all provinces
	GetProvinces(ctx context.Context) ([]*entity.Province, error)
	// GetRegencies returns all regencies
	GetRegencies(ctx context.Context, id int64) ([]*entity.Regency, error)
	// GetDistricts returns all districts
	GetDistricts(ctx context.Context, id int64) ([]*entity.District, error)
	// GetVillages returns all villages
	GetVillages(ctx context.Context, id int64) ([]*entity.Village, error)
}

// NewMasterFinder creates a new MasterFinder
func NewMasterFinder(
	cfg config.Config,
	provinceRepo repository.ProvinceRepositoryUseCase,
	regencyRepo repository.RegencyRepositoryUseCase,
	districtRepo repository.DistrictRepositoryUseCase,
	villageRepo repository.VillageRepositoryUseCase,
) *MasterFinder {
	return &MasterFinder{
		cfg:          cfg,
		provinceRepo: provinceRepo,
		regencyRepo:  regencyRepo,
		districtRepo: districtRepo,
		villageRepo:  villageRepo,
	}
}

// GetProvinces returns all provinces
func (s *MasterFinder) GetProvinces(ctx context.Context) ([]*entity.Province, error) {
	provinces, err := s.provinceRepo.FindAll(ctx)

	if err != nil {
		return nil, errors.ErrInternalServerError.Error()
	}

	return provinces, nil
}

// GetRegencies returns all regencies
func (s *MasterFinder) GetRegencies(ctx context.Context, id int64) ([]*entity.Regency, error) {
	regencies, err := s.regencyRepo.FindByProvinceID(ctx, id)

	if err != nil {
		return nil, errors.ErrInternalServerError.Error()
	}

	return regencies, nil
}

// GetDistricts returns all districts
func (s *MasterFinder) GetDistricts(ctx context.Context, id int64) ([]*entity.District, error) {
	districts, err := s.districtRepo.FindByRegencyID(ctx, id)

	if err != nil {
		return nil, errors.ErrInternalServerError.Error()
	}

	return districts, nil
}

// GetVillages returns all villages
func (s *MasterFinder) GetVillages(ctx context.Context, id int64) ([]*entity.Village, error) {
	villages, err := s.villageRepo.FindByDistrictID(ctx, id)

	if err != nil {
		return nil, errors.ErrInternalServerError.Error()
	}

	return villages, nil
}
