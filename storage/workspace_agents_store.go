package storage

import (
	"context"
	"database/sql"
	"fmt"
)

// WorkspaceAgentsStore handles persistence operations for workspace-agent associations.
type WorkspaceAgentsStore struct {
	db *sql.DB
}

// NewWorkspaceAgentsStore creates a new WorkspaceAgentsStore instance.
func NewWorkspaceAgentsStore(db *sql.DB) *WorkspaceAgentsStore {
	return &WorkspaceAgentsStore{db: db}
}

// AddAgent adds an agent to a workspace.
func (s *WorkspaceAgentsStore) AddAgent(ctx context.Context, link *WorkspaceAgent) error {
	query := `INSERT INTO workspace_agents (workspace_id, agent_id, enabled) VALUES (?, ?, ?)`
	_, err := s.db.ExecContext(ctx, query, link.WorkspaceID, link.AgentID, link.Enabled)
	if err != nil {
		return fmt.Errorf("failed to add workspace agent: %w", err)
	}
	return nil
}

// ListAgents retrieves all agents for a workspace.
func (s *WorkspaceAgentsStore) ListAgents(ctx context.Context, workspaceID int64) ([]WorkspaceAgent, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, workspace_id, agent_id, enabled FROM workspace_agents WHERE workspace_id = ?`,
		workspaceID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query workspace agents: %w", err)
	}
	defer rows.Close()

	var links []WorkspaceAgent
	for rows.Next() {
		var link WorkspaceAgent
		var enabledInt int
		if err := rows.Scan(&link.ID, &link.WorkspaceID, &link.AgentID, &enabledInt); err != nil {
			return nil, fmt.Errorf("failed to scan workspace agent: %w", err)
		}
		link.Enabled = enabledInt == 1
		links = append(links, link)
	}
	return links, nil
}

// RemoveAgent removes an agent from a workspace.
func (s *WorkspaceAgentsStore) RemoveAgent(ctx context.Context, workspaceID, agentID int64) error {
	query := `DELETE FROM workspace_agents WHERE workspace_id = ? AND agent_id = ?`
	result, err := s.db.ExecContext(ctx, query, workspaceID, agentID)
	if err != nil {
		return fmt.Errorf("failed to remove workspace agent: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("workspace agent not found")
	}
	return nil
}

// DeleteAllAgents deletes all agents for a workspace.
func (s *WorkspaceAgentsStore) DeleteAllAgents(ctx context.Context, workspaceID int64) error {
	query := `DELETE FROM workspace_agents WHERE workspace_id = ?`
	_, err := s.db.ExecContext(ctx, query, workspaceID)
	if err != nil {
		return fmt.Errorf("failed to delete all workspace agents: %w", err)
	}
	return nil
}