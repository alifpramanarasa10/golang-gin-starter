package service

import (
	"gin-starter/common/interfaces"
	"gin-starter/config"
)

// MasterDeleter is a struct that contains all the dependencies for the Master creator
type MasterDeleter struct {
	cfg          config.Config
	cloudStorage interfaces.CloudStorageUseCase
}

// MasterDeleterUseCase is a use case for the Master creator
type MasterDeleterUseCase interface {
}

// NewMasterDeleter creates a new MasterDeleter
func NewMasterDeleter(
	cfg config.Config,
	cloudStorage interfaces.CloudStorageUseCase,
) *MasterDeleter {
	return &MasterDeleter{
		cfg:          cfg,
		cloudStorage: cloudStorage,
	}
}
