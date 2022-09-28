package handler

import (
	"gin-starter/common/errors"
	"gin-starter/modules/activities/v1/service"
	"gin-starter/resource"
	"gin-starter/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ActivitiesCreatorHTTPHandler is handler for activities creator
type ActivitiesCreatorHandler struct {
	activitiesCreator service.ActivitiesCreatorUseCase
}

func NewActivitiesCreatorHandler(
	activitiesCreator service.ActivitiesCreatorUseCase,
) *ActivitiesCreatorHandler {
	return &ActivitiesCreatorHandler{
		activitiesCreator: activitiesCreator,
	}
}

// CreateActivities is a handler for creating activities
func (a *ActivitiesCreatorHandler) CreateActivities(c *gin.Context) {
	var request resource.CreateActivitiesRequest

	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorAPIResponse(http.StatusBadRequest, err.Error()))
		c.Abort()
		return
	}

	activities, err := a.activitiesCreator.CreateActivities(
		c.Request.Context(),
		request.UserID,
		request.Title,
		request.Description,
		request.ActivitiesType,
	)

	if err != nil {
		parseError := errors.ParseError(err)
		c.JSON(parseError.Code, response.ErrorAPIResponse(parseError.Code, parseError.Message))
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, response.SuccessAPIResponseList(http.StatusOK, "success", resource.NewActivities(activities)))
}
