package entity

import (
	"time"

	"github.com/google/uuid"
)

const (
	activitiesTableName = "activities.activities"
)

type Activities struct {
	ID             uuid.UUID `json:"id"`
	UserID         uuid.UUID `json:"user_id"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	ActivitiesType string    `json:"activities_type"`
	Auditable
}

// TableName specifies table name
func (model *Activities) TableName() string {
	return activitiesTableName
}

func NewActivities(
	id uuid.UUID,
	userID uuid.UUID,
	title string,
	description string,
	activities_type string,
	createdBy string,
) *Activities {
	return &Activities{
		ID:             id,
		UserID:         userID,
		Title:          title,
		Description:    description,
		ActivitiesType: activities_type,
		Auditable:      NewAuditable(createdBy),
	}
}

// MapUpdateFrom mapping from model
func (model *Activities) MapUpdateFrom(from *Activities) *map[string]interface{} {
	if from == nil {
		return &map[string]interface{}{
			"user_id":         model.UserID,
			"title":           model.Title,
			"description":     model.Description,
			"activities_type": model.ActivitiesType,
			"created_by":      model.CreatedBy,
			"updated_at":      model.UpdatedAt,
		}
	}

	mapped := make(map[string]interface{})

	if model.UserID != from.UserID {
		mapped["user_id"] = from.UserID
	}

	if model.Title != from.Title {
		mapped["title"] = from.Title
	}

	if model.Description != from.Description {
		mapped["description"] = from.Description
	}

	if model.ActivitiesType != from.ActivitiesType {
		mapped["activities_type"] = from.ActivitiesType
	}

	mapped["updated_at"] = time.Now()
	return &mapped
}
