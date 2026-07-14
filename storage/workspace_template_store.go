package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// ErrWorkspaceTemplateNotFound is returned when a workspace template does not exist.
var ErrWorkspaceTemplateNotFound = errors.New("workspace template not found")

// WorkspaceTemplate represents a workspace template.
type WorkspaceTemplate struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description sql.NullString `json:"description"`
	Personality string    `json:"personality"`
	CreatedAt   time.Time `json:"created_at"`
}

// WorkspaceTemplateStore handles persistence operations for workspace templates.
type WorkspaceTemplateStore struct {
	db *sql.DB
}

// NewWorkspaceTemplateStore creates a new WorkspaceTemplateStore instance.
func NewWorkspaceTemplateStore(db *sql.DB) *WorkspaceTemplateStore {
	return &WorkspaceTemplateStore{db: db}
}

// CreateTemplate creates a new workspace template.
func (s *WorkspaceTemplateStore) CreateTemplate(ctx context.Context, template *WorkspaceTemplate) error {
	query := `INSERT INTO workspace_templates (name, description, personality, created_at) VALUES (?, ?, ?, ?)`
	_, err := s.db.ExecContext(ctx, query, template.Name, template.Description, template.Personality, time.Now().UTC())
	if err != nil {
		return fmt.Errorf("failed to create workspace template: %w", err)
	}
	return nil
}

// GetTemplate retrieves a workspace template by ID.
func (s *WorkspaceTemplateStore) GetTemplate(ctx context.Context, id int64) (*WorkspaceTemplate, error) {
	query := `SELECT id, name, description, personality, created_at FROM workspace_templates WHERE id = ?`
	var t WorkspaceTemplate
	err := s.db.QueryRowContext(ctx, query, id).Scan(&t.ID, &t.Name, &t.Description, &t.Personality, &t.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrWorkspaceTemplateNotFound
		}
		return nil, fmt.Errorf("failed to get workspace template: %w", err)
	}
	return &t, nil
}

// ListTemplates retrieves all workspace templates.
func (s *WorkspaceTemplateStore) ListTemplates(ctx context.Context) ([]WorkspaceTemplate, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, name, description, personality, created_at FROM workspace_templates ORDER BY name`)
	if err != nil {
		return nil, fmt.Errorf("failed to query workspace templates: %w", err)
	}
	defer rows.Close()

	var templates []WorkspaceTemplate
	for rows.Next() {
		var t WorkspaceTemplate
		if err := rows.Scan(&t.ID, &t.Name, &t.Description, &t.Personality, &t.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan workspace template: %w", err)
		}
		templates = append(templates, t)
	}
	return templates, nil
}

// DeleteTemplate deletes a workspace template.
func (s *WorkspaceTemplateStore) DeleteTemplate(ctx context.Context, id int64) error {
	query := `DELETE FROM workspace_templates WHERE id = ?`
	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete workspace template: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return ErrWorkspaceTemplateNotFound
	}
	return nil
}