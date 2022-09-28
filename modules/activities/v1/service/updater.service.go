package service

import (
	"context"
	"gin-starter/common/errors"
	"gin-starter/config"
	"gin-starter/entity"
	"gin-starter/modules/activities/v1/repository"
)

// ActivitiesUpdater is a service for updating activities
type ActivitiesUpdater struct {
	cfg            config.Config
	activitiesRepo repository.ActivitiesRepositoryUseCase
}

// ActivitiesUpdaterUseCase is a use case for the Activities updater
type ActivitiesUpdaterUseCase interface {
	// UpdateActivity updates an activity
	UpdateActivity(ctx context.Context, activities *entity.Activities) error
}

// NewActivitiesUpdater is a constructor for the Activities updater
func NewActivitiesUpdater(
	cfg config.Config,
	activitiesRepo repository.ActivitiesRepositoryUseCase,
) *ActivitiesUpdater {
	return &ActivitiesUpdater{
		cfg:            cfg,
		activitiesRepo: activitiesRepo,
	}
}

// UpdateActivity updates an activity
func (a *ActivitiesUpdater) UpdateActivity(ctx context.Context, activities *entity.Activities) error {
	if err := a.activitiesRepo.Update(ctx, activities); err != nil {
		return errors.ErrInternalServerError.Error()
	}

	return nil
}
