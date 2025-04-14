package reception

import (
	"avito/internal/models"
	"avito/internal/models/core_errors"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, pvzID uuid.UUID) (*models.Reception, error)
	HasActiveReception(ctx context.Context, pvzID uuid.UUID) (bool, error)
	GetLastActiveReception(ctx context.Context, pvzId uuid.UUID) (*models.Reception, error)
	Close(ctx context.Context, pvzID uuid.UUID) error
}
type repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, pvzID uuid.UUID) (*models.Reception, error) {
	row := r.db.QueryRow(ctx, `
		INSERT INTO receptions (pvz_id, status) 
		VALUES ($1, $2) 
		RETURNING id, pvz_id, created_at, status
	`, pvzID, models.InProgress)

	var rec models.Reception
	if err := row.Scan(&rec.ID, &rec.PvzID, &rec.CreatedAt, &rec.Status); err != nil {
		return nil, err
	}
	return &rec, nil
}

func (r *repository) HasActiveReception(ctx context.Context, pvzID uuid.UUID) (bool, error) {
	row := r.db.QueryRow(ctx, `
		SELECT 1 FROM receptions 
		WHERE pvz_id = $1 AND status = $2
	`, pvzID, models.InProgress)

	var exist int
	if err := row.Scan(&exist); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *repository) GetLastActiveReception(ctx context.Context, pvzID uuid.UUID) (*models.Reception, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, pvz_id, created_at, status FROM receptions 
		WHERE pvz_id = $1 AND status = $2
		ORDER BY created_at DESC
		LIMIT 1
	`, pvzID, models.InProgress)

	var rec models.Reception
	err := row.Scan(&rec.ID, &rec.PvzID, &rec.CreatedAt, &rec.Status)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, core_errors.ErrNoActiveReceptions
		}
		return nil, err
	}
	return &rec, nil
}

func (r *repository) Close(ctx context.Context, pvzID uuid.UUID) error {
	cmd, err := r.db.Exec(ctx, `
		UPDATE receptions
		SET status = $1
		WHERE pvz_id = $2 AND status = $3
	`, models.Close, pvzID, models.InProgress)

	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return core_errors.ErrDeleteReception
	}
	return nil
}
