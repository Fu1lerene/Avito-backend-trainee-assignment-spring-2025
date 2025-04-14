package database

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func MustInit() *pgxpool.Pool {
	dsn := os.Getenv("DATABASE_URL")

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Failed to connect PostgreSQL: %v", err)
	}
	log.Println("PostgreSQL connected")

	return pool
}
