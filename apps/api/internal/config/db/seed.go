package database

import (
	"api/db/seed"
	"api/internal/config/env"
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func InitializeDBForSeed(cfg *env.Config) (*sql.DB, error) {
	dsn := cfg.Database.GetDSN()

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test connection
	ctx := context.Background()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

func SeedDatabase(cfg *env.Config) error {
	db, err := InitializeDBForSeed(cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	defer db.Close()

	// Run seeds
	seeder := seed.NewSeeder(db)
	return seeder.Seed(context.Background())
}
