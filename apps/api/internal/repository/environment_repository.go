package repository

import (
	"context"
	"database/sql"
	"time"

	"api/internal/dto"
	"api/internal/model"

	"github.com/google/uuid"
)

type environmentRepository struct {
	db *sql.DB
}

func NewEnvironmentRepository(db *sql.DB) EnvironmentRepository {
	return &environmentRepository{db: db}
}

func (r *environmentRepository) Create(ctx context.Context, env *model.Environment) error {
	query := `
		INSERT INTO environments (id, project_id, name, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	
	now := time.Now()
	env.ID = uuid.New()
	env.CreatedAt = now
	env.UpdatedAt = now
	
	_, err := r.db.ExecContext(ctx, query, env.ID, env.ProjectID, env.Name, env.CreatedAt, env.UpdatedAt)
	return err
}

func (r *environmentRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Environment, error) {
	query := `
		SELECT id, project_id, name, created_at, updated_at
		FROM environments
		WHERE id = $1
	`
	
	var env model.Environment
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&env.ID,
		&env.ProjectID,
		&env.Name,
		&env.CreatedAt,
		&env.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &env, err
}

func (r *environmentRepository) GetAll(ctx context.Context) ([]model.Environment, error) {
	query := `
		SELECT id, project_id, name, created_at, updated_at
		FROM environments
		ORDER BY created_at DESC
	`
	
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var envs []model.Environment
	for rows.Next() {
		var env model.Environment
		err := rows.Scan(
			&env.ID,
			&env.ProjectID,
			&env.Name,
			&env.CreatedAt,
			&env.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		envs = append(envs, env)
	}
	
	return envs, nil
}

func (r *environmentRepository) GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]model.Environment, error) {
	query := `
		SELECT id, project_id, name, created_at, updated_at
		FROM environments
		WHERE project_id = $1
		ORDER BY created_at DESC
	`
	
	rows, err := r.db.QueryContext(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var envs []model.Environment
	for rows.Next() {
		var env model.Environment
		err := rows.Scan(
			&env.ID,
			&env.ProjectID,
			&env.Name,
			&env.CreatedAt,
			&env.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		envs = append(envs, env)
	}
	
	return envs, nil
}

func (r *environmentRepository) Update(ctx context.Context, id uuid.UUID, req *dto.UpdateEnvironmentRequest) (*model.Environment, error) {
	query := `
		UPDATE environments
		SET name = COALESCE($1, name),
			updated_at = $2
		WHERE id = $3
		RETURNING id, project_id, name, created_at, updated_at
	`
	
	var env model.Environment
	now := time.Now()
	
	err := r.db.QueryRowContext(ctx, query, 
		req.Name, now, id).Scan(
		&env.ID,
		&env.ProjectID,
		&env.Name,
		&env.CreatedAt,
		&env.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &env, err
}

func (r *environmentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM environments WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
