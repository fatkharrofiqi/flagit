package repository

import (
	"context"
	"database/sql"
	"time"

	"api/internal/dto"
	"api/internal/model"

	"github.com/google/uuid"
)

type flagValueRepository struct {
	db *sql.DB
}

func NewFlagValueRepository(db *sql.DB) FlagValueRepository {
	return &flagValueRepository{db: db}
}

func (r *flagValueRepository) Create(ctx context.Context, flagValue *model.FlagValue) error {
	query := `
		INSERT INTO flag_values (id, flag_id, env_id, value, enabled, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (flag_id, env_id) DO UPDATE SET
			value = EXCLUDED.value,
			enabled = EXCLUDED.enabled,
			updated_at = EXCLUDED.updated_at
	`
	
	now := time.Now()
	flagValue.ID = uuid.New()
	flagValue.CreatedAt = now
	flagValue.UpdatedAt = now
	
	_, err := r.db.ExecContext(ctx, query, 
		flagValue.ID, flagValue.FlagID, flagValue.EnvID, flagValue.Value, flagValue.Enabled, flagValue.CreatedAt, flagValue.UpdatedAt)
	return err
}

func (r *flagValueRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.FlagValue, error) {
	query := `
		SELECT id, flag_id, env_id, value, enabled, created_at, updated_at
		FROM flag_values
		WHERE id = $1
	`
	
	var flagValue model.FlagValue
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&flagValue.ID,
		&flagValue.FlagID,
		&flagValue.EnvID,
		&flagValue.Value,
		&flagValue.Enabled,
		&flagValue.CreatedAt,
		&flagValue.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &flagValue, err
}

func (r *flagValueRepository) GetByFlagID(ctx context.Context, flagID uuid.UUID) ([]model.FlagValue, error) {
	query := `
		SELECT id, flag_id, env_id, value, enabled, created_at, updated_at
		FROM flag_values
		WHERE flag_id = $1
		ORDER BY created_at ASC
	`
	
	rows, err := r.db.QueryContext(ctx, query, flagID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var flagValues []model.FlagValue
	for rows.Next() {
		var flagValue model.FlagValue
		err := rows.Scan(
			&flagValue.ID,
			&flagValue.FlagID,
			&flagValue.EnvID,
			&flagValue.Value,
			&flagValue.Enabled,
			&flagValue.CreatedAt,
			&flagValue.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		flagValues = append(flagValues, flagValue)
	}
	
	return flagValues, nil
}

func (r *flagValueRepository) GetByEnvID(ctx context.Context, envID uuid.UUID) ([]model.FlagValue, error) {
	query := `
		SELECT id, flag_id, env_id, value, enabled, created_at, updated_at
		FROM flag_values
		WHERE env_id = $1
		ORDER BY created_at ASC
	`
	
	rows, err := r.db.QueryContext(ctx, query, envID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var flagValues []model.FlagValue
	for rows.Next() {
		var flagValue model.FlagValue
		err := rows.Scan(
			&flagValue.ID,
			&flagValue.FlagID,
			&flagValue.EnvID,
			&flagValue.Value,
			&flagValue.Enabled,
			&flagValue.CreatedAt,
			&flagValue.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		flagValues = append(flagValues, flagValue)
	}
	
	return flagValues, nil
}

func (r *flagValueRepository) Update(ctx context.Context, id uuid.UUID, req *dto.UpdateFlagValueRequest) (*model.FlagValue, error) {
	query := `
		UPDATE flag_values
		SET value = COALESCE($1, value),
			enabled = COALESCE($2, enabled),
			updated_at = $3
		WHERE id = $4
		RETURNING id, flag_id, env_id, value, enabled, created_at, updated_at
	`
	
	var flagValue model.FlagValue
	now := time.Now()
	
	err := r.db.QueryRowContext(ctx, query, 
		req.Value, req.Enabled, now, id).Scan(
		&flagValue.ID,
		&flagValue.FlagID,
		&flagValue.EnvID,
		&flagValue.Value,
		&flagValue.Enabled,
		&flagValue.CreatedAt,
		&flagValue.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &flagValue, err
}

func (r *flagValueRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM flag_values WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
