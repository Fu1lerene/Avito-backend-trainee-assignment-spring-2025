package auth

import (
	"avito/internal/models"
	"avito/internal/models/core_errors"
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) Create(ctx context.Context, email, hashedPassword string, role models.UserRole) error {
	args := m.Called(ctx, email, hashedPassword, role)
	if _, ok := args.Get(0).(bool); ok {
		return args.Error(1)
	}
	return args.Error(1)
}

func (m *MockUserRepo) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if rec, ok := args.Get(0).(*models.User); ok {
		return rec, args.Error(1)
	}
	return nil, args.Error(1)
}

func Test_DummyLogin_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepo)

	svc := New(mockUserRepo)

	token, err := svc.DummyLogin(models.UserRoleEmployee)

	require.NoError(t, err)
	require.NotEmpty(t, token)
}

func Test_Register_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepo)

	svc := New(mockUserRepo)

	ctx := context.Background()
	role := models.UserRoleEmployee
	password := "password"
	email := "employee@example.com"

	mockUserRepo.
		On("Create", ctx, email, mock.Anything, role).
		Return(true, nil)

	err := svc.Register(ctx, role, password, email)

	require.NoError(t, err)

	mockUserRepo.AssertExpectations(t)
}

func Test_Login_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepo)

	svc := New(mockUserRepo)

	ctx := context.Background()
	role := models.UserRoleEmployee
	password := "password"
	email := "employee@example.com"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	dummyUser := &models.User{
		ID:           uuid.New(),
		PasswordHash: string(hashedPassword),
		Role:         role,
		Email:        email,
	}

	mockUserRepo.
		On("FindUserByEmail", ctx, email).
		Return(dummyUser, nil)

	token, err := svc.Login(ctx, password, email)

	require.NoError(t, err)
	require.NotEmpty(t, token)

	mockUserRepo.AssertExpectations(t)
}

func Test_DummyLogin_InvalidRole(t *testing.T) {
	mockUserRepo := new(MockUserRepo)

	svc := New(mockUserRepo)

	token, err := svc.DummyLogin("test")

	require.Error(t, err)
	require.Empty(t, token)
	require.ErrorIs(t, err, core_errors.ErrInvalidRole)
}

func Test_Register_InvalidRole(t *testing.T) {
	mockUserRepo := new(MockUserRepo)

	svc := New(mockUserRepo)

	ctx := context.Background()
	role := models.UserRole("test")
	password := "password"
	email := "employee@example.com"

	mockUserRepo.
		On("Create", ctx, email, password, role).
		Return(true, nil)

	err := svc.Register(ctx, role, password, email)

	require.Error(t, err)
	require.ErrorIs(t, err, core_errors.ErrInvalidRole)
}

func Test_Register_InvalidRequest(t *testing.T) {
	mockUserRepo := new(MockUserRepo)

	svc := New(mockUserRepo)

	ctx := context.Background()
	role := models.UserRoleEmployee
	password := "password"
	email := "badEmail"

	mockUserRepo.
		On("Create", ctx, email, password, role).
		Return(true, nil)

	err := svc.Register(ctx, role, password, email)

	require.Error(t, err)
}

func Test_Login_InvalidCredentials(t *testing.T) {
	mockUserRepo := new(MockUserRepo)

	svc := New(mockUserRepo)

	ctx := context.Background()
	password := "password"
	email := "notexist@example.com"

	mockUserRepo.
		On("FindUserByEmail", ctx, email).
		Return(nil, core_errors.ErrUserNotFound)

	token, err := svc.Login(ctx, password, email)

	require.Error(t, err)
	require.Empty(t, token)
	require.ErrorIs(t, err, core_errors.ErrInvalidCredentials)
}

func Test_Login_InvalidRequest(t *testing.T) {
	mockUserRepo := new(MockUserRepo)

	svc := New(mockUserRepo)

	ctx := context.Background()
	role := models.UserRoleEmployee
	password := "password"
	email := "badEmail"

	dummyUser := &models.User{
		ID:           uuid.New(),
		PasswordHash: password,
		Role:         role,
		Email:        email,
	}
	mockUserRepo.
		On("FindUserByEmail", ctx, email).
		Return(dummyUser, nil)

	token, err := svc.Login(ctx, password, email)

	require.Error(t, err)
	require.Empty(t, token)
}
