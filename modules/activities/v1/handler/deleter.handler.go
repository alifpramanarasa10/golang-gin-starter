package handler

import (
	"gin-starter/common/errors"
	"gin-starter/modules/activities/v1/service"
	"gin-starter/resource"
	"gin-starter/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ActivitiesDeleterHandler is handler for activities deleter
type ActivitiesDeleterHandler struct {
	activitiesDeleter service.ActivitiesDeleterUseCase
}

func NewActivitiesDeleterHandler(
	activitiesDeleter service.ActivitiesDeleterUseCase,
) *ActivitiesDeleterHandler {
	return &ActivitiesDeleterHandler{
		activitiesDeleter: activitiesDeleter,
	}
}

// DeleteActivities is a handler for deleting activities
func (a *ActivitiesDeleterHandler) DeleteActivities(c *gin.Context) {
	var request resource.DeleteActivitiesRequest

	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorAPIResponse(http.StatusBadRequest, err.Error()))
		c.Abort()
		return
	}

	err := a.activitiesDeleter.DeleteActivity(
		c.Request.Context(),
		request.ID,
	)

	if err != nil {
		parseError := errors.ParseError(err)
		c.JSON(parseError.Code, response.ErrorAPIResponse(parseError.Code, parseError.Message))
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, response.SuccessAPIResponseList(http.StatusOK, "success", nil))
}
