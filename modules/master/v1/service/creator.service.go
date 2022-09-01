package service

import (
	"gin-starter/common/interfaces"
	"gin-starter/config"
)

// MasterCreator is a struct that contains all the dependencies for the Master creator
type MasterCreator struct {
	cfg          config.Config
	cloudStorage interfaces.CloudStorageUseCase
}

// MasterCreatorUseCase is a use case for the Master creator
type MasterCreatorUseCase interface {
}

// NewMasterCreator creates a new MasterCreator
func NewMasterCreator(
	cfg config.Config,
	cloudStorage interfaces.CloudStorageUseCase,
) *MasterCreator {
	return &MasterCreator{
		cfg:          cfg,
		cloudStorage: cloudStorage,
	}
}
