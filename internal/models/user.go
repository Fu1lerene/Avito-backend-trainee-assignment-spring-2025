package models

import "github.com/google/uuid"

type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	Role         UserRole
}

type UserRole string

const (
	UserRoleEmployee  UserRole = "employee"
	UserRoleModerator UserRole = "moderator"
)
