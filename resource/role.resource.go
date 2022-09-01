package resource

import (
	"gin-starter/entity"

	"github.com/google/uuid"
)

// GetRoleResponse returns a role response
type GetRoleResponse struct {
	List  []*Role `json:"list"`
	Total int64   `json:"total"`
	Meta  *Meta   `json:"meta"`
}

// GetRoleByID returns a role by id request
type GetRoleByID struct {
	ID string `uri:"id" binding:"required"`
}

// CreateRoleRequest is a request for create role
type CreateRoleRequest struct {
	Name          string   `form:"name" json:"name" binding:"required"`
	PermissionIDs []string `form:"permission_id"`
}

// UpdateRoleRequest is a request for update role
type UpdateRoleRequest struct {
	Name          string   `form:"name" json:"name" binding:"required"`
	PermissionIDs []string `form:"permission_id"`
}

// DeleteRoleRequest is a request for delete role
type DeleteRoleRequest struct {
	ID string `uri:"id" binding:"required"`
}

// Role is a base response for role
type Role struct {
	ID         uuid.UUID     `json:"id"`
	Name       string        `json:"name"`
	Permission []*Permission `json:"permissions"`
}

// NewRoleResponse returns a role response
func NewRoleResponse(role *entity.Role) *Role {
	if role == nil {
		return nil
	}

	permissions := make([]*Permission, 0)

	if len(role.RolePermissions) > 0 {
		for _, v := range role.RolePermissions {
			if v.Permission != nil {
				permissions = append(permissions, &Permission{
					ID:    v.Permission.ID,
					Name:  v.Permission.Name,
					Label: v.Permission.Label,
				})
			}
		}
	}

	return &Role{
		ID:         role.ID,
		Name:       role.Name,
		Permission: permissions,
	}
}
