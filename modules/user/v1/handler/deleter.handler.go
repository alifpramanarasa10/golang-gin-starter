package handler

import (
	"gin-starter/common/errors"
	"gin-starter/common/interfaces"
	"gin-starter/modules/user/v1/service"
	"gin-starter/resource"
	"gin-starter/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

// UserDeleterHandler is a handler for user finder
type UserDeleterHandler struct {
	userDeleter  service.UserDeleterUseCase
	cloudStorage interfaces.CloudStorageUseCase
}

// NewUserDeleterHandler is a constructor for UserDeleterHandler
func NewUserDeleterHandler(
	userDeleter service.UserDeleterUseCase,
	cloudStorage interfaces.CloudStorageUseCase,
) *UserDeleterHandler {
	return &UserDeleterHandler{
		userDeleter:  userDeleter,
		cloudStorage: cloudStorage,
	}
}

// DeleteAdmin is a handler for deleting an admin
func (ud *UserDeleterHandler) DeleteAdmin(c *gin.Context) {
	var request resource.DeleteAdminRequest

	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorAPIResponse(http.StatusBadRequest, err.Error()))
		c.Abort()
		return
	}

	reqID, err := uuid.Parse(request.ID)

	if err != nil {
		c.JSON(errors.ErrInvalidArgument.Code, response.ErrorAPIResponse(errors.ErrInvalidArgument.Code, errors.ErrInvalidArgument.Message))
		c.Abort()
		return
	}

	if err := ud.userDeleter.DeleteAdmin(c, reqID); err != nil {
		parseError := errors.ParseError(err)
		c.JSON(parseError.Code, response.ErrorAPIResponse(parseError.Code, parseError.Message))
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, response.SuccessAPIResponseList(http.StatusOK, "success", nil))
}

// DeleteRole is a handler for delete role
func (ud *UserDeleterHandler) DeleteRole(c *gin.Context) {
	var request resource.DeleteRoleRequest

	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorAPIResponse(http.StatusBadRequest, err.Error()))
		c.Abort()
		return
	}

	reqID, err := uuid.Parse(request.ID)

	if err != nil {
		c.JSON(errors.ErrInvalidArgument.Code, response.ErrorAPIResponse(errors.ErrInvalidArgument.Code, errors.ErrInvalidArgument.Message))
		c.Abort()
		return
	}

	if err := ud.userDeleter.DeleteRole(c, reqID, "system"); err != nil {
		parseError := errors.ParseError(err)
		c.JSON(parseError.Code, response.ErrorAPIResponse(parseError.Code, parseError.Message))
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, response.SuccessAPIResponseList(http.StatusOK, "success", nil))
}
