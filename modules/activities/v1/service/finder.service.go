package service

import (
	"context"
	"gin-starter/common/errors"
	"gin-starter/config"
	"gin-starter/entity"
	"gin-starter/modules/activities/v1/repository"

	"github.com/google/uuid"
)

// ActivitiesFinder is a service for finding activities
type ActivitiesFinder struct {
	cfg            config.Config
	activitiesRepo repository.ActivitiesRepositoryUseCase
}

// ActivitiesFinderUseCase is a use case for the Activities finder
type ActivitiesFinderUseCase interface {
	// GetActivities gets all activities
	GetActivities(ctx context.Context, query, order, sort string, limit, offset int) ([]*entity.Activities, int64, error)
	// GetActivitiesByID gets a activities by ID
	GetActivitiesByID(ctx context.Context, id uuid.UUID) (*entity.Activities, error)
}

// NewActivitiesFinder is a constructor for the Activities finder
func NewActivitiesFinder(
	cfg config.Config,
	activitiesRepo repository.ActivitiesRepositoryUseCase,
) *ActivitiesFinder {
	return &ActivitiesFinder{
		cfg:            cfg,
		activitiesRepo: activitiesRepo,
	}
}

// GetActivities gets all activities
func (a *ActivitiesFinder) GetActivities(ctx context.Context, query, order, sort string, limit, offset int) ([]*entity.Activities, int64, error) {
	activities, total, err := a.activitiesRepo.GetActivities(ctx, query, order, sort, limit, offset)
	if err != nil {
		return nil, 0, errors.ErrInternalServerError.Error()
	}

	return activities, total, nil
}

// GetActivitiesByID gets a activities by ID
func (a *ActivitiesFinder) GetActivitiesByID(ctx context.Context, id uuid.UUID) (*entity.Activities, error) {
	activities, err := a.activitiesRepo.GetActivitiesByID(ctx, id)
	if err != nil {
		return nil, errors.ErrInternalServerError.Error()
	}

	return activities, nil
}
