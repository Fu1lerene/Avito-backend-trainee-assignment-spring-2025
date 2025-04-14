package product

import (
	"avito/internal/models"
	"avito/internal/models/core_errors"
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, receptionID uuid.UUID, productType models.ProductType) (*models.Product, error)
	DeleteLast(ctx context.Context, receptionID uuid.UUID) error
}

type repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, receptionID uuid.UUID, productType models.ProductType) (*models.Product, error) {
	row := r.db.QueryRow(ctx, `
		INSERT INTO products (reception_id, type)
		VALUES ($1, $2)
		RETURNING id, reception_id, type, created_at
	`, receptionID, productType)

	var product models.Product
	if err := row.Scan(&product.ID, &product.ReceptionID, &product.Type, &product.CreatedAt); err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *repository) DeleteLast(ctx context.Context, receptionID uuid.UUID) error {
	cmd, err := r.db.Exec(ctx, `
		DELETE FROM products
		WHERE id = (
			SELECT id FROM products
			WHERE reception_id = $1
			ORDER BY created_at DESC
			LIMIT 1
		)
	`, receptionID)

	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return core_errors.ErrNoProductsToDelete
	}
	return nil
}
