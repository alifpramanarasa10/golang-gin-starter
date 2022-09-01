package interfaces

import (
	"context"
	"gin-starter/entity"

	"github.com/google/uuid"
)

// MiddlewareTool is a service for middleware tool
type MiddlewareTool interface {
	RequiredSignedIn(ctx context.Context, userID uuid.UUID, permission string) (*entity.UserRole, error)
	RequiredPermission(ctx context.Context, roleID uuid.UUID, permission string) (*entity.Permission, error)
}
