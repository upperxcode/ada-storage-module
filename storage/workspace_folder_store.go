package storage

import (
	"context"
	"database/sql"
	"errors"
)

// ErrWorkspaceFolderNotFound is returned when a workspace folder does not exist.
var ErrWorkspaceFolderNotFound = errors.New("workspace folder not found")

// WorkspaceFolder represents a folder associated with a workspace.
type WorkspaceFolder struct {
	ID         int64  `json:"id"`
	WorkspaceID int64 `json:"workspace_id"`
	FolderPath string `json:"folder_path"`
}

// WorkspaceFolderStore handles persistence operations for workspace folders.
type WorkspaceFolderStore struct {
	db *sql.DB
}

// NewWorkspaceFolderStore creates a new WorkspaceFolderStore instance.
func NewWorkspaceFolderStore(db *sql.DB) *WorkspaceFolderStore {
	return &WorkspaceFolderStore{db: db}
}

// Create adds a new folder to a workspace.
func (s *WorkspaceFolderStore) Create(ctx context.Context, folder *WorkspaceFolder) error {
	query := `INSERT INTO workspace_folders (workspace_id, folder_path) VALUES (?, ?)`
	_, err := s.db.ExecContext(ctx, query, folder.WorkspaceID, folder.FolderPath)
	if err != nil {
		return err
	}
	return nil
}

// Delete removes a specific folder from a workspace.
func (s *WorkspaceFolderStore) Delete(ctx context.Context, workspaceID int64, folderPath string) error {
	query := `DELETE FROM workspace_folders WHERE workspace_id = ? AND folder_path = ?`
	result, err := s.db.ExecContext(ctx, query, workspaceID, folderPath)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrWorkspaceFolderNotFound
	}
	return nil
}

// ListByWorkspace retrieves all folders for a workspace.
func (s *WorkspaceFolderStore) ListByWorkspace(ctx context.Context, workspaceID int64) ([]WorkspaceFolder, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, workspace_id, folder_path FROM workspace_folders WHERE workspace_id = ?`,
		workspaceID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var folders []WorkspaceFolder
	for rows.Next() {
		var f WorkspaceFolder
		if err := rows.Scan(&f.ID, &f.WorkspaceID, &f.FolderPath); err != nil {
			return nil, err
		}
		folders = append(folders, f)
	}
	return folders, nil
}

// DeleteAllByWorkspace removes all folders from a workspace (used when deleting a workspace).
func (s *WorkspaceFolderStore) DeleteAllByWorkspace(ctx context.Context, workspaceID int64) error {
	query := `DELETE FROM workspace_folders WHERE workspace_id = ?`
	_, err := s.db.ExecContext(ctx, query, workspaceID)
	if err != nil {
		return err
	}
	return nil
}