package service

import (
	"gin-starter/common/interfaces"
	"gin-starter/config"
)

// MasterUpdater is a struct that contains all the dependencies for the Master creator
type MasterUpdater struct {
	cfg          config.Config
	cloudStorage interfaces.CloudStorageUseCase
}

// MasterUpdaterUseCase is a use case for the Master creator
type MasterUpdaterUseCase interface {
}

// NewMasterUpdater creates a new MasterUpdater
func NewMasterUpdater(
	cfg config.Config,
	cloudStorage interfaces.CloudStorageUseCase,
) *MasterUpdater {
	return &MasterUpdater{
		cfg:          cfg,
		cloudStorage: cloudStorage,
	}
}
