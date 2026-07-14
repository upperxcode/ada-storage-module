package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// ErrFixedModelNotFound is returned when a fixed model does not exist.
var ErrFixedModelNotFound = errors.New("fixed model not found")

// FixedModel represents a fixed model configuration.
type FixedModel struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Provider  string `json:"provider"`
	Model     string `json:"model"`
}

// FixedModelTool represents a tool associated with a fixed model.
type FixedModelTool struct {
	ID          int64  `json:"id"`
	FixedModelID int64 `json:"fixed_model_id"`
	Tool        string `json:"tool"`
}

// FixedModelStore handles persistence operations for fixed models and their tools.
type FixedModelStore struct {
	db *sql.DB
}

// NewFixedModelStore creates a new FixedModelStore instance.
func NewFixedModelStore(db *sql.DB) *FixedModelStore {
	return &FixedModelStore{db: db}
}

// CreateFixedModel creates a new fixed model.
func (s *FixedModelStore) CreateFixedModel(ctx context.Context, model *FixedModel) error {
	query := `INSERT INTO fixed_models (name, provider, model) VALUES (?, ?, ?)`
	_, err := s.db.ExecContext(ctx, query, model.Name, model.Provider, model.Model)
	if err != nil {
		return fmt.Errorf("failed to create fixed model: %w", err)
	}
	return nil
}

// GetFixedModel retrieves a fixed model by name.
func (s *FixedModelStore) GetFixedModel(ctx context.Context, name string) (*FixedModel, error) {
	query := `SELECT id, name, provider, model FROM fixed_models WHERE name = ?`
	var model FixedModel
	err := s.db.QueryRowContext(ctx, query, name).Scan(&model.ID, &model.Name, &model.Provider, &model.Model)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrFixedModelNotFound
		}
		return nil, fmt.Errorf("failed to get fixed model: %w", err)
	}
	return &model, nil
}

// ListFixedModels retrieves all fixed models.
func (s *FixedModelStore) ListFixedModels(ctx context.Context) ([]FixedModel, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, name, provider, model FROM fixed_models`)
	if err != nil {
		return nil, fmt.Errorf("failed to query fixed models: %w", err)
	}
	defer rows.Close()

	var models []FixedModel
	for rows.Next() {
		var m FixedModel
		if err := rows.Scan(&m.ID, &m.Name, &m.Provider, &m.Model); err != nil {
			return nil, fmt.Errorf("failed to scan fixed model: %w", err)
		}
		models = append(models, m)
	}
	return models, nil
}

// DeleteFixedModel deletes a fixed model and its tools (CASCADE).
func (s *FixedModelStore) DeleteFixedModel(ctx context.Context, id int64) error {
	query := `DELETE FROM fixed_models WHERE id = ?`
	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete fixed model: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return ErrFixedModelNotFound
	}
	return nil
}

// AddTool adds a tool to a fixed model.
func (s *FixedModelStore) AddTool(ctx context.Context, tool *FixedModelTool) error {
	query := `INSERT INTO fixed_model_tools (fixed_model_id, tool) VALUES (?, ?)`
	_, err := s.db.ExecContext(ctx, query, tool.FixedModelID, tool.Tool)
	if err != nil {
		return fmt.Errorf("failed to add tool: %w", err)
	}
	return nil
}

// ListTools retrieves all tools for a fixed model.
func (s *FixedModelStore) ListTools(ctx context.Context, modelID int64) ([]FixedModelTool, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, fixed_model_id, tool FROM fixed_model_tools WHERE fixed_model_id = ?`,
		modelID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query tools: %w", err)
	}
	defer rows.Close()

	var tools []FixedModelTool
	for rows.Next() {
		var t FixedModelTool
		if err := rows.Scan(&t.ID, &t.FixedModelID, &t.Tool); err != nil {
			return nil, fmt.Errorf("failed to scan tool: %w", err)
		}
		tools = append(tools, t)
	}
	return tools, nil
}