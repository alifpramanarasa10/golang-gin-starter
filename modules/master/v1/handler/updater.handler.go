package handler

import (
	"gin-starter/common/interfaces"
	"gin-starter/modules/master/v1/service"
)

// MasterUpdaterHandler is a handler for master finder
type MasterUpdaterHandler struct {
	masterUpdater service.MasterUpdaterUseCase
	cloudStorage  interfaces.CloudStorageUseCase
}

// NewMasterUpdaterHandler is a constructor for MasterUpdaterHandler
func NewMasterUpdaterHandler(
	masterUpdater service.MasterUpdaterUseCase,
	cloudStorage interfaces.CloudStorageUseCase,
) *MasterUpdaterHandler {
	return &MasterUpdaterHandler{
		masterUpdater: masterUpdater,
		cloudStorage:  cloudStorage,
	}
}
