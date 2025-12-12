package main

import (
	database "api/internal/config/db"
	"api/internal/config/env"

	"log"

	_ "github.com/lib/pq"
)

func main() {
	// Load configuration
	cfg, err := env.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Seed database
	err = database.SeedDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}

	log.Println("Database seeded successfully!")
}
