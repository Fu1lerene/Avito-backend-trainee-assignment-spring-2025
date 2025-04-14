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
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
	"time"
)

func Test(t *testing.T) {
	db := InitTestDB()

	pvzRepo := pvzrepo.New(db)
	recRepo := receptionrepo.New(db)
	prodRepo := productrepo.New(db)

	productService := product.New(recRepo, prodRepo)
	receptionService := reception.New(recRepo)
	pvzService := pvz.New(pvzRepo)
	ctxUser := roles.WithRole(context.Background(), models.UserRoleEmployee)
	ctxMod := roles.WithRole(context.Background(), models.UserRoleModerator)
	pvzID := uuid.New()
	pvzRequest := &models.Pvz{
		ID:        pvzID,
		City:      models.Moscow,
		CreatedAt: time.Now(),
	}

	pvzResponse, err := pvzService.Create(ctxMod, pvzRequest)
	require.NoError(t, err)
	assert.Equal(t, pvzRequest, pvzResponse)

	_, err = receptionService.Open(ctxUser, pvzID)
	require.NoError(t, err)

	for i := 0; i < 50; i++ {
		_, err := productService.Add(ctxUser, models.Electronics, pvzID)
		require.NoError(t, err)
	}

	ok, err := receptionService.Close(ctxUser, pvzID)
	require.NoError(t, err)
	assert.True(t, ok)
}

func InitTestDB() *pgxpool.Pool {
	pool, err := pgxpool.New(context.Background(), "postgres://postgres:postgres@localhost:5544/avito?sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to connect PostgreSQL: %v", err)
	}
	log.Println("PostgreSQL connected")

	return pool
}
