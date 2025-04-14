package reception

import (
	"avito/internal/models"
	"avito/internal/repositories/reception"
	core_errors "avito/internal/utils/errors"
	"avito/internal/utils/roles"
	"context"
	"github.com/google/uuid"
)

type Service interface {
	Open(ctx context.Context, pvzId uuid.UUID) (*models.Reception, error)
	Close(ctx context.Context, pvzId uuid.UUID) (bool, error)
}

type service struct {
	receptionRepo reception.Repository
}

func New(receptionRepo reception.Repository) Service {
	return &service{
		receptionRepo: receptionRepo,
	}
}

func (s *service) Open(ctx context.Context, pvzId uuid.UUID) (*models.Reception, error) {
	role, ok := roles.RoleFromContext(ctx)
	if !ok || role != models.UserRoleEmployee {
		return nil, core_errors.ErrAccessDenied
	}

	hasActive, err := s.receptionRepo.HasActiveReception(ctx, pvzId)
	if err != nil {
		return nil, err
	}
	if hasActive {
		return nil, core_errors.ErrReceptionNotClosed
	}

	rec, err := s.receptionRepo.Create(ctx, pvzId)
	if err != nil {
		return nil, err
	}

	return rec, nil
}
func (s *service) Close(ctx context.Context, pvzId uuid.UUID) (bool, error) {
	role, ok := roles.RoleFromContext(ctx)
	if !ok || role != models.UserRoleEmployee {
		return false, core_errors.ErrAccessDenied
	}

	_, err := s.receptionRepo.Close(ctx, pvzId)
	if err != nil {
		return false, err
	}

	return true, nil
}
