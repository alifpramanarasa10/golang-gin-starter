package handler

import (
	"gin-starter/common/errors"
	"gin-starter/common/interfaces"
	"gin-starter/modules/user/v1/service"
	"gin-starter/resource"
	"gin-starter/response"
	"gin-starter/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

// UserCreatorHandler is a handler for user finder
type UserCreatorHandler struct {
	userCreator  service.UserCreatorUseCase
	cloudStorage interfaces.CloudStorageUseCase
}

// NewUserCreatorHandler is a constructor for UserCreatorHandler
func NewUserCreatorHandler(
	userCreator service.UserCreatorUseCase,
	cloudStorage interfaces.CloudStorageUseCase,
) *UserCreatorHandler {
	return &UserCreatorHandler{
		userCreator:  userCreator,
		cloudStorage: cloudStorage,
	}
}

// CreateUser is a handler for creating user
func (uc *UserCreatorHandler) CreateUser(c *gin.Context) {
	var request resource.CreateUserRequest

	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorAPIResponse(http.StatusBadRequest, err.Error()))
		c.Abort()
		return
	}

	imagePath, err := uc.cloudStorage.Upload(request.Photo, "users/user/profile")

	if err != nil {
		parseError := errors.ParseError(err)
		c.JSON(parseError.Code, response.ErrorAPIResponse(parseError.Code, parseError.Message))
		c.Abort()
		return
	}

	dob, err := utils.DateStringToTime(request.DOB)

	if err != nil {
		parseError := errors.ParseError(err)
		c.JSON(parseError.Code, response.ErrorAPIResponse(parseError.Code, parseError.Message))
		c.Abort()
		return
	}

	user, err := uc.userCreator.CreateUser(
		c,
		request.Name,
		request.Email,
		request.Password,
		request.PhoneNumber,
		imagePath,
		dob,
	)

	if err != nil {
		parseError := errors.ParseError(err)
		c.JSON(parseError.Code, response.ErrorAPIResponse(parseError.Code, parseError.Message))
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, response.SuccessAPIResponseList(http.StatusOK, "success", resource.NewUserProfile(user)))
}

// CreateAdmin is a handler for creating admin
func (uc *UserCreatorHandler) CreateAdmin(c *gin.Context) {
	var request resource.CreateAdminRequest

	if err := c.ShouldBind(&request); err != nil {
		parseError := errors.ParseError(err)
		c.JSON(parseError.Code, response.ErrorAPIResponse(parseError.Code, parseError.Message))
		c.Abort()
		return
	}

	imagePath, err := uc.cloudStorage.Upload(request.Photo, "users/admin/profile")

	if err != nil {
		parseError := errors.ParseError(err)
		c.JSON(parseError.Code, response.ErrorAPIResponse(parseError.Code, parseError.Message))
		c.Abort()
		return
	}

	dob, err := utils.DateStringToTime(request.DOB)

	if err != nil {
		parseError := errors.ParseError(err)
		c.JSON(parseError.Code, response.ErrorAPIResponse(parseError.Code, parseError.Message))
		c.Abort()
		return
	}

	roleID, err := uuid.Parse(request.RoleID)

	if err != nil {
		parseError := errors.ParseError(err)
		c.JSON(parseError.Code, response.ErrorAPIResponse(parseError.Code, parseError.Message))
		c.Abort()
		return
	}

	user, err := uc.userCreator.CreateAdmin(
		c,
		request.Name,
		request.Email,
		request.Password,
		request.PhoneNumber,
		imagePath,
		dob,
		roleID,
	)

	if err != nil {
		parseError := errors.ParseError(err)
		c.JSON(parseError.Code, response.ErrorAPIResponse(parseError.Code, parseError.Message))
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, response.SuccessAPIResponseList(http.StatusOK, "success", resource.NewUserProfile(user)))
}

// CreatePermission is a handler for creating permission data
func (uc *UserCreatorHandler) CreatePermission(c *gin.Context) {
	var request resource.CreatePermissionRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorAPIResponse(http.StatusBadRequest, err.Error()))
		c.Abort()
		return
	}

	permission, err := uc.userCreator.CreatePermission(c, request.Name, request.Label)
	if err != nil {
		parseError := errors.ParseError(err)
		c.JSON(parseError.Code, response.ErrorAPIResponse(parseError.Code, parseError.Message))
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, response.SuccessAPIResponseList(http.StatusOK, "success", resource.NewPermissionResponse(permission)))
}

// CreateRole is a handler for creating role data
func (uc *UserCreatorHandler) CreateRole(c *gin.Context) {
	var request resource.CreateRoleRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorAPIResponse(http.StatusBadRequest, err.Error()))
		c.Abort()
		return
	}

	var permissionIDs []uuid.UUID
	if len(request.PermissionIDs) > 0 {
		for _, permissionID := range request.PermissionIDs {
			valid, err := uuid.Parse(permissionID)
			if err != nil {
				parseError := errors.ParseError(err)
				c.JSON(parseError.Code, response.ErrorAPIResponse(parseError.Code, parseError.Message))
				c.Abort()
				return
			}
			permissionIDs = append(permissionIDs, valid)
		}
	}

	role, err := uc.userCreator.CreateRole(
		c,
		request.Name,
		permissionIDs,
		"system",
	)
	if err != nil {
		parseError := errors.ParseError(err)
		c.JSON(parseError.Code, response.ErrorAPIResponse(parseError.Code, parseError.Message))
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, response.SuccessAPIResponseList(http.StatusOK, "success", resource.NewRoleResponse(role)))
}
