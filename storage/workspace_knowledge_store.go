package storage

import (
	"context"
	"database/sql"
	"fmt"
)

// WorkspaceKnowledge represents a knowledge item for a workspace.
type WorkspaceKnowledge struct {
	ID          int64  `json:"id"`
	WorkspaceID int64  `json:"workspace_id"`
	KnowledgeItem string `json:"knowledge_item"`
}

// WorkspaceKnowledgeStore handles persistence operations for workspace knowledge items.
type WorkspaceKnowledgeStore struct {
	db *sql.DB
}

// NewWorkspaceKnowledgeStore creates a new WorkspaceKnowledgeStore instance.
func NewWorkspaceKnowledgeStore(db *sql.DB) *WorkspaceKnowledgeStore {
	return &WorkspaceKnowledgeStore{db: db}
}

// AddKnowledge adds a knowledge item to a workspace.
func (s *WorkspaceKnowledgeStore) AddKnowledge(ctx context.Context, knowledge *WorkspaceKnowledge) error {
	query := `INSERT INTO workspace_knowledge (workspace_id, knowledge_item) VALUES (?, ?)`
	_, err := s.db.ExecContext(ctx, query, knowledge.WorkspaceID, knowledge.KnowledgeItem)
	if err != nil {
		return fmt.Errorf("failed to add knowledge: %w", err)
	}
	return nil
}

// ListKnowledge retrieves all knowledge items for a workspace.
func (s *WorkspaceKnowledgeStore) ListKnowledge(ctx context.Context, workspaceID int64) ([]WorkspaceKnowledge, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, workspace_id, knowledge_item FROM workspace_knowledge WHERE workspace_id = ?`,
		workspaceID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query knowledge: %w", err)
	}
	defer rows.Close()

	var knowledge []WorkspaceKnowledge
	for rows.Next() {
		var k WorkspaceKnowledge
		if err := rows.Scan(&k.ID, &k.WorkspaceID, &k.KnowledgeItem); err != nil {
			return nil, fmt.Errorf("failed to scan knowledge: %w", err)
		}
		knowledge = append(knowledge, k)
	}
	return knowledge, nil
}

// DeleteKnowledge deletes a knowledge item.
func (s *WorkspaceKnowledgeStore) DeleteKnowledge(ctx context.Context, workspaceID int64, knowledgeItem string) error {
	query := `DELETE FROM workspace_knowledge WHERE workspace_id = ? AND knowledge_item = ?`
	_, err := s.db.ExecContext(ctx, query, workspaceID, knowledgeItem)
	if err != nil {
		return fmt.Errorf("failed to delete knowledge: %w", err)
	}
	return nil
}

// DeleteAllKnowledge deletes all knowledge items for a workspace.
func (s *WorkspaceKnowledgeStore) DeleteAllKnowledge(ctx context.Context, workspaceID int64) error {
	query := `DELETE FROM workspace_knowledge WHERE workspace_id = ?`
	_, err := s.db.ExecContext(ctx, query, workspaceID)
	if err != nil {
		return fmt.Errorf("failed to delete all knowledge: %w", err)
	}
	return nil
}