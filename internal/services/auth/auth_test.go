package auth

import (
	"avito/internal/models"
	core_errors "avito/internal/utils/errors"
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) Create(ctx context.Context, email, hashedPassword string, role models.UserRole) (bool, error) {
	args := m.Called(ctx, email, hashedPassword, role)
	if rec, ok := args.Get(0).(bool); ok {
		return rec, args.Error(1)
	}
	return false, args.Error(1)
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

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

//
//func Test_Register_Success(t *testing.T) {
//	mockUserRepo := new(MockUserRepo)
//
//	svc := New(mockUserRepo)
//
//	ctx := context.Background()
//	role := models.UserRoleEmployee
//	password := "password"
//	email := "employee@example.com"
//	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
//
//	mockUserRepo.
//		On("Create", ctx, email, string(hashedPassword), role).
//		Return(true, nil)
//
//	ok, err := svc.Register(ctx, role, password, email)
//
//	assert.NoError(t, err)
//	assert.True(t, ok)
//
//	mockUserRepo.AssertExpectations(t)
//}

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

	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	mockUserRepo.AssertExpectations(t)
}

func Test_DummyLogin_InvalidRole(t *testing.T) {
	mockUserRepo := new(MockUserRepo)

	svc := New(mockUserRepo)

	token, err := svc.DummyLogin("test")

	assert.Error(t, err)
	assert.Empty(t, token)
	assert.True(t, errors.Is(err, core_errors.ErrInvalidRole))
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

	ok, err := svc.Register(ctx, role, password, email)

	assert.Error(t, err)
	assert.False(t, ok)
	assert.True(t, errors.Is(err, core_errors.ErrInvalidRole))
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

	ok, err := svc.Register(ctx, role, password, email)

	assert.Error(t, err)
	assert.False(t, ok)
	assert.True(t, errors.Is(err, core_errors.ErrInvalidRequest))
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

	assert.Error(t, err)
	assert.Empty(t, token)
	assert.True(t, errors.Is(err, core_errors.ErrInvalidCredentials))

	mockUserRepo.AssertExpectations(t)
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

	assert.Error(t, err)
	assert.Empty(t, token)
	assert.True(t, errors.Is(err, core_errors.ErrInvalidRequest))

	mockUserRepo.AssertExpectations(t)
}
