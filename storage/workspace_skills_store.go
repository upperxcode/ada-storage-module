package storage

import (
	"context"
	"database/sql"
)

type WorkspaceSkill struct {
	ID         int64 `json:"id"`
	WorkspaceID int64 `json:"workspace_id"`
	SkillID     int64 `json:"skill_id"`
	Enabled     bool  `json:"enabled"`
}

type WorkspaceSkillsStore struct {
	db *sql.DB
}

func NewWorkspaceSkillsStore(db *sql.DB) *WorkspaceSkillsStore {
	return &WorkspaceSkillsStore{db: db}
}

func (s *WorkspaceSkillsStore) Create(ctx context.Context, link *WorkspaceSkill) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO workspace_skills (workspace_id, skill_id, enabled) VALUES (?, ?, ?)`,
		link.WorkspaceID, link.SkillID, link.Enabled,
	)
	return err
}

func (s *WorkspaceSkillsStore) DeleteByWorkspace(ctx context.Context, workspaceID int64) error {
	_, err := s.db.ExecContext(ctx,
		`DELETE FROM workspace_skills WHERE workspace_id = ?`,
		workspaceID,
	)
	return err
}

func (s *WorkspaceSkillsStore) ListByWorkspace(ctx context.Context, workspaceID int64) ([]WorkspaceSkill, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, workspace_id, skill_id, enabled FROM workspace_skills WHERE workspace_id = ?`,
		workspaceID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []WorkspaceSkill
	for rows.Next() {
		var l WorkspaceSkill
		if err := rows.Scan(&l.ID, &l.WorkspaceID, &l.SkillID, &l.Enabled); err != nil {
			return nil, err
		}
		links = append(links, l)
	}
	return links, nil
}