package resource

import (
	"gin-starter/entity"

	"github.com/google/uuid"
)

// CreatePermissionRequest returns a create permission request
type CreatePermissionRequest struct {
	Name  string `form:"name" json:"name" binding:"required"`
	Label string `form:"label" json:"label" binding:"required"`
}

// UpdatePermissionRequest returns a update permission request
type UpdatePermissionRequest struct {
	ID    string `form:"id" json:"id"`
	Name  string `form:"name" json:"name"`
	Label string `form:"label" json:"label"`
}

// Permission is a base response for permission
type Permission struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Label string    `json:"label"`
}

// NewPermissionResponse returns a permission response
func NewPermissionResponse(permission *entity.Permission) *Permission {
	return &Permission{
		ID:    permission.ID,
		Name:  permission.Name,
		Label: permission.Label,
	}
}

// GetPermissionResponse returns a list of permission
type GetPermissionResponse struct {
	List  []*Permission `json:"list"`
	Total int64         `json:"total"`
	Meta  *Meta         `json:"meta"`
}
