package handler

import (
	"gin-starter/common/errors"
	"gin-starter/entity"
	"gin-starter/modules/activities/v1/service"
	"gin-starter/resource"
	"gin-starter/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ActivitiesUpdaterHandler is handler for activities updater
type ActivitiesUpdaterHandler struct {
	activitiesUpdater service.ActivitiesUpdaterUseCase
}

func NewActivitiesUpdaterHandler(
	activitiesUpdater service.ActivitiesUpdaterUseCase,
) *ActivitiesUpdaterHandler {
	return &ActivitiesUpdaterHandler{
		activitiesUpdater: activitiesUpdater,
	}
}

// UpdateActivities is a handler for updating activities
func (a *ActivitiesUpdaterHandler) UpdateActivities(c *gin.Context) {
	var request resource.UpdateActivitiesRequest

	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorAPIResponse(http.StatusBadRequest, err.Error()))
		c.Abort()
		return
	}

	activities := entity.NewActivities(
		request.ID,
		request.UserID,
		request.Title,
		request.Description,
		request.ActivitiesType,
		"system",
	)

	if err := a.activitiesUpdater.UpdateActivity(c, activities); err != nil {
		parseError := errors.ParseError(err)
		c.JSON(parseError.Code, response.ErrorAPIResponse(parseError.Code, parseError.Message))
		c.Abort()
		return
	}
}
