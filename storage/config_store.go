package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// ConfigStore handles persistence operations for generic configurations and API keys.
type ConfigStore struct {
	db *sql.DB
}

// NewConfigStore creates a new ConfigStore instance.
func NewConfigStore(db *sql.DB) *ConfigStore {
	return &ConfigStore{db: db}
}

// SetConfig stores a configuration key-value pair.
func (s *ConfigStore) SetConfig(ctx context.Context, key string, value string) error {
	query := `INSERT INTO config (key, value)
			  VALUES (?, ?)
			  ON CONFLICT(key) DO UPDATE SET value = excluded.value`

	_, err := s.db.ExecContext(ctx, query, key, value)
	if err != nil {
		return fmt.Errorf("failed to set config: %w", err)
	}
	return nil
}

// GetConfig retrieves a configuration value by key.
func (s *ConfigStore) GetConfig(ctx context.Context, key string) (string, error) {
	var value string
	err := s.db.QueryRowContext(ctx,
		`SELECT value FROM config WHERE key = ?`,
		key,
	).Scan(&value)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrConfigNotFound
		}
		return "", fmt.Errorf("failed to get config: %w", err)
	}

	return value, nil
}

// DeleteConfig removes a configuration key.
func (s *ConfigStore) DeleteConfig(ctx context.Context, key string) error {
	query := `DELETE FROM config WHERE key = ?`
	result, err := s.db.ExecContext(ctx, query, key)
	if err != nil {
		return fmt.Errorf("failed to delete config: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return ErrConfigNotFound
	}

	return nil
}

// ListConfigs retrieves all configuration key-value pairs.
func (s *ConfigStore) ListConfigs(ctx context.Context) (map[string]string, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT key, value FROM config`,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query configs: %w", err)
	}
	defer rows.Close()

	configs := make(map[string]string)
	for rows.Next() {
		var key, value string
		if err := rows.Scan(&key, &value); err != nil {
			return nil, fmt.Errorf("failed to scan config: %w", err)
		}
		configs[key] = value
	}

	return configs, nil
}

// SetAPIKey stores an API key for a provider.
func (s *ConfigStore) SetAPIKey(ctx context.Context, providerID int64, apiKey string) error {
	query := `INSERT INTO provider_apikeys (provider_id, apikey)
			  VALUES (?, ?)
			  ON CONFLICT(provider_id, apikey) DO NOTHING`

	_, err := s.db.ExecContext(ctx, query, providerID, apiKey)
	if err != nil {
		return fmt.Errorf("failed to set API key: %w", err)
	}
	return nil
}

// GetAPIKeys retrieves all API keys for a provider.
func (s *ConfigStore) GetAPIKeys(ctx context.Context, providerID int64) ([]string, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT apikey FROM provider_apikeys WHERE provider_id = ?`,
		providerID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query API keys: %w", err)
	}
	defer rows.Close()

	var apikeys []string
	for rows.Next() {
		var key string
		if err := rows.Scan(&key); err != nil {
			return nil, fmt.Errorf("failed to scan API key: %w", err)
		}
		apikeys = append(apikeys, key)
	}

	return apikeys, nil
}

// DeleteAPIKey removes an API key for a provider.
func (s *ConfigStore) DeleteAPIKey(ctx context.Context, providerID int64, apiKey string) error {
	query := `DELETE FROM provider_apikeys WHERE provider_id = ? AND apikey = ?`
	result, err := s.db.ExecContext(ctx, query, providerID, apiKey)
	if err != nil {
		return fmt.Errorf("failed to delete API key: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("API key not found")
	}

	return nil
}