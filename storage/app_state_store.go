package storage

import (
	"context"
	"database/sql"
	"fmt"
)

// AppState represents the application state (singleton row with id=1).
type AppState struct {
	ID                     int64  `json:"id"`
	ActiveWorkspacePath    string `json:"active_workspace_path"`
	ActiveWorkspaceIndex   int    `json:"active_workspace_index"`
}

// AppStateStore handles persistence operations for application state.
type AppStateStore struct {
	db *sql.DB
}

// NewAppStateStore creates a new AppStateStore instance.
func NewAppStateStore(db *sql.DB) *AppStateStore {
	return &AppStateStore{db: db}
}

// GetState retrieves the application state.
func (s *AppStateStore) GetState(ctx context.Context) (*AppState, error) {
	query := `SELECT id, active_workspace_path, active_workspace_index FROM app_state WHERE id = 1`
	var state AppState
	err := s.db.QueryRowContext(ctx, query).Scan(&state.ID, &state.ActiveWorkspacePath, &state.ActiveWorkspaceIndex)
	if err != nil {
		if err == sql.ErrNoRows {
			// Initialize with default state
			state = AppState{ID: 1, ActiveWorkspacePath: "", ActiveWorkspaceIndex: 0}
			s.db.ExecContext(ctx, `INSERT INTO app_state (id, active_workspace_path, active_workspace_index) VALUES (1, '', 0)`)
		} else {
			return nil, fmt.Errorf("failed to get app state: %w", err)
		}
	}
	return &state, nil
}

// SetActiveWorkspace sets the active workspace.
func (s *AppStateStore) SetActiveWorkspace(ctx context.Context, path string, index int) error {
	query := `INSERT INTO app_state (id, active_workspace_path, active_workspace_index) VALUES (1, ?, ?)
			  ON CONFLICT(id) DO UPDATE SET active_workspace_path = excluded.active_workspace_path, active_workspace_index = excluded.active_workspace_index`
	_, err := s.db.ExecContext(ctx, query, path, index)
	if err != nil {
		return fmt.Errorf("failed to set active workspace: %w", err)
	}
	return nil
}