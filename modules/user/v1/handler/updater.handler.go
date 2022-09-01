package handler

import (
	"gin-starter/common/errors"
	"gin-starter/common/interfaces"
	"gin-starter/entity"
	"gin-starter/middleware"
	"gin-starter/modules/user/v1/service"
	"gin-starter/resource"
	"gin-starter/response"
	"gin-starter/utils"
	"github.com/google/uuid"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserUpdaterHandler is a handler for user updater
type UserUpdaterHandler struct {
	userUpdater  service.UserUpdaterUseCase
	userFinder   service.UserFinderUseCase
	cloudStorage interfaces.CloudStorageUseCase
}

// NewUserUpdaterHandler is a constructor for UserUpdaterHandler
func NewUserUpdaterHandler(
	userUpdater service.UserUpdaterUseCase,
	userFinder service.UserFinderUseCase,
	cloudStorage interfaces.CloudStorageUseCase,
) *UserUpdaterHandler {
	return &UserUpdaterHandler{
		userUpdater:  userUpdater,
		userFinder:   userFinder,
		cloudStorage: cloudStorage,
	}
}

// ChangePassword is a handler for changing password
func (uu *UserUpdaterHandler) ChangePassword(c *gin.Context) {
	var request resource.ChangePasswordRequest

	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorAPIResponse(http.StatusBadRequest, err.Error()))
		c.Abort()
		return
	}

	if request.NewPassword != request.NewPasswordConfirmation {
		c.JSON(http.StatusBadRequest, response.ErrorAPIResponse(http.StatusBadRequest, errors.ErrWrongPasswordConfirmation.Message))
		c.Abort()
		return
	}

	if err := uu.userUpdater.ChangePassword(
		c,
		middleware.UserID,
		request.OldPassword,
		request.NewPassword,
	); err != nil {
		parseError := errors.ParseError(err)
		c.JSON(parseError.Code, response.ErrorAPIResponse(parseError.Code, parseError.Message))
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, response.SuccessAPIResponseList(http.StatusOK, "success", nil))
}

// ForgotPasswordRequest is a handler for forgot password request
func (uu *UserUpdaterHandler) ForgotPasswordRequest(c *gin.Context) {
	var request resource.ForgotPasswordRequest

	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorAPIResponse(http.StatusBadRequest, err.Error()))
		c.Abort()
		return
	}

	if err := uu.userUpdater.ForgotPasswordRequest(
		c,
		request.Email,
	); err != nil {
		parseError := errors.ParseError(err)
		c.JSON(parseError.Code, response.ErrorAPIResponse(parseError.Code, parseError.Message))
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, response.SuccessAPIResponseList(http.StatusOK, "success", nil))
}

// VerifyOTP is a handler for verifying OTP
func (uu *UserUpdaterHandler) VerifyOTP(c *gin.Context) {
	var request resource.VerifyOTPRequest

	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorAPIResponse(http.StatusBadRequest, err.Error()))
		c.Abort()
		return
	}

	verify, err := uu.userUpdater.VerifyOTP(
		c,
		middleware.UserID,
		request.Code,
	)

	if err != nil {
		parseError := errors.ParseError(err)
		c.JSON(parseError.Code, response.ErrorAPIResponse(parseError.Code, parseError.Message))
		c.Abort()
		return
	}

	if !verify {
		c.JSON(http.StatusBadRequest, response.ErrorAPIResponse(http.StatusBadRequest, errors.ErrOTPMismatch.Message))
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, response.SuccessAPIResponseList(http.StatusOK, "success", nil))
}

// ResendOTP is a handler for resending OTP
func (uu *UserUpdaterHandler) ResendOTP(c *gin.Context) {
	if err := uu.userUpdater.ResendOTP(
		c,
		middleware.UserID,
	); err != nil {
		parseError := errors.ParseError(err)
		c.JSON(parseError.Code, response.ErrorAPIResponse(parseError.Code, parseError.Message))
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, response.SuccessAPIResponseList(http.StatusOK, "success", nil))
}

// ForgotPassword is a handler for forgot password
func (uu *UserUpdaterHandler) ForgotPassword(c *gin.Context) {
	var request resource.ForgotPasswordChangeRequest

	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorAPIResponse(http.StatusBadRequest, err.Error()))
		c.Abort()
		return
	}

	if request.NewPassword != request.NewPasswordConfirmation {
		c.JSON(http.StatusBadRequest, response.ErrorAPIResponse(http.StatusBadRequest, errors.ErrWrongPasswordConfirmation.Message))
		c.Abort()
		return
	}

	res, err := uu.userFinder.GetUserByForgotPasswordToken(c, request.Token)

	if err != nil {
		parseError := errors.ParseError(err)
		c.JSON(parseError.Code, response.ErrorAPIResponse(parseError.Code, parseError.Message))
		c.Abort()
		return
	}

	if err := uu.userUpdater.ForgotPassword(
		c,
		res.ID,
		request.NewPassword,
	); err != nil {
		parseError := errors.ParseError(err)
		c.JSON(parseError.Code, response.ErrorAPIResponse(parseError.Code, parseError.Message))
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, response.SuccessAPIResponseList(http.StatusOK, "success", nil))
}

// UpdateUser is a handler for updating user
func (uu *UserUpdaterHandler) UpdateUser(c *gin.Context) {
	var request resource.UpdateUserRequest

	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorAPIResponse(http.StatusBadRequest, err.Error()))
		c.Abort()
		return
	}

	_, err := uu.userFinder.GetUserByID(c, middleware.UserID)

	if err != nil {
		parseError := errors.ParseError(err)
		c.JSON(parseError.Code, response.ErrorAPIResponse(parseError.Code, parseError.Message))
		c.Abort()
		return
	}

	imagePath, err := uu.cloudStorage.Upload(request.Photo, "users/user/profile")

	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorAPIResponse(http.StatusBadRequest, err.Error()))
		c.Abort()
		return
	}

	dob, err := utils.DateStringToTime(request.DOB)

	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorAPIResponse(http.StatusBadRequest, err.Error()))
		c.Abort()
		return
	}

	user := entity.NewUser(
		middleware.UserID,
		request.Name,
		request.Email,
		request.Name,
		utils.TimeToNullTime(dob),
		imagePath,
		request.PhoneNumber,
		"system",
	)

	if err := uu.userUpdater.Update(c, user); err != nil {
		parseError := errors.ParseError(err)
		c.JSON(parseError.Code, response.ErrorAPIResponse(parseError.Code, parseError.Message))
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, response.SuccessAPIResponseList(http.StatusOK, "success", nil))
}

// ActivateDeactivateUser is a handler for activating and deactivating user
func (uu *UserUpdaterHandler) ActivateDeactivateUser(c *gin.Context) {
	var request resource.DeactivateUserRequest

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

	if err := uu.userUpdater.ActivateDeactivateUser(c, reqID); err != nil {
		parseError := errors.ParseError(err)
		c.JSON(parseError.Code, response.ErrorAPIResponse(parseError.Code, parseError.Message))
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, response.SuccessAPIResponseList(http.StatusOK, "success", nil))
}

// UpdateAdmin is a handler for updating admin
func (uu *UserUpdaterHandler) UpdateAdmin(c *gin.Context) {
	var request resource.UpdateAdminRequest

	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorAPIResponse(http.StatusBadRequest, err.Error()))
		c.Abort()
		return
	}

	userID, err := uuid.Parse(c.Param("id"))

	if err != nil {
		c.JSON(errors.ErrInvalidArgument.Code, response.ErrorAPIResponse(errors.ErrInvalidArgument.Code, errors.ErrInvalidArgument.Message))
		c.Abort()
		return
	}

	_, err = uu.userFinder.GetUserByID(c, userID)

	if err != nil {
		parseError := errors.ParseError(err)
		c.JSON(parseError.Code, response.ErrorAPIResponse(parseError.Code, parseError.Message))
		c.Abort()
		return
	}

	imagePath, err := uu.cloudStorage.Upload(request.Photo, "users/admin/profile")

	if err != nil {
		c.JSON(errors.ErrInternalServerError.Code, response.ErrorAPIResponse(errors.ErrInternalServerError.Code, errors.ErrInternalServerError.Message))
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

	user := entity.NewUser(
		userID,
		request.Name,
		request.Email,
		request.Name,
		utils.TimeToNullTime(dob),
		imagePath,
		request.PhoneNumber,
		"system",
	)

	roleID, err := uuid.Parse(request.RoleID)

	if err != nil {
		c.JSON(errors.ErrInvalidArgument.Code, response.ErrorAPIResponse(errors.ErrInvalidArgument.Code, errors.ErrInvalidArgument.Message))
		c.Abort()
		return
	}

	if err := uu.userUpdater.UpdateAdmin(c, user, roleID); err != nil {
		parseError := errors.ParseError(err)
		c.JSON(parseError.Code, response.ErrorAPIResponse(parseError.Code, parseError.Message))
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, response.SuccessAPIResponseList(http.StatusOK, "success", nil))
}

// UpdateRole is a handler for updating role
func (uu *UserUpdaterHandler) UpdateRole(c *gin.Context) {
	var request resource.UpdateRoleRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorAPIResponse(http.StatusBadRequest, err.Error()))
		c.Abort()
		return
	}

	roleID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(errors.ErrInvalidArgument.Code, response.ErrorAPIResponse(errors.ErrInvalidArgument.Code, errors.ErrInvalidArgument.Message))
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

	if err := uu.userUpdater.UpdateRole(c, roleID, request.Name, permissionIDs); err != nil {
		parseError := errors.ParseError(err)
		c.JSON(parseError.Code, response.ErrorAPIResponse(parseError.Code, parseError.Message))
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, response.SuccessAPIResponseList(http.StatusOK, "success", nil))
}

// UpdatePermission is a handler for updating permission
func (uu *UserUpdaterHandler) UpdatePermission(c *gin.Context) {
	var request resource.UpdatePermissionRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorAPIResponse(http.StatusBadRequest, err.Error()))
		c.Abort()
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(errors.ErrInvalidArgument.Code, response.ErrorAPIResponse(errors.ErrInvalidArgument.Code, errors.ErrInvalidArgument.Message))
		c.Abort()
		return
	}

	if err := uu.userUpdater.UpdatePermission(c, id, request.Name, request.Label); err != nil {
		parseError := errors.ParseError(err)
		c.JSON(parseError.Code, response.ErrorAPIResponse(parseError.Code, parseError.Message))
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, response.SuccessAPIResponseList(http.StatusOK, "success", nil))
}
