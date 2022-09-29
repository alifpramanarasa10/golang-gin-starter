package handler

import (
	"gin-starter/common/errors"
	"gin-starter/modules/activities/v1/service"
	"gin-starter/resource"
	"gin-starter/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ActivitiesFinderHandler is handler for activities finder
type ActivitiesFinderHandler struct {
	activitiesFinder service.ActivitiesFinderUseCase
}

func NewActivitiesFinderHandler(
	activitiesFinder service.ActivitiesFinderUseCase,
) *ActivitiesFinderHandler {
	return &ActivitiesFinderHandler{
		activitiesFinder: activitiesFinder,
	}
}

// GetActivities is a handler for getting activities
func (a *ActivitiesFinderHandler) GetActivities(c *gin.Context) {
	var request resource.GetActivitiesRequest

	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorAPIResponse(http.StatusBadRequest, err.Error()))
		c.Abort()
		return
	}

	activities, total, err := a.activitiesFinder.GetActivities(c, request.Query, request.Sort, request.Order, request.Limit, request.Offset)

	if err != nil {
		parseError := errors.ParseError(err)
		c.JSON(parseError.Code, response.ErrorAPIResponse(parseError.Code, parseError.Message))
		c.Abort()
		return
	}

	res := make([]*resource.Activities, 0)

	for _, activity := range activities {
		res = append(res, resource.NewActivities(activity))
	}

	c.JSON(http.StatusOK, response.SuccessAPIResponseList(http.StatusOK, "success", &resource.GetActivitiesResponse{
		List:  res,
		Total: total,
	}))
}

// GetActivityByID is a handler for getting activity by ID
func (a *ActivitiesFinderHandler) GetActivityByID(c *gin.Context) {
	var request resource.GetActivitiesByIDRequest

	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorAPIResponse(http.StatusBadRequest, err.Error()))
		c.Abort()
		return
	}

	activity, err := a.activitiesFinder.GetActivitiesByID(c, request.ID)

	if err != nil {
		parseError := errors.ParseError(err)
		c.JSON(parseError.Code, response.ErrorAPIResponse(parseError.Code, parseError.Message))
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, response.SuccessAPIResponseList(http.StatusOK, "success", resource.NewActivities(activity)))
}

// GetActivityByUserID is a handler for getting activity by user ID
func (a *ActivitiesFinderHandler) GetActivityByUserID(c *gin.Context) {
	var request resource.GetActivitiesByUserIDRequest

	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorAPIResponse(http.StatusBadRequest, err.Error()))
		c.Abort()
		return
	}

	activities, err := a.activitiesFinder.GetActivitiesByUserID(c, request.UserID)

	if err != nil {
		parseError := errors.ParseError(err)
		c.JSON(parseError.Code, response.ErrorAPIResponse(parseError.Code, parseError.Message))
		c.Abort()
		return
	}

	res := make([]*resource.Activities, 0)

	for _, activity := range activities {
		res = append(res, resource.NewActivities(activity))
	}

	c.JSON(http.StatusOK, response.SuccessAPIResponseList(http.StatusOK, "success", &resource.GetActivitiesWithoutTotalResponse{
		List: res,
	}))
}
