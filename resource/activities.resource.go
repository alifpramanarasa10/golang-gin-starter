package resource

import (
	"gin-starter/entity"

	"github.com/google/uuid"
)

type CreateActivitiesRequest struct {
	ID             uuid.UUID `form:"id" json:"id"`
	UserID         uuid.UUID `form:"user_id" json:"user_id"`
	Title          string    `form:"title" json:"title"`
	Description    string    `form:"description" json:"description"`
	ActivitiesType string    `form:"activities_type" json:"activities_type"`
}

type UpdateActivitiesRequest struct {
	ID             uuid.UUID `form:"id" json:"id"`
	UserID         uuid.UUID `form:"user_id" json:"user_id"`
	Title          string    `form:"title" json:"title"`
	Description    string    `form:"description" json:"description"`
	ActivitiesType string    `form:"activities_type" json:"activities_type"`
}

type DeleteActivitiesRequest struct {
	ID uuid.UUID `uri:"id" binding:"required"`
}

type Activities struct {
	ID             uuid.UUID `json:"id"`
	UserID         uuid.UUID `json:"user_id"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	ActivitiesType string    `json:"activities_type"`
}

type GetActivitiesRequest struct {
	Query  string `form:"query" json:"query"`
	Sort   string `form:"sort" json:"sort"`
	Order  string `form:"order" json:"order"`
	Limit  int    `form:"limit,default=10" json:"limit"`
	Offset int    `form:"offset,default=0" json:"offset"`
}

type GetActivitiesResponse struct {
	List  []*Activities `json:"list"`
	Total int64         `json:"total"`
}

type GetActivitiesWithoutTotalResponse struct {
	List []*Activities `json:"list"`
}

type GetActivitiesByIDRequest struct {
	ID uuid.UUID `uri:"id" binding:"required"`
}

type GetActivitiesByUserIDRequest struct {
	UserID uuid.UUID `uri:"id" binding:"required"`
}

func NewActivities(activities *entity.Activities) *Activities {
	if activities == nil {
		return nil
	}

	return &Activities{
		ID:             activities.ID,
		UserID:         activities.UserID,
		Title:          activities.Title,
		Description:    activities.Description,
		ActivitiesType: activities.ActivitiesType,
	}
}
