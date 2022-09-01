package handler

import (
	"gin-starter/common/errors"
	"gin-starter/middleware"
	"gin-starter/modules/user/v1/service"
	"gin-starter/resource"
	"gin-starter/response"
	"github.com/google/uuid"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserFinderHandler is a handler for user finder
type UserFinderHandler struct {
	userFinder service.UserFinderUseCase
}

// NewUserFinderHandler is a constructor for UserFinderHandler
func NewUserFinderHandler(
	userFinder service.UserFinderUseCase,
) *UserFinderHandler {
	return &UserFinderHandler{
		userFinder: userFinder,
	}
}

// GetUserProfile is a handler for getting user profile
func (uf *UserFinderHandler) GetUserProfile(c *gin.Context) {
	res, err := uf.userFinder.GetUserByID(c, middleware.UserID)

	if err != nil {
		parseError := errors.ParseError(err)
		c.JSON(parseError.Code, response.ErrorAPIResponse(parseError.Code, parseError.Message))
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, response.SuccessAPIResponseList(http.StatusOK, "success", resource.NewUserProfile(res)))
}

// GetUserByForgotPasswordToken is a handler for getting user by forgot password token
func (uf *UserFinderHandler) GetUserByForgotPasswordToken(c *gin.Context) {
	var request resource.GetUserByForgotPasswordTokenRequest

	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorAPIResponse(http.StatusBadRequest, err.Error()))
		c.Abort()
		return
	}

	res, err := uf.userFinder.GetUserByForgotPasswordToken(c, request.Token)

	if err != nil {
		parseError := errors.ParseError(err)
		c.JSON(parseError.Code, response.ErrorAPIResponse(parseError.Code, parseError.Message))
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, response.SuccessAPIResponseList(http.StatusOK, "success", resource.NewUserProfile(res)))
}

// GetAdminProfile is a handler for getting admin profile
func (uf *UserFinderHandler) GetAdminProfile(c *gin.Context) {
	user, err := uf.userFinder.GetUserByID(c, middleware.UserID)

	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorAPIResponse(http.StatusBadRequest, err.Error()))
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, response.SuccessAPIResponseList(http.StatusOK, "success", resource.NewUserAdmin(user)))
}

// GetUsers is a handler for getting users
func (uf *UserFinderHandler) GetUsers(c *gin.Context) {
	var request resource.GetAdminUsersRequest

	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorAPIResponse(http.StatusBadRequest, err.Error()))
		c.Abort()
		return
	}

	users, total, err := uf.userFinder.GetUsers(c, request.Query, request.Sort, request.Order, request.Limit, request.Offset)

	if err != nil {
		parseError := errors.ParseError(err)
		c.JSON(parseError.Code, response.ErrorAPIResponse(parseError.Code, parseError.Message))
		c.Abort()
		return
	}

	res := make([]*resource.UserProfile, 0)

	for _, u := range users {
		res = append(res, resource.NewUserProfile(u))
	}

	c.JSON(http.StatusOK, response.SuccessAPIResponseList(http.StatusOK, "success", &resource.GetUsersResponse{
		List:  res,
		Total: total,
	}))
}

// GetUserByID is a handler for getting user by ID
func (uf *UserFinderHandler) GetUserByID(c *gin.Context) {
	var request resource.GetUserByIDRequest

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

	user, err := uf.userFinder.GetUserByID(c, reqID)

	if err != nil {
		parseError := errors.ParseError(err)
		c.JSON(parseError.Code, response.ErrorAPIResponse(parseError.Code, parseError.Message))
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, response.SuccessAPIResponseList(http.StatusOK, "success", resource.NewUserProfile(user)))
}

// GetAdminUsers is a handler for getting admin users
func (uf *UserFinderHandler) GetAdminUsers(c *gin.Context) {
	var request resource.GetAdminUsersRequest

	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorAPIResponse(http.StatusBadRequest, err.Error()))
		c.Abort()
		return
	}

	users, total, err := uf.userFinder.GetAdminUsers(c, request.Query, request.Sort, request.Order, request.Limit, request.Offset)

	if err != nil {
		parseError := errors.ParseError(err)
		c.JSON(parseError.Code, response.ErrorAPIResponse(parseError.Code, parseError.Message))
		c.Abort()
		return
	}

	res := make([]*resource.UserAdmin, 0)

	for _, u := range users {
		res = append(res, resource.NewUserAdmin(u))
	}

	c.JSON(http.StatusOK, response.SuccessAPIResponseList(http.StatusOK, "success", &resource.GetAdminUsersResponse{
		List:  res,
		Total: total,
	}))
}

// GetAdminUserByID is a handler for getting admin user by ID
func (uf *UserFinderHandler) GetAdminUserByID(c *gin.Context) {
	var request resource.GetAdminUserByIDRequest

	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorAPIResponse(http.StatusBadRequest, err.Error()))
		c.Abort()
		return
	}

	userID, err := uuid.Parse(request.ID)

	if err != nil {
		c.JSON(errors.ErrInvalidArgument.Code, response.ErrorAPIResponse(errors.ErrInvalidArgument.Code, errors.ErrInvalidArgument.Message))
		c.Abort()
		return
	}

	user, err := uf.userFinder.GetUserByID(c, userID)

	if err != nil {
		parseError := errors.ParseError(err)
		c.JSON(parseError.Code, response.ErrorAPIResponse(parseError.Code, parseError.Message))
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, response.SuccessAPIResponseList(http.StatusOK, "success", resource.NewUserAdmin(user)))
}

// GetRoles is a handler for getting roles
func (uf *UserFinderHandler) GetRoles(c *gin.Context) {
	var request resource.PaginationQueryParam
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorAPIResponse(http.StatusBadRequest, err.Error()))
		c.Abort()
		return
	}

	page, err := uf.userFinder.GetRoles(c, request.Query, request.Sort, request.Order, request.Limit, request.Offset)

	if err != nil {
		parseError := errors.ParseError(err)
		c.JSON(parseError.Code, response.ErrorAPIResponse(parseError.Code, parseError.Message))
		c.Abort()
		return
	}

	res := make([]*resource.Role, 0)

	for _, v := range page {
		res = append(res, resource.NewRoleResponse(v))
	}

	currentPage := request.Offset/request.Limit + 1
	totalPage := len(res) / request.Limit
	if len(res)%request.Limit > 0 {
		totalPage++
	}

	meta := &resource.Meta{
		Total:       len(res),
		Limit:       request.Limit,
		Offset:      request.Offset,
		CurrentPage: currentPage,
		TotalPage:   totalPage,
	}

	c.JSON(http.StatusOK, response.SuccessAPIResponseList(http.StatusOK, "success", &resource.GetRoleResponse{
		List:  res,
		Total: int64(len(res)),
		Meta:  meta,
	}))
}

// GetPermissions is a handler for getting permission
func (uf *UserFinderHandler) GetPermissions(c *gin.Context) {
	permissions, err := uf.userFinder.GetPermissions(c)

	if err != nil {
		parseError := errors.ParseError(err)
		c.JSON(parseError.Code, response.ErrorAPIResponse(parseError.Code, parseError.Message))
		c.Abort()
		return
	}

	res := make([]*resource.Permission, 0)

	for _, p := range permissions {
		res = append(res, resource.NewPermissionResponse(p))
	}

	c.JSON(http.StatusOK, response.SuccessAPIResponseList(http.StatusOK, "success", &resource.GetPermissionResponse{
		List:  res,
		Total: int64(len(res)),
	}))
}

// GetUserPermissions is a handler for get list permission of user
func (uf *UserFinderHandler) GetUserPermissions(c *gin.Context) {
	adminPermission, err := uf.userFinder.GetUserPermissions(c, middleware.UserID)
	if err != nil {
		parseError := errors.ParseError(err)
		c.JSON(parseError.Code, response.ErrorAPIResponse(parseError.Code, parseError.Message))
		c.Abort()
		return
	}

	res := make([]*resource.Permission, 0)
	for _, ap := range adminPermission {
		res = append(res, resource.NewPermissionResponse(ap))
	}

	c.JSON(http.StatusOK, response.SuccessAPIResponseList(http.StatusOK, "success", &resource.GetPermissionResponse{
		List:  res,
		Total: int64(len(res)),
	}))
}
