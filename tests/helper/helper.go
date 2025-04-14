package helper

import (
	"avito/internal/migrations"
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/require"
	pg "github.com/testcontainers/testcontainers-go/modules/postgres"
)

func SetupPG(
	t *testing.T,
) (*pg.PostgresContainer, *pgxpool.Pool) {
	ctx := context.Background()
	cont, err := pg.Run(ctx,
		"postgres:16-alpine",
		pg.WithDatabase("postgres_test"),
		pg.WithUsername("postgres"),
		pg.WithPassword("postgres"),
		pg.BasicWaitStrategies(),
	)
	require.NoError(t, err)

	dsn, err := cont.ConnectionString(ctx)
	require.NoError(t, err)
	pool, err := pgxpool.New(ctx, dsn)
	require.NoError(t, err)

	goose.SetBaseFS(migrations.All)
	db := stdlib.OpenDBFromPool(pool)

	require.NoError(t, goose.Up(db, "."))
	require.NoError(t, db.Close())

	return cont, pool
}
