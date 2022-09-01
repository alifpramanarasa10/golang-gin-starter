package entity

import (
	"time"

	"github.com/google/uuid"
)

const (
	userRoleTableName = "main.user_roles"
)

// UserRole define for table user_roles
type UserRole struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
	RoleID uuid.UUID `json:"role_id"`
	Role   *Role     `foreignKey:"RoleID"`
	User   *User     `foreignKey:"UserID"`
	Auditable
}

// TableName specifies table name
func (model *UserRole) TableName() string {
	return userRoleTableName
}

// NewUserRole create new entity UserRole
func NewUserRole(
	id uuid.UUID,
	userID uuid.UUID,
	roleID uuid.UUID,
	createdBy string,
) *UserRole {
	return &UserRole{
		ID:        id,
		UserID:    userID,
		RoleID:    roleID,
		Auditable: NewAuditable(createdBy),
	}
}

// MapUpdateFrom mapping from model
func (model *UserRole) MapUpdateFrom(from *UserRole) *map[string]interface{} {
	if from == nil {
		return &map[string]interface{}{
			"role_id":    model.RoleID,
			"updated_at": model.UpdatedAt,
		}
	}

	mapped := make(map[string]interface{})

	if model.RoleID != from.RoleID {
		mapped["role_id"] = from.RoleID
	}

	mapped["updated_at"] = time.Now()
	return &mapped
}
