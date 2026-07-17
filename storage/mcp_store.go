package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// ErrMcpNotFound is returned when an MCP does not exist.
var ErrMcpNotFound = errors.New("mcp not found")

// Mcp represents an MCP (Model Context Protocol) configuration.
type Mcp struct {
	ID            int64          `json:"id"`
	Name          string         `json:"nome"`
	ConnectType   string         `json:"connect_type"`
	Command       sql.NullString `json:"command"`
	Arguments     sql.NullString `json:"arguments"`
	Environment   sql.NullString `json:"environment"`
	URL           sql.NullString `json:"url"`
	Enabled       bool           `json:"enabled"`
	Timeout       int            `json:"timeout"`
	OAuthClientID string         `json:"oauth_client_id"`
	Color         string         `json:"color"`
	Icon          string         `json:"icon"`
}

// McpStore handles persistence operations for MCPs.
type McpStore struct {
	db *sql.DB
}

// NewMcpStore creates a new McpStore instance.
func NewMcpStore(db *sql.DB) *McpStore {
	return &McpStore{db: db}
}

// CreateMcp creates a new MCP.
func (s *McpStore) CreateMcp(ctx context.Context, mcp *Mcp) error {
	query := `INSERT INTO mcps (nome, connect_type, command, arguments, environment, url, enabled, timeout, oauth_client_id, color, icon) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	result, err := s.db.ExecContext(ctx, query, mcp.Name, mcp.ConnectType, mcp.Command, mcp.Arguments, mcp.Environment, mcp.URL, mcp.Enabled, mcp.Timeout, mcp.OAuthClientID, mcp.Color, mcp.Icon)
	if err != nil {
		return fmt.Errorf("failed to create MCP: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}
	mcp.ID = id
	return nil
}

// UpdateMcp updates an existing MCP by name.
func (s *McpStore) UpdateMcp(ctx context.Context, mcp *Mcp) error {
	query := `UPDATE mcps SET nome=?, connect_type=?, command=?, arguments=?, environment=?, url=?, enabled=?, timeout=?, oauth_client_id=?, color=?, icon=? WHERE id=?`
	result, err := s.db.ExecContext(ctx, query, mcp.Name, mcp.ConnectType, mcp.Command, mcp.Arguments, mcp.Environment, mcp.URL, mcp.Enabled, mcp.Timeout, mcp.OAuthClientID, mcp.Color, mcp.Icon, mcp.ID)
	if err != nil {
		return fmt.Errorf("failed to update MCP: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return ErrMcpNotFound
	}
	return nil
}

// GetMcp retrieves an MCP by ID.
func (s *McpStore) GetMcp(ctx context.Context, id int64) (*Mcp, error) {
	query := `SELECT id, nome, connect_type, command, arguments, environment, url, enabled, timeout, oauth_client_id, color, icon FROM mcps WHERE id = ?`
	var m Mcp
	err := s.db.QueryRowContext(ctx, query, id).Scan(&m.ID, &m.Name, &m.ConnectType, &m.Command, &m.Arguments, &m.Environment, &m.URL, &m.Enabled, &m.Timeout, &m.OAuthClientID, &m.Color, &m.Icon)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrMcpNotFound
		}
		return nil, fmt.Errorf("failed to get MCP: %w", err)
	}
	return &m, nil
}

// GetMcpByName retrieves an MCP by name.
func (s *McpStore) GetMcpByName(ctx context.Context, name string) (*Mcp, error) {
	query := `SELECT id, nome, connect_type, command, arguments, environment, url, enabled, timeout, oauth_client_id, color, icon FROM mcps WHERE nome = ?`
	var m Mcp
	err := s.db.QueryRowContext(ctx, query, name).Scan(&m.ID, &m.Name, &m.ConnectType, &m.Command, &m.Arguments, &m.Environment, &m.URL, &m.Enabled, &m.Timeout, &m.OAuthClientID, &m.Color, &m.Icon)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrMcpNotFound
		}
		return nil, fmt.Errorf("failed to get MCP: %w", err)
	}
	return &m, nil
}

// ListMcps retrieves all MCPs.
func (s *McpStore) ListMcps(ctx context.Context) ([]Mcp, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, nome, connect_type, command, arguments, environment, url, enabled, timeout, oauth_client_id, color, icon FROM mcps ORDER BY nome`)
	if err != nil {
		return nil, fmt.Errorf("failed to query MCPs: %w", err)
	}
	defer rows.Close()

	var mcps []Mcp
	for rows.Next() {
		var m Mcp
		if err := rows.Scan(&m.ID, &m.Name, &m.ConnectType, &m.Command, &m.Arguments, &m.Environment, &m.URL, &m.Enabled, &m.Timeout, &m.OAuthClientID, &m.Color, &m.Icon); err != nil {
			return nil, fmt.Errorf("failed to scan MCP: %w", err)
		}
		mcps = append(mcps, m)
	}
	return mcps, nil
}

// DeleteMcp deletes an MCP by ID.
func (s *McpStore) DeleteMcp(ctx context.Context, id int64) error {
	query := `DELETE FROM mcps WHERE id = ?`
	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete MCP: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return ErrMcpNotFound
	}
	return nil
}

// DeleteMcpByName deletes an MCP by name.
func (s *McpStore) DeleteMcpByName(ctx context.Context, name string) error {
	query := `DELETE FROM mcps WHERE nome = ?`
	result, err := s.db.ExecContext(ctx, query, name)
	if err != nil {
		return fmt.Errorf("failed to delete MCP by name: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return ErrMcpNotFound
	}
	return nil
}