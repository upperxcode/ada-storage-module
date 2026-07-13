package storage

import (
	"context"
	"database/sql"
	"fmt"
)

// Memory represents a workspace memory.
type Memory struct {
	ID          int64          `json:"id"`
	WorkspacePath sql.NullString `json:"workspace_path"`
	Content     string         `json:"content"`
	Importance  int            `json:"importance"`
	Embedding   []byte         `json:"embedding"`
	CreatedAt   string         `json:"created_at"`
	UpdatedAt   string         `json:"updated_at"`
}

// MemoryStore handles persistence operations for workspace memories.
type MemoryStore struct {
	db *sql.DB
}

// NewMemoryStore creates a new MemoryStore instance.
func NewMemoryStore(db *sql.DB) *MemoryStore {
	return &MemoryStore{db: db}
}

// CreateMemory creates a new memory for a workspace.
func (s *MemoryStore) CreateMemory(ctx context.Context, memory *Memory) error {
	query := `INSERT INTO memories (workspace_path, content, importance, embedding)
			  VALUES (?, ?, ?, ?)`

	_, err := s.db.ExecContext(ctx, query,
		memory.WorkspacePath, memory.Content, memory.Importance, memory.Embedding,
	)

	if err != nil {
		return fmt.Errorf("failed to create memory: %w", err)
	}
	return nil
}

// GetMemories retrieves all memories for a workspace.
func (s *MemoryStore) GetMemories(ctx context.Context, workspacePath string) ([]Memory, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, workspace_path, content, importance, embedding, created_at, updated_at
		 FROM memories WHERE workspace_path = ? ORDER BY created_at DESC`,
		workspacePath,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query memories: %w", err)
	}
	defer rows.Close()

	var memories []Memory
	for rows.Next() {
		var memory Memory
		if err := rows.Scan(&memory.ID, &memory.WorkspacePath, &memory.Content,
			&memory.Importance, &memory.Embedding, &memory.CreatedAt, &memory.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan memory: %w", err)
		}
		memories = append(memories, memory)
	}

	return memories, nil
}

// DeleteMemories deletes all memories for a workspace.
func (s *MemoryStore) DeleteMemories(ctx context.Context, workspacePath string) error {
	query := `DELETE FROM memories WHERE workspace_path = ?`
	_, err := s.db.ExecContext(ctx, query, workspacePath)
	if err != nil {
		return fmt.Errorf("failed to delete memories: %w", err)
	}
	return nil
}

// WorkspaceAgent represents a workspace-agent association.
type WorkspaceAgent struct {
	ID          int64  `json:"id"`
	WorkspaceID int64  `json:"workspace_id"`
	AgentID     int64  `json:"agent_id"`
	Enabled     bool   `json:"enabled"`
}

// WorkspaceAgentStore handles persistence operations for workspace-agent associations.
type WorkspaceAgentStore struct {
	db *sql.DB
}

// NewWorkspaceAgentStore creates a new WorkspaceAgentStore instance.
func NewWorkspaceAgentStore(db *sql.DB) *WorkspaceAgentStore {
	return &WorkspaceAgentStore{db: db}
}

// AddWorkspaceAgent links an agent to a workspace.
func (s *WorkspaceAgentStore) AddWorkspaceAgent(ctx context.Context, link *WorkspaceAgent) error {
	query := `INSERT INTO workspace_agents (workspace_id, agent_id, enabled)
			  VALUES (?, ?, ?)`

	_, err := s.db.ExecContext(ctx, query, link.WorkspaceID, link.AgentID, link.Enabled)
	if err != nil {
		return fmt.Errorf("failed to add workspace agent: %w", err)
	}
	return nil
}

// GetWorkspaceAgents retrieves all agents for a workspace.
func (s *WorkspaceAgentStore) GetWorkspaceAgents(ctx context.Context, workspaceID int64) ([]WorkspaceAgent, error) {
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

// RemoveWorkspaceAgent removes an agent from a workspace.
func (s *WorkspaceAgentStore) RemoveWorkspaceAgent(ctx context.Context, workspaceID, agentID int64) error {
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