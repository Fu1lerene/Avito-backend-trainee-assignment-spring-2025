package roles

import (
	"avito/internal/models"
	"context"
)

type ctxKey string

const RoleKey ctxKey = "role"

func WithRole(ctx context.Context, role models.UserRole) context.Context {
	return context.WithValue(ctx, RoleKey, role)
}

func RoleFromContext(ctx context.Context) (models.UserRole, bool) {
	role, ok := ctx.Value(RoleKey).(models.UserRole)
	return role, ok
}

var allowedRoles = map[models.UserRole]bool{
	models.UserRoleModerator: true,
	models.UserRoleEmployee:  true,
}

func IsRoleAllowed(role models.UserRole) bool {
	return allowedRoles[role]
}
