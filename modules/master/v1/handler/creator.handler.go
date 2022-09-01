package handler

import (
	"gin-starter/common/interfaces"
	"gin-starter/modules/master/v1/service"
)

// MasterCreatorHandler is a handler for master finder
type MasterCreatorHandler struct {
	masterCreator service.MasterCreatorUseCase
	cloudStorage  interfaces.CloudStorageUseCase
}

// NewMasterCreatorHandler is a constructor for MasterCreatorHandler
func NewMasterCreatorHandler(
	masterCreator service.MasterCreatorUseCase,
	cloudStorage interfaces.CloudStorageUseCase,
) *MasterCreatorHandler {
	return &MasterCreatorHandler{
		masterCreator: masterCreator,
		cloudStorage:  cloudStorage,
	}
}
