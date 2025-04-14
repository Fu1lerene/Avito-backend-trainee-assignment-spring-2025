package product

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

func (m *MockProductRepo) DeleteLast(ctx context.Context, receptionID uuid.UUID) (bool, error) {
	args := m.Called(ctx, receptionID)
	if rec, ok := args.Get(0).(bool); ok {
		return rec, args.Error(1)
	}
	return false, args.Error(1)
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

	assert.NoError(t, err)
	assert.Equal(t, dummyProduct, productResult)

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

	result, err := svc.DeleteLastProduct(ctx, pvzId)

	assert.NoError(t, err)
	assert.True(t, result)

	mockReceptionRepo.AssertExpectations(t)
	mockProductRepo.AssertExpectations(t)
}

func TestService_Delete_AccessDenied(t *testing.T) {
	mockReceptionRepo := new(MockReceptionRepo)
	mockProductRepo := new(MockProductRepo)

	svc := New(mockReceptionRepo, mockProductRepo)

	ctx := context.Background()

	ok, err := svc.DeleteLastProduct(ctx, uuid.New())

	assert.False(t, ok)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, core_errors.ErrAccessDenied))
}

func TestService_Add_AccessDenied(t *testing.T) {
	mockReceptionRepo := new(MockReceptionRepo)
	mockProductRepo := new(MockProductRepo)

	svc := New(mockReceptionRepo, mockProductRepo)

	ctx := context.Background()

	productResult, err := svc.Add(ctx, models.Shoes, uuid.New())

	assert.Nil(t, productResult)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, core_errors.ErrAccessDenied))
}

func TestService_Add_InvalidProductType(t *testing.T) {
	mockReceptionRepo := new(MockReceptionRepo)
	mockProductRepo := new(MockProductRepo)

	svc := New(mockReceptionRepo, mockProductRepo)

	ctx := roles.WithRole(context.Background(), models.UserRoleEmployee)

	productResult, err := svc.Add(ctx, "Книги", uuid.New())

	assert.Nil(t, productResult)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, core_errors.ErrInvalidProductType))
}
