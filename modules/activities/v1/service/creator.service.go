package service

import (
	"context"
	"gin-starter/config"
	"gin-starter/entity"
	"gin-starter/modules/activities/v1/repository"

	"github.com/google/uuid"
)

// ActivitiesCreator is a service for creating activities
type ActivitiesCreator struct {
	cfg            config.Config
	activitiesRepo repository.ActivitiesRepositoryUseCase
}

// ActivitiesCreatorUseCase is a use case for the Activities creator
type ActivitiesCreatorUseCase interface {
	// CreateActivities creates a new activity
	CreateActivities(ctx context.Context, userID uuid.UUID, title, description, activities_type string) (*entity.Activities, error)
}

// NewActivitiesCreator is a constructor for the Activities creator
func NewActivitiesCreator(
	cfg config.Config,
	activitiesRepo repository.ActivitiesRepositoryUseCase,
) *ActivitiesCreator {
	return &ActivitiesCreator{
		cfg:            cfg,
		activitiesRepo: activitiesRepo,
	}
}

// CreateActivities creates a new activity
func (a *ActivitiesCreator) CreateActivities(ctx context.Context, userID uuid.UUID, title, description, activities_type string) (*entity.Activities, error) {
	activities := entity.NewActivities(
		uuid.New(),
		userID,
		title,
		description,
		activities_type,
		"system",
	)

	err := a.activitiesRepo.Create(ctx, activities)
	if err != nil {
		return nil, err
	}

	return activities, nil
}
