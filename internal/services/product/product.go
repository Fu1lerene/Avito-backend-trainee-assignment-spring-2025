package product

import (
	"avito/internal/models"
	"avito/internal/models/core_errors"
	"avito/internal/repositories/product"
	"avito/internal/repositories/reception"
	"avito/internal/utils"
	"avito/internal/utils/roles"
	"context"
	"github.com/google/uuid"
)

type Service interface {
	Add(ctx context.Context, productType models.ProductType, PvzId uuid.UUID) (*models.Product, error)
	DeleteLastProduct(ctx context.Context, pvzId uuid.UUID) error
}

type service struct {
	receptionRepo reception.Repository
	productRepo   product.Repository
}

func New(receptionRepo reception.Repository, productRepo product.Repository) Service {
	return &service{
		receptionRepo: receptionRepo,
		productRepo:   productRepo,
	}
}

func (s *service) Add(ctx context.Context, productType models.ProductType, pvzId uuid.UUID) (*models.Product, error) {
	role, ok := roles.RoleFromContext(ctx)
	if !ok || role != models.UserRoleEmployee {
		return nil, core_errors.ErrAccessDenied
	}

	if !utils.IsProductTypeAllowed(productType) {
		return nil, core_errors.ErrInvalidProductType
	}

	reception, err := s.receptionRepo.GetLastActiveReception(ctx, pvzId)
	if err != nil {
		return nil, err
	}

	product, err := s.productRepo.Create(ctx, reception.ID, productType)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (s *service) DeleteLastProduct(ctx context.Context, pvzId uuid.UUID) error {
	role, ok := roles.RoleFromContext(ctx)
	if !ok || role != models.UserRoleEmployee {
		return core_errors.ErrAccessDenied
	}

	reception, err := s.receptionRepo.GetLastActiveReception(ctx, pvzId)
	if err != nil {
		return err
	}

	err = s.productRepo.DeleteLast(ctx, reception.ID)
	if err != nil {
		return err
	}

	return nil
}
