package user

import (
	"avito/internal/models"
	core_errors "avito/internal/utils/errors"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, email, hashedPassword string, role models.UserRole) (bool, error)
	FindUserByEmail(ctx context.Context, email string) (*models.User, error)
}

type repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, email, hashedPassword string, role models.UserRole) (bool, error) {
	_, err := r.db.Exec(ctx, `
		INSERT INTO users (email, password, role) 
		VALUES ($1, $2, $3)
	`, email, hashedPassword, role)

	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *repository) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, email, password, role 
		FROM users 
		WHERE email = $1
	`, email)

	var user models.User
	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Role)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, core_errors.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}
