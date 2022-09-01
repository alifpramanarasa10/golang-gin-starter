package entity

import (
	"time"

	"github.com/google/uuid"
)

const (
	permissionTableName = "main.permissions"
)

// Permission defines table permission
type Permission struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Label string    `json:"label"`
	Auditable
}

// TableName specifies table name
func (model *Permission) TableName() string {
	return permissionTableName
}

// NewPermission creating new permission entity
func NewPermission(
	id uuid.UUID,
	name string,
	label string,
	createdBy string,
) *Permission {
	return &Permission{
		ID:        id,
		Name:      name,
		Label:     label,
		Auditable: NewAuditable(createdBy),
	}
}

// MapUpdateFrom mapping from model
func (model *Permission) MapUpdateFrom(from *Permission) *map[string]interface{} {
	if from == nil {
		return &map[string]interface{}{
			"name":       model.Name,
			"created_by": model.CreatedBy,
			"updated_at": model.UpdatedAt,
		}
	}

	mapped := make(map[string]interface{})

	if model.Name != from.Name {
		mapped["name"] = from.Name
	}

	if model.Label != from.Label {
		mapped["label"] = from.Label
	}

	mapped["updated_at"] = time.Now()
	return &mapped
}
