package pvz

import (
	"avito/internal/models"
	"context"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, pvz *models.Pvz) (*models.Pvz, error)
	GetWithFilter(
		ctx context.Context,
		startDate, endDate *time.Time,
		page, limit int,
	) (*[]models.PzvDto, error)
}

type repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, pvz *models.Pvz) (*models.Pvz, error) {
	row := r.db.QueryRow(ctx, `
		INSERT INTO pvz_points (id, city, created_at) 
		VALUES ($1, $2, $3) 
		RETURNING id, city, created_at
	`, pvz.ID, pvz.City, pvz.CreatedAt)

	if err := row.Scan(&pvz.ID, &pvz.City, &pvz.CreatedAt); err != nil {
		return nil, err
	}
	return pvz, nil
}

func (r *repository) GetWithFilter(
	ctx context.Context,
	startDate, endDate *time.Time,
	page, limit int,
) (*[]models.PzvDto, error) {
	query := `
WITH pvz_json AS (
  SELECT json_build_object(
    'pvz', json_build_object(
              'city', p.city,
              'id', p.id,
              'registrationDate', to_char(p.created_at, 'YYYY-MM-DD"T"HH24:MI:SS.MS"Z"')
           ),
    'receptions', (
      SELECT COALESCE(
        json_agg(
          json_build_object(
            'reception', json_build_object(
                            'id', r.id,
                            'pvzId', r.pvz_id,
                            'status', r.status,
                            'dateTime', to_char(r.created_at, 'YYYY-MM-DD"T"HH24:MI:SS.MS"Z"')
                         ),
            'products', (
              SELECT COALESCE(
                json_agg(
                  json_build_object(
                    'id', pr.id,
                    'receptionId', pr.reception_id,
                    'type', pr.type,
                    'dateTime', to_char(pr.created_at, 'YYYY-MM-DD"T"HH24:MI:SS.MS"Z"')
                  )
                ),
                '[]'::json
              )
              FROM products pr
              WHERE pr.reception_id = r.id
            )
          )
        ),
        '[]'::json
      )
      FROM receptions r
      WHERE r.pvz_id = p.id
        AND ($1::timestamp IS NULL OR r.created_at >= $1)
        AND ($2::timestamp IS NULL OR r.created_at <= $2)
    )
  ) AS pvz_data
  FROM pvz_points p
  WHERE 
      (EXISTS (
         SELECT 1 FROM receptions r 
         WHERE r.pvz_id = p.id 
         	AND ($1::timestamp IS NULL OR r.created_at >= $1)
         	AND ($2::timestamp IS NULL OR r.created_at <= $2)
      )
		OR (($1::timestamp IS NULL OR p.created_at >= $1)
            AND ($2::timestamp IS NULL OR p.created_at <= $2)))
  ORDER BY p.created_at DESC
  LIMIT $3 OFFSET $4
)
SELECT json_agg(pvz_data) AS result FROM pvz_json;
    `
	offset := (page - 1) * limit
	var jsonResult []byte
	if err := r.db.QueryRow(ctx, query, startDate, endDate, limit, offset).Scan(&jsonResult); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &[]models.PzvDto{}, nil
		}
		return nil, err
	}

	if len(jsonResult) <= 0 {
		return &[]models.PzvDto{}, nil
	}

	var pvzList *[]models.PzvDto
	if err := json.Unmarshal(jsonResult, &pvzList); err != nil {
		return nil, err
	}

	return pvzList, nil
}
