package storage

import (
	"context"
	"database/sql"
	"fmt"
)

// WorkspaceWorker represents a workspace-worker association.
type WorkspaceWorker struct {
	ID        int64 `json:"id"`
	WorkspaceID int64 `json:"workspace_id"`
	WorkerID  int64 `json:"worker_id"`
	Enabled   bool  `json:"enabled"`
}

// WorkspaceWorkersStore handles persistence operations for workspace-worker associations.
type WorkspaceWorkersStore struct {
	db *sql.DB
}

// NewWorkspaceWorkersStore creates a new WorkspaceWorkersStore instance.
func NewWorkspaceWorkersStore(db *sql.DB) *WorkspaceWorkersStore {
	return &WorkspaceWorkersStore{db: db}
}

// AddWorker adds a worker to a workspace.
func (s *WorkspaceWorkersStore) AddWorker(ctx context.Context, link *WorkspaceWorker) error {
	query := `INSERT INTO workspace_workers (workspace_id, worker_id, enabled) VALUES (?, ?, ?)`
	_, err := s.db.ExecContext(ctx, query, link.WorkspaceID, link.WorkerID, link.Enabled)
	if err != nil {
		return fmt.Errorf("failed to add workspace worker: %w", err)
	}
	return nil
}

// ListWorkers retrieves all workers for a workspace.
func (s *WorkspaceWorkersStore) ListWorkers(ctx context.Context, workspaceID int64) ([]WorkspaceWorker, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, workspace_id, worker_id, enabled FROM workspace_workers WHERE workspace_id = ?`,
		workspaceID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query workspace workers: %w", err)
	}
	defer rows.Close()

	var links []WorkspaceWorker
	for rows.Next() {
		var link WorkspaceWorker
		var enabledInt int
		if err := rows.Scan(&link.ID, &link.WorkspaceID, &link.WorkerID, &enabledInt); err != nil {
			return nil, fmt.Errorf("failed to scan workspace worker: %w", err)
		}
		link.Enabled = enabledInt == 1
		links = append(links, link)
	}
	return links, nil
}

// RemoveWorker removes a worker from a workspace.
func (s *WorkspaceWorkersStore) RemoveWorker(ctx context.Context, workspaceID, workerID int64) error {
	query := `DELETE FROM workspace_workers WHERE workspace_id = ? AND worker_id = ?`
	result, err := s.db.ExecContext(ctx, query, workspaceID, workerID)
	if err != nil {
		return fmt.Errorf("failed to remove workspace worker: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("workspace worker not found")
	}
	return nil
}

// DeleteAllWorkers deletes all workers for a workspace.
func (s *WorkspaceWorkersStore) DeleteAllWorkers(ctx context.Context, workspaceID int64) error {
	query := `DELETE FROM workspace_workers WHERE workspace_id = ?`
	_, err := s.db.ExecContext(ctx, query, workspaceID)
	if err != nil {
		return fmt.Errorf("failed to delete all workspace workers: %w", err)
	}
	return nil
}