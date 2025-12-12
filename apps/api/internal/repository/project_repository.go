package repository

import (
	"context"
	"database/sql"
	"time"

	"api/internal/dto"
	"api/internal/model"

	"github.com/google/uuid"
)

type projectRepository struct {
	db *sql.DB
}

func NewProjectRepository(db *sql.DB) ProjectRepository {
	return &projectRepository{db: db}
}

func (r *projectRepository) Create(ctx context.Context, project *model.Project) error {
	query := `
		INSERT INTO projects (id, name, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	now := time.Now()
	project.ID = uuid.New()
	project.CreatedAt = now
	project.UpdatedAt = now

	_, err := r.db.ExecContext(ctx, query, project.ID, project.Name, project.Description, project.CreatedAt, project.UpdatedAt)
	return err
}

func (r *projectRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Project, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM projects
		WHERE id = $1
	`

	var project model.Project
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&project.ID,
		&project.Name,
		&project.Description,
		&project.CreatedAt,
		&project.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &project, err
}

func (r *projectRepository) GetAll(ctx context.Context) ([]model.Project, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM projects
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []model.Project
	for rows.Next() {
		var project model.Project
		err := rows.Scan(
			&project.ID,
			&project.Name,
			&project.Description,
			&project.CreatedAt,
			&project.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}

	return projects, nil
}

func (r *projectRepository) Update(ctx context.Context, id uuid.UUID, req *dto.UpdateProjectRequest) (*model.Project, error) {
	query := `
		UPDATE projects
		SET name = COALESCE($1, name),
			description = COALESCE($2, description),
			updated_at = $3
		WHERE id = $4
		RETURNING id, name, description, created_at, updated_at
	`

	var project model.Project
	now := time.Now()

	err := r.db.QueryRowContext(ctx, query,
		req.Name, req.Description, now, id).Scan(
		&project.ID,
		&project.Name,
		&project.Description,
		&project.CreatedAt,
		&project.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &project, err
}

func (r *projectRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM projects WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
