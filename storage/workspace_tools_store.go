package storage

import (
	"context"
	"database/sql"
	"fmt"
)

// WorkspaceTool represents a workspace-tool association.
type WorkspaceTool struct {
	ID        int64  `json:"id"`
	WorkspaceID int64 `json:"workspace_id"`
	ToolName  string `json:"tool_name"`
	Enabled   bool   `json:"enabled"`
}

// WorkspaceToolsStore handles persistence operations for workspace-tool associations.
type WorkspaceToolsStore struct {
	db *sql.DB
}

// NewWorkspaceToolsStore creates a new WorkspaceToolsStore instance.
func NewWorkspaceToolsStore(db *sql.DB) *WorkspaceToolsStore {
	return &WorkspaceToolsStore{db: db}
}

// AddTool adds a tool to a workspace.
func (s *WorkspaceToolsStore) AddTool(ctx context.Context, link *WorkspaceTool) error {
	query := `INSERT INTO workspace_tools (workspace_id, tool_name, enabled) VALUES (?, ?, ?)`
	_, err := s.db.ExecContext(ctx, query, link.WorkspaceID, link.ToolName, link.Enabled)
	if err != nil {
		return fmt.Errorf("failed to add workspace tool: %w", err)
	}
	return nil
}

// ListTools retrieves all tools for a workspace.
func (s *WorkspaceToolsStore) ListTools(ctx context.Context, workspaceID int64) ([]WorkspaceTool, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, workspace_id, tool_name, enabled FROM workspace_tools WHERE workspace_id = ?`,
		workspaceID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query workspace tools: %w", err)
	}
	defer rows.Close()

	var links []WorkspaceTool
	for rows.Next() {
		var link WorkspaceTool
		var enabledInt int
		if err := rows.Scan(&link.ID, &link.WorkspaceID, &link.ToolName, &enabledInt); err != nil {
			return nil, fmt.Errorf("failed to scan workspace tool: %w", err)
		}
		link.Enabled = enabledInt == 1
		links = append(links, link)
	}
	return links, nil
}

// RemoveTool removes a tool from a workspace.
func (s *WorkspaceToolsStore) RemoveTool(ctx context.Context, workspaceID int64, toolName string) error {
	query := `DELETE FROM workspace_tools WHERE workspace_id = ? AND tool_name = ?`
	result, err := s.db.ExecContext(ctx, query, workspaceID, toolName)
	if err != nil {
		return fmt.Errorf("failed to remove workspace tool: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("workspace tool not found")
	}
	return nil
}

// DeleteAllTools deletes all tools for a workspace.
func (s *WorkspaceToolsStore) DeleteAllTools(ctx context.Context, workspaceID int64) error {
	query := `DELETE FROM workspace_tools WHERE workspace_id = ?`
	_, err := s.db.ExecContext(ctx, query, workspaceID)
	if err != nil {
		return fmt.Errorf("failed to delete all workspace tools: %w", err)
	}
	return nil
}