package pvz

import (
	"avito/internal/models"
	"avito/internal/models/core_errors"
	"avito/internal/repositories/pvz"
	"avito/internal/utils"
	"avito/internal/utils/roles"
	"context"
)

type Service interface {
	Create(ctx context.Context, pvz *models.Pvz) (*models.Pvz, error)
	GetWithFilter(ctx context.Context, filter models.Filter) (*[]models.PzvDto, error)
}

type service struct {
	pvzRepo pvz.Repository
}

func New(pvzRepo pvz.Repository) Service {
	return &service{
		pvzRepo: pvzRepo,
	}
}

func (s *service) Create(ctx context.Context, pvz *models.Pvz) (*models.Pvz, error) {
	role, ok := roles.RoleFromContext(ctx)
	if !ok || role != models.UserRoleModerator {
		return nil, core_errors.ErrAccessDenied
	}

	if !utils.IsCityAllowed(pvz.City) {
		return nil, core_errors.ErrInvalidCity
	}

	pvz, err := s.pvzRepo.Create(ctx, pvz)
	if err != nil {
		return nil, err
	}

	return pvz, nil
}

func (s *service) GetWithFilter(ctx context.Context, filter models.Filter) (*[]models.PzvDto, error) {
	role, ok := roles.RoleFromContext(ctx)
	if !ok || !roles.IsRoleAllowed(role) {
		return nil, core_errors.ErrAccessDenied
	}

	res, err := s.pvzRepo.GetWithFilter(ctx, filter.StartDate, filter.EndDate, filter.Page, filter.Limit)
	if err != nil {
		return nil, err
	}
	return res, nil
}
