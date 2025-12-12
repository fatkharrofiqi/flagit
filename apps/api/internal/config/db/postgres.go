package database

import (
	"api/db/seed"
	"api/internal/config/env"
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os/exec"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

func InitializeDB(cfg *env.Config) (*sql.DB, error) {
	dsn := cfg.Database.GetDSN()

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Run migrations
	// if err := runMigrations(cfg); err != nil {
	// 	return nil, fmt.Errorf("failed to run migrations: %w", err)
	// }

	// Run seeds
	if err := runSeeds(db); err != nil {
		log.Printf("Warning: failed to run seeds: %v", err)
	}

	return db, nil
}

func runMigrations(cfg *env.Config) error {
	// Extract connection details from DSN
	u, err := url.Parse(cfg.Database.GetDSN())
	if err != nil {
		return fmt.Errorf("failed to parse database URL: %w", err)
	}

	dbName := strings.TrimPrefix(u.Path, "/")
	host := u.Hostname()
	port := u.Port()
	user := u.User.Username()
	password, _ := u.User.Password()

	// Build migration command
	migratePath := "file://db/migration"
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		user, password, host, port, dbName, cfg.Database.SSLMode)

	cmd := exec.Command("migrate", "-path", migratePath, "-database", dbURL, "up")
	cmd.Dir = "../../apps/api" // Set working directory

	output, err := cmd.CombinedOutput()
	if err != nil {
		if strings.Contains(string(output), "no change") {
			log.Println("Database is already up to date")
			return nil
		}
		return fmt.Errorf("migration failed: %w\nOutput: %s", err, string(output))
	}

	log.Printf("Migration completed successfully")
	return nil
}

func runSeeds(db *sql.DB) error {
	seeder := seed.NewSeeder(db)
	return seeder.Seed(context.Background())
}
