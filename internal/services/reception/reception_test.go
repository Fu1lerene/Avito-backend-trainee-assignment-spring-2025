package reception

import (
	"avito/internal/models"
	core_errors "avito/internal/utils/errors"
	"avito/internal/utils/roles"
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

type MockReceptionRepo struct {
	mock.Mock
}

func (m *MockReceptionRepo) GetLastActiveReception(ctx context.Context, pvzId uuid.UUID) (*models.Reception, error) {
	args := m.Called(ctx, pvzId)
	if rec, ok := args.Get(0).(*models.Reception); ok {
		return rec, args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockReceptionRepo) HasActiveReception(ctx context.Context, pvzID uuid.UUID) (bool, error) {
	args := m.Called(ctx, pvzID)
	if rec, ok := args.Get(0).(bool); ok {
		return rec, args.Error(1)
	}
	return false, args.Error(1)
}
func (m *MockReceptionRepo) Create(ctx context.Context, pvzID uuid.UUID) (*models.Reception, error) {
	args := m.Called(ctx, pvzID)
	if rec, ok := args.Get(0).(*models.Reception); ok {
		return rec, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockReceptionRepo) Close(ctx context.Context, pvzID uuid.UUID) (bool, error) {
	args := m.Called(ctx, pvzID)
	if rec, ok := args.Get(0).(bool); ok {
		return rec, args.Error(1)
	}
	return false, args.Error(1)
}

func TestService_Open_Success(t *testing.T) {
	mockReceptionRepo := new(MockReceptionRepo)

	svc := New(mockReceptionRepo)

	ctx := roles.WithRole(context.Background(), models.UserRoleEmployee)

	pvzId := uuid.New()
	receptionID := uuid.New()
	dummyReception := &models.Reception{
		ID:        receptionID,
		PvzID:     pvzId,
		Status:    models.InProgress,
		CreatedAt: time.Now(),
	}

	mockReceptionRepo.
		On("HasActiveReception", ctx, pvzId).
		Return(false, nil)

	mockReceptionRepo.
		On("Create", ctx, pvzId).
		Return(dummyReception, nil)

	rec, err := svc.Open(ctx, pvzId)

	assert.NoError(t, err)
	assert.Equal(t, dummyReception, rec)

	mockReceptionRepo.AssertExpectations(t)
}

func TestService_Close_Success(t *testing.T) {
	mockReceptionRepo := new(MockReceptionRepo)

	svc := New(mockReceptionRepo)

	ctx := roles.WithRole(context.Background(), models.UserRoleEmployee)

	pvzId := uuid.New()

	mockReceptionRepo.
		On("Close", ctx, pvzId).
		Return(false, nil)

	ok, err := svc.Close(ctx, pvzId)

	assert.NoError(t, err)
	assert.True(t, ok)

	mockReceptionRepo.AssertExpectations(t)
}

func TestService_Open_AccessDenied(t *testing.T) {
	mockReceptionRepo := new(MockReceptionRepo)

	svc := New(mockReceptionRepo)

	ctx := roles.WithRole(context.Background(), models.UserRoleModerator)

	pvzId := uuid.New()
	receptionID := uuid.New()
	dummyReception := &models.Reception{
		ID:        receptionID,
		PvzID:     pvzId,
		Status:    models.InProgress,
		CreatedAt: time.Now(),
	}

	mockReceptionRepo.
		On("HasActiveReception", ctx, pvzId).
		Return(false, nil)

	mockReceptionRepo.
		On("Create", ctx, pvzId).
		Return(dummyReception, nil)

	rec, err := svc.Open(ctx, pvzId)

	assert.Error(t, err)
	assert.Nil(t, rec)
	assert.True(t, errors.Is(err, core_errors.ErrAccessDenied))
}

func TestService_Close_AccessDenied(t *testing.T) {
	mockReceptionRepo := new(MockReceptionRepo)

	svc := New(mockReceptionRepo)

	ctx := roles.WithRole(context.Background(), models.UserRoleModerator)

	pvzId := uuid.New()
	mockReceptionRepo.
		On("Close", ctx, pvzId).
		Return(false, nil)

	ok, err := svc.Close(ctx, pvzId)

	assert.Error(t, err)
	assert.False(t, ok)
	assert.True(t, errors.Is(err, core_errors.ErrAccessDenied))
}

func TestService_Open_ReceptionNotClosed(t *testing.T) {
	mockReceptionRepo := new(MockReceptionRepo)

	svc := New(mockReceptionRepo)

	ctx := roles.WithRole(context.Background(), models.UserRoleEmployee)

	pvzId := uuid.New()
	mockReceptionRepo.
		On("HasActiveReception", ctx, pvzId).
		Return(true, nil)

	rec, err := svc.Open(ctx, pvzId)

	assert.Error(t, err)
	assert.Nil(t, rec)
	assert.True(t, errors.Is(err, core_errors.ErrReceptionNotClosed))
}
