package auth

import (
	"avito/internal/models"
	"avito/internal/repositories/user"
	"avito/internal/utils/errors"
	"avito/internal/utils/jwt"
	"avito/internal/utils/roles"
	"context"
	"golang.org/x/crypto/bcrypt"
	"net/mail"
)

type Service interface {
	DummyLogin(role models.UserRole) (string, error)
	Register(ctx context.Context, role models.UserRole, password, email string) (bool, error)
	Login(ctx context.Context, password, email string) (string, error)
}

type service struct {
	userRepo user.Repository
}

func New(userRepo user.Repository) Service {
	return &service{userRepo: userRepo}
}

func (s *service) DummyLogin(role models.UserRole) (string, error) {
	if !roles.IsRoleAllowed(role) {
		return "", errors.ErrInvalidRole
	}

	token, err := jwt.Generate(role)
	if err != nil {
		return "", errors.ErrInvalidRequest
	}

	return token, err
}

func (s *service) Register(ctx context.Context, role models.UserRole, password, email string) (bool, error) {
	if !roles.IsRoleAllowed(role) {
		return false, errors.ErrInvalidRole
	}
	_, err := mail.ParseAddress(email)
	if err != nil {
		return false, errors.ErrInvalidRequest
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return false, errors.ErrInvalidRequest
	}

	_, err = s.userRepo.Create(ctx, email, string(hashed), role)
	if err != nil {
		return false, errors.ErrInvalidRequest
	}

	return true, nil
}

func (s *service) Login(ctx context.Context, password, email string) (string, error) {
	user, err := s.userRepo.FindUserByEmail(ctx, email)
	if err != nil {
		return "", errors.ErrInvalidCredentials
	}

	_, err = mail.ParseAddress(email)
	if err != nil {
		return "", errors.ErrInvalidRequest
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.ErrInvalidCredentials
	}

	token, err := jwt.Generate(user.Role)
	if err != nil {
		return "", errors.ErrInvalidRequest
	}

	return token, nil
}
