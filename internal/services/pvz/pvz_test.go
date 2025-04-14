package pvz

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

type MockPvzRepo struct {
	mock.Mock
}

func (m *MockPvzRepo) Create(ctx context.Context, pvz *models.Pvz) (*models.Pvz, error) {
	args := m.Called(ctx, pvz)
	if rec, ok := args.Get(0).(*models.Pvz); ok {
		return rec, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockPvzRepo) GetWithFilter(
	ctx context.Context,
	startDate, endDate *time.Time,
	page, limit int,
) (*[]models.PzvDto, error) {
	args := m.Called(ctx, startDate, endDate, page, limit)
	if rec, ok := args.Get(0).(*[]models.PzvDto); ok {
		return rec, args.Error(1)
	}
	return nil, args.Error(1)
}

func TestService_Create_Success(t *testing.T) {
	mockPvzRepo := new(MockPvzRepo)

	svc := New(mockPvzRepo)

	ctx := roles.WithRole(context.Background(), models.UserRoleModerator)

	pvzId := uuid.New()

	dummyPvz := &models.Pvz{
		ID:        pvzId,
		City:      models.Moscow,
		CreatedAt: time.Now(),
	}

	mockPvzRepo.
		On("Create", ctx, dummyPvz).
		Return(dummyPvz, nil)

	pvzResult, err := svc.Create(ctx, dummyPvz)

	assert.NoError(t, err)
	assert.Equal(t, dummyPvz, pvzResult)

	mockPvzRepo.AssertExpectations(t)
}

func TestService_GetWithFilter_Success(t *testing.T) {
	mockPvzRepo := new(MockPvzRepo)

	svc := New(mockPvzRepo)

	ctx := roles.WithRole(context.Background(), models.UserRoleModerator)

	page := 1
	limit := 10
	filter := models.Filter{
		Page:  page,
		Limit: limit,
	}

	dummyPvzs := &[]models.PzvDto{}

	mockPvzRepo.
		On("GetWithFilter", ctx, filter.StartDate, filter.EndDate, page, limit).
		Return(dummyPvzs, nil)

	pvzResult, err := svc.GetWithFilter(ctx, filter)

	assert.NoError(t, err)
	assert.Equal(t, dummyPvzs, pvzResult)

	mockPvzRepo.AssertExpectations(t)
}

func TestService_Create_AccessDenied(t *testing.T) {
	mockPvzRepo := new(MockPvzRepo)

	svc := New(mockPvzRepo)

	ctx := roles.WithRole(context.Background(), models.UserRoleEmployee)

	pvzId := uuid.New()

	dummyPvz := &models.Pvz{
		ID:        pvzId,
		City:      models.Moscow,
		CreatedAt: time.Now(),
	}

	mockPvzRepo.
		On("Create", ctx, dummyPvz).
		Return(dummyPvz, nil)

	pvzResult, err := svc.Create(ctx, dummyPvz)

	assert.Error(t, err)
	assert.Nil(t, pvzResult)
	assert.True(t, errors.Is(err, core_errors.ErrAccessDenied))
}

func TestService_Create_InvalidCity(t *testing.T) {
	mockPvzRepo := new(MockPvzRepo)

	svc := New(mockPvzRepo)

	ctx := roles.WithRole(context.Background(), models.UserRoleModerator)

	pvzId := uuid.New()

	dummyPvz := &models.Pvz{
		ID:        pvzId,
		City:      "Tyumen",
		CreatedAt: time.Now(),
	}

	mockPvzRepo.
		On("Create", ctx, dummyPvz).
		Return(dummyPvz, nil)

	pvzResult, err := svc.Create(ctx, dummyPvz)

	assert.Error(t, err)
	assert.Nil(t, pvzResult)
	assert.True(t, errors.Is(err, core_errors.ErrInvalidCity))
}

func TestService_GetWithFilter_AccessDenied(t *testing.T) {
	mockPvzRepo := new(MockPvzRepo)

	svc := New(mockPvzRepo)

	ctx := context.Background()

	page := 1
	limit := 10
	filter := models.Filter{
		Page:  page,
		Limit: limit,
	}

	dummyPvzs := &[]models.PzvDto{}

	mockPvzRepo.
		On("GetWithFilter", ctx, filter.StartDate, filter.EndDate, page, limit).
		Return(dummyPvzs, nil)

	pvzResult, err := svc.GetWithFilter(ctx, filter)

	assert.Error(t, err)
	assert.Nil(t, pvzResult)
	assert.True(t, errors.Is(err, core_errors.ErrAccessDenied))
}
