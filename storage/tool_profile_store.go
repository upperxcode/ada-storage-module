package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// ErrToolProfileNotFound is returned when a tool profile does not exist.
var ErrToolProfileNotFound = errors.New("tool profile not found")

// ToolProfile represents a tool profile configuration.
type ToolProfile struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description sql.NullString `json:"description"`
	Color       string `json:"color"`
	Icon        string `json:"icon"`
}

// ToolProfileTool represents a tool in a profile.
type ToolProfileTool struct {
	ID        int64  `json:"id"`
	ProfileID int64  `json:"profile_id"`
	ToolName  string `json:"tool_name"`
}

// ToolProfileStore handles persistence operations for tool profiles.
type ToolProfileStore struct {
	db *sql.DB
}

// NewToolProfileStore creates a new ToolProfileStore instance.
func NewToolProfileStore(db *sql.DB) *ToolProfileStore {
	return &ToolProfileStore{db: db}
}

// CreateProfile creates a new tool profile and sets its ID to the auto-incremented value.
func (s *ToolProfileStore) CreateProfile(ctx context.Context, profile *ToolProfile) error {
	query := `INSERT INTO tool_profiles (name, description, color, icon) VALUES (?, ?, ?, ?)`
	result, err := s.db.ExecContext(ctx, query, profile.Name, profile.Description, profile.Color, profile.Icon)
	if err != nil {
		return fmt.Errorf("failed to create tool profile: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}
	profile.ID = id
	return nil
}

// GetProfile retrieves a tool profile by ID.
func (s *ToolProfileStore) GetProfile(ctx context.Context, id int64) (*ToolProfile, error) {
	query := `SELECT id, name, description, color, icon FROM tool_profiles WHERE id = ?`
	var profile ToolProfile
	err := s.db.QueryRowContext(ctx, query, id).Scan(&profile.ID, &profile.Name, &profile.Description, &profile.Color, &profile.Icon)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrToolProfileNotFound
		}
		return nil, fmt.Errorf("failed to get tool profile: %w", err)
	}
	return &profile, nil
}

// ListProfiles retrieves all tool profiles.
func (s *ToolProfileStore) ListProfiles(ctx context.Context) ([]ToolProfile, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, name, description, color, icon FROM tool_profiles ORDER BY name`)
	if err != nil {
		return nil, fmt.Errorf("failed to query tool profiles: %w", err)
	}
	defer rows.Close()

	var profiles []ToolProfile
	for rows.Next() {
		var p ToolProfile
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Color, &p.Icon); err != nil {
			return nil, fmt.Errorf("failed to scan tool profile: %w", err)
		}
		profiles = append(profiles, p)
	}
	return profiles, nil
}

// DeleteProfile deletes a tool profile and its tools (CASCADE).
func (s *ToolProfileStore) DeleteProfile(ctx context.Context, id int64) error {
	query := `DELETE FROM tool_profiles WHERE id = ?`
	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete tool profile: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return ErrToolProfileNotFound
	}
	return nil
}

// AddTool adds a tool to a profile.
func (s *ToolProfileStore) AddTool(ctx context.Context, tool *ToolProfileTool) error {
	query := `INSERT INTO tool_profile_tools (profile_id, tool_name) VALUES (?, ?)`
	_, err := s.db.ExecContext(ctx, query, tool.ProfileID, tool.ToolName)
	if err != nil {
		return fmt.Errorf("failed to add tool to profile: %w", err)
	}
	return nil
}

// RemoveTool removes a tool from a profile.
func (s *ToolProfileStore) RemoveTool(ctx context.Context, profileID int64, toolName string) error {
	query := `DELETE FROM tool_profile_tools WHERE profile_id = ? AND tool_name = ?`
	_, err := s.db.ExecContext(ctx, query, profileID, toolName)
	if err != nil {
		return fmt.Errorf("failed to remove tool from profile: %w", err)
	}
	return nil
}

// ListTools retrieves all tools for a profile.
func (s *ToolProfileStore) ListTools(ctx context.Context, profileID int64) ([]ToolProfileTool, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, profile_id, tool_name FROM tool_profile_tools WHERE profile_id = ?`,
		profileID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query profile tools: %w", err)
	}
	defer rows.Close()

	var tools []ToolProfileTool
	for rows.Next() {
		var t ToolProfileTool
		if err := rows.Scan(&t.ID, &t.ProfileID, &t.ToolName); err != nil {
			return nil, fmt.Errorf("failed to scan profile tool: %w", err)
		}
		tools = append(tools, t)
	}
	return tools, nil
}

// UpdateProfile updates an existing tool profile.
func (s *ToolProfileStore) UpdateProfile(ctx context.Context, profile *ToolProfile) error {
	query := `UPDATE tool_profiles SET name = ?, description = ?, color = ?, icon = ? WHERE id = ?`
	result, err := s.db.ExecContext(ctx, query, profile.Name, profile.Description, profile.Color, profile.Icon, profile.ID)
	if err != nil {
		return fmt.Errorf("failed to update tool profile: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return ErrToolProfileNotFound
	}
	return nil
}

// GetProfileByName retrieves a tool profile by name.
func (s *ToolProfileStore) GetProfileByName(ctx context.Context, name string) (*ToolProfile, error) {
	query := `SELECT id, name, description, color, icon FROM tool_profiles WHERE name = ?`
	var profile ToolProfile
	err := s.db.QueryRowContext(ctx, query, name).Scan(&profile.ID, &profile.Name, &profile.Description, &profile.Color, &profile.Icon)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrToolProfileNotFound
		}
		return nil, fmt.Errorf("failed to get tool profile by name: %w", err)
	}
	return &profile, nil
}