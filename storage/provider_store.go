package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// ProviderStore handles persistence operations for providers.
type ProviderStore struct {
	db *sql.DB
}

// NewProviderStore creates a new ProviderStore instance.
func NewProviderStore(db *sql.DB) *ProviderStore {
	return &ProviderStore{db: db}
}

// CreateProvider creates a new provider.
func (s *ProviderStore) CreateProvider(ctx context.Context, provider *Provider) error {
	query := `INSERT INTO providers (name, api_url, connection_types, color, icon)
			  VALUES (?, ?, ?, ?, ?)`

	_, err := s.db.ExecContext(ctx, query,
		provider.Name, provider.APIURL, provider.ConnectionTypes, provider.Color, provider.Icon,
	)
	if err != nil {
		return fmt.Errorf("failed to create provider: %w", err)
	}
	return nil
}

// GetProvider retrieves a provider by ID.
func (s *ProviderStore) GetProvider(ctx context.Context, id int64) (*Provider, error) {
	query := `SELECT id, name, api_url, connection_types, color, icon FROM providers WHERE id = ?`

	var provider Provider
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&provider.ID, &provider.Name, &provider.APIURL, &provider.ConnectionTypes,
		&provider.Color, &provider.Icon,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrProviderNotFound
		}
		return nil, fmt.Errorf("failed to get provider: %w", err)
	}

	return &provider, nil
}

// GetProviderByName retrieves a provider by name.
func (s *ProviderStore) GetProviderByName(ctx context.Context, name string) (*Provider, error) {
	query := `SELECT id, name, api_url, connection_types, color, icon FROM providers WHERE name = ?`

	var provider Provider
	err := s.db.QueryRowContext(ctx, query, name).Scan(
		&provider.ID, &provider.Name, &provider.APIURL, &provider.ConnectionTypes,
		&provider.Color, &provider.Icon,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrProviderNotFound
		}
		return nil, fmt.Errorf("failed to get provider by name: %w", err)
	}

	return &provider, nil
}

// ListProviders retrieves all providers.
func (s *ProviderStore) ListProviders(ctx context.Context) ([]Provider, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, name, api_url, connection_types, color, icon FROM providers ORDER BY name ASC`,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query providers: %w", err)
	}
	defer rows.Close()

	var providers []Provider
	for rows.Next() {
		var p Provider
		if err := rows.Scan(&p.ID, &p.Name, &p.APIURL, &p.ConnectionTypes, &p.Color, &p.Icon); err != nil {
			return nil, fmt.Errorf("failed to scan provider: %w", err)
		}
		providers = append(providers, p)
	}

	return providers, nil
}

// DeleteProvider deletes a provider.
func (s *ProviderStore) DeleteProvider(ctx context.Context, id int64) error {
	query := `DELETE FROM providers WHERE id = ?`
	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete provider: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return ErrProviderNotFound
	}

	return nil
}