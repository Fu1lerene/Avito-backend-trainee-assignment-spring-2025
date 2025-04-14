package tests

import (
	"avito/internal/models"
	productrepo "avito/internal/repositories/product"
	pvzrepo "avito/internal/repositories/pvz"
	receptionrepo "avito/internal/repositories/reception"
	"avito/internal/services/product"
	"avito/internal/services/pvz"
	"avito/internal/services/reception"
	"avito/internal/utils/roles"
	"avito/tests/helper"
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"testing"
	"time"
)

type Suite struct {
	suite.Suite
	cont             *postgres.PostgresContainer
	pool             *pgxpool.Pool
	productService   product.Service
	receptionService reception.Service
	pvzService       pvz.Service
}

func (s *Suite) SetupSuite() {
	cont, pool := helper.SetupPG(s.T())
	s.cont, s.pool = cont, pool
}

var ctx = context.Background()

func (s *Suite) SetupTest() {
	_, err := s.pool.Exec(ctx, `
		TRUNCATE products, pvz_points, receptions, users RESTART IDENTITY CASCADE
	`)
	s.Require().NoError(err)

	pvzRepo := pvzrepo.New(s.pool)
	recRepo := receptionrepo.New(s.pool)
	prodRepo := productrepo.New(s.pool)

	s.productService = product.New(recRepo, prodRepo)
	s.receptionService = reception.New(recRepo)
	s.pvzService = pvz.New(pvzRepo)
}

func (s *Suite) TearDownSuite() {
	s.pool.Close()
	s.Require().NoError(s.cont.Terminate(ctx))
}

func (s *Suite) TestIntegration() {
	ctxUser := roles.WithRole(context.Background(), models.UserRoleEmployee)
	ctxMod := roles.WithRole(context.Background(), models.UserRoleModerator)
	pvzID := uuid.New()
	pvzRequest := &models.Pvz{
		ID:        pvzID,
		City:      models.Moscow,
		CreatedAt: time.Now(),
	}

	pvzResponse, err := s.pvzService.Create(ctxMod, pvzRequest)
	s.Require().NoError(err)
	s.Require().Equal(pvzRequest, pvzResponse)

	_, err = s.receptionService.Open(ctxUser, pvzID)
	s.Require().NoError(err)

	for i := 0; i < 50; i++ {
		_, err := s.productService.Add(ctxUser, models.Electronics, pvzID)
		s.Require().NoError(err)
	}

	err = s.receptionService.Close(ctxUser, pvzID)
	s.Require().NoError(err)
}

func TestEndToEndSuccess(t *testing.T) {
	suite.Run(t, new(Suite))
}
