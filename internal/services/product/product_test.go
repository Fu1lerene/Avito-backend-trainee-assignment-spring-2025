package product

import (
	"avito/internal/models"
	"avito/internal/models/core_errors"
	"avito/internal/utils/roles"
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
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

func (m *MockReceptionRepo) Close(ctx context.Context, pvzID uuid.UUID) error {
	args := m.Called(ctx, pvzID)
	if _, ok := args.Get(0).(bool); ok {
		return args.Error(1)
	}
	return args.Error(1)
}

type MockProductRepo struct {
	mock.Mock
}

func (m *MockProductRepo) Create(ctx context.Context, receptionID uuid.UUID, productType models.ProductType) (*models.Product, error) {
	args := m.Called(ctx, receptionID, productType)
	if prod, ok := args.Get(0).(*models.Product); ok {
		return prod, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockProductRepo) DeleteLast(ctx context.Context, receptionID uuid.UUID) error {
	args := m.Called(ctx, receptionID)
	if _, ok := args.Get(0).(bool); ok {
		return args.Error(1)
	}
	return args.Error(1)
}

func TestService_Add_Success(t *testing.T) {
	mockReceptionRepo := new(MockReceptionRepo)
	mockProductRepo := new(MockProductRepo)

	svc := New(mockReceptionRepo, mockProductRepo)

	ctx := roles.WithRole(context.Background(), models.UserRoleEmployee)

	pvzId := uuid.New()
	receptionID := uuid.New()
	productType := models.Shoes
	dummyReception := &models.Reception{
		ID:        receptionID,
		PvzID:     pvzId,
		Status:    models.InProgress,
		CreatedAt: time.Now(),
	}
	dummyProduct := &models.Product{
		ID:          uuid.New(),
		ReceptionID: receptionID,
		Type:        productType,
		CreatedAt:   time.Now(),
	}

	mockReceptionRepo.
		On("GetLastActiveReception", ctx, pvzId).
		Return(dummyReception, nil)

	mockProductRepo.
		On("Create", ctx, receptionID, productType).
		Return(dummyProduct, nil)

	productResult, err := svc.Add(ctx, productType, pvzId)

	require.NoError(t, err)
	require.Equal(t, dummyProduct, productResult)

	mockReceptionRepo.AssertExpectations(t)
	mockProductRepo.AssertExpectations(t)
}

func TestService_DeleteLastProduct_Success(t *testing.T) {
	mockReceptionRepo := new(MockReceptionRepo)
	mockProductRepo := new(MockProductRepo)

	svc := New(mockReceptionRepo, mockProductRepo)

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
		On("GetLastActiveReception", ctx, pvzId).
		Return(dummyReception, nil)

	mockProductRepo.
		On("DeleteLast", ctx, receptionID).
		Return(true, nil)

	err := svc.DeleteLastProduct(ctx, pvzId)

	require.NoError(t, err)

	mockReceptionRepo.AssertExpectations(t)
	mockProductRepo.AssertExpectations(t)
}

func TestService_Delete_AccessDenied(t *testing.T) {
	mockReceptionRepo := new(MockReceptionRepo)
	mockProductRepo := new(MockProductRepo)

	svc := New(mockReceptionRepo, mockProductRepo)

	ctx := context.Background()

	err := svc.DeleteLastProduct(ctx, uuid.New())

	require.Error(t, err)
	require.ErrorIs(t, err, core_errors.ErrAccessDenied)
}

func TestService_Add_AccessDenied(t *testing.T) {
	mockReceptionRepo := new(MockReceptionRepo)
	mockProductRepo := new(MockProductRepo)

	svc := New(mockReceptionRepo, mockProductRepo)

	ctx := context.Background()

	productResult, err := svc.Add(ctx, models.Shoes, uuid.New())

	require.Nil(t, productResult)
	require.Error(t, err)
	require.ErrorIs(t, err, core_errors.ErrAccessDenied)
}

func TestService_Add_InvalidProductType(t *testing.T) {
	mockReceptionRepo := new(MockReceptionRepo)
	mockProductRepo := new(MockProductRepo)

	svc := New(mockReceptionRepo, mockProductRepo)

	ctx := roles.WithRole(context.Background(), models.UserRoleEmployee)

	productResult, err := svc.Add(ctx, "Книги", uuid.New())

	require.Nil(t, productResult)
	require.Error(t, err)
	require.ErrorIs(t, err, core_errors.ErrInvalidProductType)
}
