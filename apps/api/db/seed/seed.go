package seed

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type Seeder struct {
	db *sql.DB
}

func NewSeeder(db *sql.DB) *Seeder {
	return &Seeder{db: db}
}

func (s *Seeder) Seed(ctx context.Context) error {
	// Create demo project
	projectID := uuid.New()
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO projects (id, name, description) 
		VALUES ($1, $2, $3) 
		ON CONFLICT DO NOTHING`,
		projectID, "Demo Project", "A demonstration project for Flagit")
	if err != nil {
		return err
	}

	// Create demo environments
	envDevID := uuid.New()
	envProdID := uuid.New()
	
	_, err = s.db.ExecContext(ctx, `
		INSERT INTO environments (id, project_id, name) VALUES 
		($1, $2, $3),
		($4, $2, $5) 
		ON CONFLICT DO NOTHING`,
		envDevID, projectID, "Development", envProdID, projectID, "Production")
	if err != nil {
		return err
	}

	// Create demo flags
	flagBetaFeaturesID := uuid.New()
	flagMaintenanceID := uuid.New()
	
	_, err = s.db.ExecContext(ctx, `
		INSERT INTO flags (id, project_id, key, description, type) VALUES 
		($1, $2, $3, $4, $5),
		($6, $2, $7, $8, $9) 
		ON CONFLICT DO NOTHING`,
		flagBetaFeaturesID, projectID, "beta_features", "Enable beta features", "boolean",
		flagMaintenanceID, projectID, "maintenance_mode", "Enable maintenance mode", "boolean")
	if err != nil {
		return err
	}

	// Create demo flag values
	flagValueBetaDevID := uuid.New()
	flagValueBetaProdID := uuid.New()
	flagValueMaintDevID := uuid.New()
	flagValueMaintProdID := uuid.New()
	
	_, err = s.db.ExecContext(ctx, `
		INSERT INTO flag_values (id, flag_id, env_id, value, enabled) VALUES 
		($1, $2, $3, $4, $5),
		($6, $7, $8, $9, $10),
		($11, $12, $13, $14, $15),
		($16, $17, $18, $19, $20) 
		ON CONFLICT DO NOTHING`,
		flagValueBetaDevID, flagBetaFeaturesID, envDevID, "true", true,
		flagValueBetaProdID, flagBetaFeaturesID, envProdID, "false", false,
		flagValueMaintDevID, flagMaintenanceID, envDevID, "false", false,
		flagValueMaintProdID, flagMaintenanceID, envProdID, "false", false)
	if err != nil {
		return err
	}

	return nil
}
