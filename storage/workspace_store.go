package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// ErrWorkspaceNotFound is returned when a workspace does not exist.
var ErrWorkspaceNotFound = errors.New("workspace not found")

// Workspace represents a workspace entity.
type Workspace struct {
	ID                 int64          `json:"id"`
	Nome               string         `json:"nome"`
	Description        sql.NullString `json:"description"`
	Path               sql.NullString `json:"path"`
	MaxPrompt          int            `json:"max_prompt"`
	MaxContent         int            `json:"max_content"`
	Commit             bool           `json:"commit"`
	SpecProvider       sql.NullString `json:"spec_provider"`
	SpecWizardID       sql.NullString `json:"spec_wizard_id"`
	Personality        sql.NullString `json:"personality"`
	Color              string         `json:"color"`
	Icon               string         `json:"icon"`
	Summary            sql.NullString `json:"summary"`
	Enabled            bool           `json:"enabled"`
	MaxPromptSend      int            `json:"max_prompt_send"`
	CommitChanges      bool           `json:"commit_changes"`
	MaxContextLength   int            `json:"max_context_length"`
	EmbeddingModel     sql.NullString `json:"embedding_model"`
	EmbeddingProvider  sql.NullString `json:"embedding_provider"`
	RoutingRules       sql.NullString `json:"routing_rules"`
}

// WorkspaceStore handles persistence operations for workspaces.
type WorkspaceStore struct {
	db *sql.DB
}

// NewWorkspaceStore creates a new WorkspaceStore instance.
func NewWorkspaceStore(db *sql.DB) *WorkspaceStore {
	return &WorkspaceStore{db: db}
}

// CreateWorkspace creates a new workspace.
func (s *WorkspaceStore) CreateWorkspace(ctx context.Context, workspace *Workspace) error {
	query := `INSERT INTO workspaces 
		(nome, description, path, max_prompt, max_content, "commit", spec_provider,
		 spec_wizard_id, personality, color, icon, summary, enabled, max_prompt_send,
		 commit_changes, max_context_length, embedding_model, embedding_provider, routing_rules)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := s.db.ExecContext(ctx, query,
		workspace.Nome, workspace.Description, workspace.Path, workspace.MaxPrompt,
		workspace.MaxContent, workspace.Commit, workspace.SpecProvider, workspace.SpecWizardID,
		workspace.Personality, workspace.Color, workspace.Icon, workspace.Summary,
		workspace.Enabled, workspace.MaxPromptSend, workspace.CommitChanges,
		workspace.MaxContextLength, workspace.EmbeddingModel, workspace.EmbeddingProvider,
		workspace.RoutingRules,
	)

	if err != nil {
		return fmt.Errorf("failed to create workspace: %w", err)
	}
	return nil
}

// GetWorkspace retrieves a workspace by ID.
func (s *WorkspaceStore) GetWorkspace(ctx context.Context, id int64) (*Workspace, error) {
	query := `SELECT 
		id, nome, description, path, max_prompt, max_content, "commit", spec_provider,
		spec_wizard_id, personality, color, icon, summary, enabled, max_prompt_send,
		commit_changes, max_context_length, embedding_model, embedding_provider, routing_rules
		FROM workspaces WHERE id = ?`

	var workspace Workspace
	var commitInt, enabledInt int
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&workspace.ID, &workspace.Nome, &workspace.Description, &workspace.Path,
		&workspace.MaxPrompt, &workspace.MaxContent, &commitInt, &workspace.SpecProvider,
		&workspace.SpecWizardID, &workspace.Personality, &workspace.Color, &workspace.Icon,
		&workspace.Summary, &enabledInt, &workspace.MaxPromptSend, &workspace.CommitChanges,
		&workspace.MaxContextLength, &workspace.EmbeddingModel, &workspace.EmbeddingProvider,
		&workspace.RoutingRules,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrWorkspaceNotFound
		}
		return nil, fmt.Errorf("failed to get workspace: %w", err)
	}

	workspace.Commit = commitInt == 1
	workspace.Enabled = enabledInt == 1
	return &workspace, nil
}

// GetWorkspaceByPath retrieves a workspace by its path.
func (s *WorkspaceStore) GetWorkspaceByPath(ctx context.Context, path string) (*Workspace, error) {
	query := `SELECT 
		id, nome, description, path, max_prompt, max_content, "commit", spec_provider,
		spec_wizard_id, personality, color, icon, summary, enabled, max_prompt_send,
		commit_changes, max_context_length, embedding_model, embedding_provider, routing_rules
		FROM workspaces WHERE path = ?`

	var workspace Workspace
	var commitInt, enabledInt int
	err := s.db.QueryRowContext(ctx, query, path).Scan(
		&workspace.ID, &workspace.Nome, &workspace.Description, &workspace.Path,
		&workspace.MaxPrompt, &workspace.MaxContent, &commitInt, &workspace.SpecProvider,
		&workspace.SpecWizardID, &workspace.Personality, &workspace.Color, &workspace.Icon,
		&workspace.Summary, &enabledInt, &workspace.MaxPromptSend, &workspace.CommitChanges,
		&workspace.MaxContextLength, &workspace.EmbeddingModel, &workspace.EmbeddingProvider,
		&workspace.RoutingRules,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrWorkspaceNotFound
		}
		return nil, fmt.Errorf("failed to get workspace by path: %w", err)
	}

	workspace.Commit = commitInt == 1
	workspace.Enabled = enabledInt == 1
	return &workspace, nil
}

// UpdateWorkspace updates an existing workspace.
func (s *WorkspaceStore) UpdateWorkspace(ctx context.Context, workspace *Workspace) error {
	query := `UPDATE workspaces SET
		nome = ?, description = ?, path = ?, max_prompt = ?, max_content = ?, "commit" = ?,
		spec_provider = ?, spec_wizard_id = ?, personality = ?, color = ?, icon = ?, summary = ?,
		enabled = ?, max_prompt_send = ?, commit_changes = ?, max_context_length = ?,
		embedding_model = ?, embedding_provider = ?, routing_rules = ?
		WHERE id = ?`

	_, err := s.db.ExecContext(ctx, query,
		workspace.Nome, workspace.Description, workspace.Path, workspace.MaxPrompt,
		workspace.MaxContent, workspace.Commit, workspace.SpecProvider, workspace.SpecWizardID,
		workspace.Personality, workspace.Color, workspace.Icon, workspace.Summary,
		workspace.Enabled, workspace.MaxPromptSend, workspace.CommitChanges,
		workspace.MaxContextLength, workspace.EmbeddingModel, workspace.EmbeddingProvider,
		workspace.RoutingRules, workspace.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update workspace: %w", err)
	}
	return nil
}

// DeleteWorkspace deletes a workspace and all its associated data (due to CASCADE).
func (s *WorkspaceStore) DeleteWorkspace(ctx context.Context, id int64) error {
	query := `DELETE FROM workspaces WHERE id = ?`
	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete workspace: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return ErrWorkspaceNotFound
	}

	return nil
}

// ListWorkspaces retrieves all workspaces.
func (s *WorkspaceStore) ListWorkspaces(ctx context.Context) ([]Workspace, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT 
			id, nome, description, path, max_prompt, max_content, "commit", spec_provider,
			spec_wizard_id, personality, color, icon, summary, enabled, max_prompt_send,
			commit_changes, max_context_length, embedding_model, embedding_provider, routing_rules
			FROM workspaces ORDER BY nome ASC`,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query workspaces: %w", err)
	}
	defer rows.Close()

	var workspaces []Workspace
	for rows.Next() {
		var workspace Workspace
		var commitInt, enabledInt int
		if err := rows.Scan(
			&workspace.ID, &workspace.Nome, &workspace.Description, &workspace.Path,
			&workspace.MaxPrompt, &workspace.MaxContent, &commitInt, &workspace.SpecProvider,
			&workspace.SpecWizardID, &workspace.Personality, &workspace.Color, &workspace.Icon,
			&workspace.Summary, &enabledInt, &workspace.MaxPromptSend, &workspace.CommitChanges,
			&workspace.MaxContextLength, &workspace.EmbeddingModel, &workspace.EmbeddingProvider,
			&workspace.RoutingRules,
		); err != nil {
			return nil, fmt.Errorf("failed to scan workspace: %w", err)
		}
		workspace.Commit = commitInt == 1
		workspace.Enabled = enabledInt == 1
		workspaces = append(workspaces, workspace)
	}

	return workspaces, nil
}