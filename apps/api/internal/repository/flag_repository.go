package repository

import (
	"context"
	"database/sql"
	"time"

	"api/internal/dto"
	"api/internal/model"

	"github.com/google/uuid"
)

type flagRepository struct {
	db *sql.DB
}

func NewFlagRepository(db *sql.DB) FlagRepository {
	return &flagRepository{db: db}
}

func (r *flagRepository) Create(ctx context.Context, flag *model.Flag) error {
	query := `
		INSERT INTO flags (id, project_id, key, description, type, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	
	now := time.Now()
	flag.ID = uuid.New()
	flag.CreatedAt = now
	flag.UpdatedAt = now
	
	_, err := r.db.ExecContext(ctx, query, flag.ID, flag.ProjectID, flag.Key, flag.Description, flag.Type, flag.CreatedAt, flag.UpdatedAt)
	return err
}

func (r *flagRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Flag, error) {
	query := `
		SELECT id, project_id, key, description, type, created_at, updated_at
		FROM flags
		WHERE id = $1
	`
	
	var flag model.Flag
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&flag.ID,
		&flag.ProjectID,
		&flag.Key,
		&flag.Description,
		&flag.Type,
		&flag.CreatedAt,
		&flag.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &flag, err
}

func (r *flagRepository) GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]model.Flag, error) {
	query := `
		SELECT id, project_id, key, description, type, created_at, updated_at
		FROM flags
		WHERE project_id = $1
		ORDER BY created_at DESC
	`
	
	rows, err := r.db.QueryContext(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var flags []model.Flag
	for rows.Next() {
		var flag model.Flag
		err := rows.Scan(
			&flag.ID,
			&flag.ProjectID,
			&flag.Key,
			&flag.Description,
			&flag.Type,
			&flag.CreatedAt,
			&flag.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		flags = append(flags, flag)
	}
	
	return flags, nil
}

func (r *flagRepository) GetWithValuesByProjectID(ctx context.Context, projectID uuid.UUID) ([]model.FlagWithValues, error) {
	flagQuery := `
		SELECT id, project_id, key, description, type, created_at, updated_at
		FROM flags
		WHERE project_id = $1
		ORDER BY created_at DESC
	`
	
	flagRows, err := r.db.QueryContext(ctx, flagQuery, projectID)
	if err != nil {
		return nil, err
	}
	defer flagRows.Close()
	
	var flagsWithValues []model.FlagWithValues
	for flagRows.Next() {
		var flag model.Flag
		err := flagRows.Scan(
			&flag.ID,
			&flag.ProjectID,
			&flag.Key,
			&flag.Description,
			&flag.Type,
			&flag.CreatedAt,
			&flag.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		
		// Get values for this flag
		valueQuery := `
			SELECT id, flag_id, env_id, value, enabled, created_at, updated_at
			FROM flag_values
			WHERE flag_id = $1
			ORDER BY created_at ASC
		`
		
		valueRows, err := r.db.QueryContext(ctx, valueQuery, flag.ID)
		if err != nil {
			return nil, err
		}
		
		var values []model.FlagValue
		for valueRows.Next() {
			var value model.FlagValue
			err := valueRows.Scan(
				&value.ID,
				&value.FlagID,
				&value.EnvID,
				&value.Value,
				&value.Enabled,
				&value.CreatedAt,
				&value.UpdatedAt,
			)
			if err != nil {
				valueRows.Close()
				return nil, err
			}
			values = append(values, value)
		}
		valueRows.Close()
		
		flagsWithValues = append(flagsWithValues, model.FlagWithValues{
			Flag:   flag,
			Values: values,
		})
	}
	
	return flagsWithValues, nil
}

func (r *flagRepository) Update(ctx context.Context, id uuid.UUID, req *dto.UpdateFlagRequest) (*model.Flag, error) {
	query := `
		UPDATE flags
		SET key = COALESCE($1, key),
			description = COALESCE($2, description),
			type = COALESCE($3, type),
			updated_at = $4
		WHERE id = $5
		RETURNING id, project_id, key, description, type, created_at, updated_at
	`
	
	var flag model.Flag
	now := time.Now()
	
	err := r.db.QueryRowContext(ctx, query, 
		req.Key, req.Description, req.Type, now, id).Scan(
		&flag.ID,
		&flag.ProjectID,
		&flag.Key,
		&flag.Description,
		&flag.Type,
		&flag.CreatedAt,
		&flag.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &flag, err
}

func (r *flagRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM flags WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
