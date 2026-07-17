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