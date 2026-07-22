package storage

import (
	"context"
	"database/sql"
	"fmt"
)

// ProviderModelStore handles persistence operations for provider models.
type ProviderModelStore struct {
	db *sql.DB
}

// NewProviderModelStore creates a new ProviderModelStore instance.
func NewProviderModelStore(db *sql.DB) *ProviderModelStore {
	return &ProviderModelStore{db: db}
}

// CreateProviderModel creates a new provider model.
func (s *ProviderModelStore) CreateProviderModel(ctx context.Context, model *ProviderModel) error {
	query := `INSERT INTO provider_models 
		(provider_id, model, free, thinking, tool, embedding, vision, health, context_size, max_tokens)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := s.db.ExecContext(ctx, query,
		model.ProviderID, model.Model, model.Free, model.Thinking, model.Tool,
		model.Embedding, model.Vision, model.Health, model.ContextSize, model.MaxTokens,
	)
	if err != nil {
		return fmt.Errorf("failed to create provider model: %w", err)
	}
	return nil
}

// GetProviderModels retrieves all models for a provider.
func (s *ProviderModelStore) GetProviderModels(ctx context.Context, providerID int64) ([]ProviderModel, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, provider_id, model, free, thinking, tool, embedding, vision, health, context_size, max_tokens
		 FROM provider_models WHERE provider_id = ? ORDER BY model ASC`,
		providerID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query provider models: %w", err)
	}
	defer rows.Close()

	var models []ProviderModel
	for rows.Next() {
		var m ProviderModel
		var freeInt, thinkingInt, toolInt, embeddingInt, visionInt int
		if err := rows.Scan(&m.ID, &m.ProviderID, &m.Model, &freeInt, &thinkingInt,
			&toolInt, &embeddingInt, &visionInt, &m.Health, &m.ContextSize, &m.MaxTokens); err != nil {
			return nil, fmt.Errorf("failed to scan provider model: %w", err)
		}
		m.Free = freeInt == 1
		m.Thinking = thinkingInt == 1
		m.Tool = toolInt == 1
		m.Embedding = embeddingInt == 1
		m.Vision = visionInt == 1
		models = append(models, m)
	}

	return models, nil
}

// UpdateModelHealth updates the health score of a specific provider model.
func (s *ProviderModelStore) UpdateModelHealth(ctx context.Context, providerID int64, model string, health int) error {
	query := `UPDATE provider_models SET health = ? WHERE provider_id = ? AND model = ?`
	result, err := s.db.ExecContext(ctx, query, health, providerID, model)
	if err != nil {
		return fmt.Errorf("failed to update model health: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return ErrProviderModelNotFound
	}
	return nil
}

// DeleteProviderModel deletes a model.
func (s *ProviderModelStore) DeleteProviderModel(ctx context.Context, id int64) error {
	query := `DELETE FROM provider_models WHERE id = ?`
	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete provider model: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return ErrProviderModelNotFound
	}

	return nil
}

// GetEnabledModelsByProvider returns all enabled models grouped by provider name.
func (s *ProviderModelStore) GetEnabledModelsByProvider(ctx context.Context, db *sql.DB) (map[string][]ProviderModel, error) {
	query := `SELECT p.name, pm.id, pm.provider_id, pm.model, pm.free, pm.thinking, 
			   pm.tool, pm.embedding, pm.vision, pm.health, pm.context_size, pm.max_tokens
			FROM provider_models pm
			JOIN providers p ON pm.provider_id = p.id
			WHERE pm.health > 0
			ORDER BY p.name, pm.model`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query enabled models: %w", err)
	}
	defer rows.Close()

	grouped := make(map[string][]ProviderModel)
	for rows.Next() {
		var providerName string
		var m ProviderModel
		var freeInt, thinkingInt, toolInt, embeddingInt, visionInt int

		if err := rows.Scan(&providerName, &m.ID, &m.ProviderID, &m.Model,
			&freeInt, &thinkingInt, &toolInt, &embeddingInt, &visionInt, &m.Health,
			&m.ContextSize, &m.MaxTokens); err != nil {
			return nil, fmt.Errorf("failed to scan model: %w", err)
		}

		m.Free = freeInt == 1
		m.Thinking = thinkingInt == 1
		m.Tool = toolInt == 1
		m.Embedding = embeddingInt == 1
		m.Vision = visionInt == 1

		grouped[providerName] = append(grouped[providerName], m)
	}

	return grouped, nil
}