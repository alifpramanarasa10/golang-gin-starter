package service

import (
	"context"
	"gin-starter/common/errors"
	"gin-starter/config"
	"gin-starter/modules/activities/v1/repository"

	"github.com/google/uuid"
)

// ActivitiesDeleter is a service for deleting activities
type ActivitiesDeleter struct {
	cfg            config.Config
	activitiesRepo repository.ActivitiesRepositoryUseCase
}

// ActivitiesDeleterUseCase is a use case for the Activities deleter

type ActivitiesDeleterUseCase interface {
	// DeleteActivity deletes an activity
	DeleteActivity(ctx context.Context, id uuid.UUID) error
}

// NewActivitiesDeleter is a constructor for the Activities deleter

func NewActivitiesDeleter(
	cfg config.Config,
	activitiesRepo repository.ActivitiesRepositoryUseCase,
) *ActivitiesDeleter {
	return &ActivitiesDeleter{
		cfg:            cfg,
		activitiesRepo: activitiesRepo,
	}
}

// DeleteActivity deletes an activity

func (a *ActivitiesDeleter) DeleteActivity(ctx context.Context, id uuid.UUID) error {
	if err := a.activitiesRepo.Delete(ctx, id); err != nil {
		return errors.ErrInternalServerError.Error()
	}

	return nil
}
