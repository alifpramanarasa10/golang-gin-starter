package handler

import (
	"gin-starter/common/interfaces"
	"gin-starter/modules/master/v1/service"
)

// MasterDeleterHandler is a handler for master finder
type MasterDeleterHandler struct {
	masterDeleter service.MasterDeleterUseCase
	cloudStorage  interfaces.CloudStorageUseCase
}

// NewMasterDeleterHandler is a constructor for MasterDeleterHandler
func NewMasterDeleterHandler(
	masterDeleter service.MasterDeleterUseCase,
	cloudStorage interfaces.CloudStorageUseCase,
) *MasterDeleterHandler {
	return &MasterDeleterHandler{
		masterDeleter: masterDeleter,
		cloudStorage:  cloudStorage,
	}
}
