package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// SkillStore handles persistence operations for skills.
type SkillStore struct {
	db *sql.DB
}

// NewSkillStore creates a new SkillStore instance.
func NewSkillStore(db *sql.DB) *SkillStore {
	return &SkillStore{db: db}
}

// CreateSkill creates a new skill.
func (s *SkillStore) CreateSkill(ctx context.Context, skill *Skill) error {
	query := `INSERT INTO skills (name, description, tags, content)
			  VALUES (?, ?, ?, ?)`

	_, err := s.db.ExecContext(ctx, query, skill.Name, skill.Description, skill.Tags, skill.Content)
	if err != nil {
		return fmt.Errorf("failed to create skill: %w", err)
	}
	return nil
}

// GetSkill retrieves a skill by ID.
func (s *SkillStore) GetSkill(ctx context.Context, id int64) (*Skill, error) {
	query := `SELECT id, name, description, tags, content FROM skills WHERE id = ?`

	var skill Skill
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&skill.ID, &skill.Name, &skill.Description, &skill.Tags, &skill.Content,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrSkillNotFound
		}
		return nil, fmt.Errorf("failed to get skill: %w", err)
	}

	return &skill, nil
}

// UpdateSkill updates an existing skill.
func (s *SkillStore) UpdateSkill(ctx context.Context, skill *Skill) error {
	query := `UPDATE skills SET name = ?, description = ?, tags = ?, content = ? WHERE id = ?`

	_, err := s.db.ExecContext(ctx, query,
		skill.Name, skill.Description, skill.Tags, skill.Content, skill.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update skill: %w", err)
	}
	return nil
}

// DeleteSkill deletes a skill.
func (s *SkillStore) DeleteSkill(ctx context.Context, id int64) error {
	query := `DELETE FROM skills WHERE id = ?`
	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete skill: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return ErrSkillNotFound
	}

	return nil
}

// ListSkills retrieves all skills.
func (s *SkillStore) ListSkills(ctx context.Context) ([]Skill, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, name, description, tags, content FROM skills ORDER BY name ASC`,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query skills: %w", err)
	}
	defer rows.Close()

	var skills []Skill
	for rows.Next() {
		var skill Skill
		if err := rows.Scan(&skill.ID, &skill.Name, &skill.Description, &skill.Tags, &skill.Content); err != nil {
			return nil, fmt.Errorf("failed to scan skill: %w", err)
		}
		skills = append(skills, skill)
	}

	return skills, nil
}

// SearchSkills searches skills by keyword in name or content.
func (s *SkillStore) SearchSkills(ctx context.Context, keyword string) ([]Skill, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, name, description, tags, content 
		 FROM skills 
		 WHERE name LIKE ? OR content LIKE ? 
		 ORDER BY name ASC`,
		"%"+keyword+"%", "%"+keyword+"%",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to search skills: %w", err)
	}
	defer rows.Close()

	var skills []Skill
	for rows.Next() {
		var skill Skill
		if err := rows.Scan(&skill.ID, &skill.Name, &skill.Description, &skill.Tags, &skill.Content); err != nil {
			return nil, fmt.Errorf("failed to scan skill: %w", err)
		}
		skills = append(skills, skill)
	}

	return skills, nil
}

// AgentStore handles persistence operations for agents.
type AgentStore struct {
	db *sql.DB
}

// NewAgentStore creates a new AgentStore instance.
func NewAgentStore(db *sql.DB) *AgentStore {
	return &AgentStore{db: db}
}

// CreateAgent creates a new agent.
func (s *AgentStore) CreateAgent(ctx context.Context, agent *Agent) error {
	query := `INSERT INTO agents 
		(name, description, type, provider_id, model_id, max_iteration, temperature,
		 system_prompt, color, icon)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := s.db.ExecContext(ctx, query,
		agent.Name, agent.Description, string(agent.Type), agent.ProviderID,
		agent.ModelID, agent.MaxIteration, agent.Temperature, agent.SystemPrompt,
		agent.Color, agent.Icon,
	)

	if err != nil {
		return fmt.Errorf("failed to create agent: %w", err)
	}
	return nil
}

// GetAgent retrieves an agent by ID.
func (s *AgentStore) GetAgent(ctx context.Context, id int64) (*Agent, error) {
	query := `SELECT 
		id, name, description, type, provider_id, model_id, max_iteration, temperature,
		system_prompt, color, icon
		FROM agents WHERE id = ?`

	var agent Agent
	var typeStr string
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&agent.ID, &agent.Name, &agent.Description, &typeStr, &agent.ProviderID,
		&agent.ModelID, &agent.MaxIteration, &agent.Temperature, &agent.SystemPrompt,
		&agent.Color, &agent.Icon,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrAgentNotFound
		}
		return nil, fmt.Errorf("failed to get agent: %w", err)
	}

	agent.Type = AgentType(typeStr)
	return &agent, nil
}

// UpdateAgent updates an existing agent.
func (s *AgentStore) UpdateAgent(ctx context.Context, agent *Agent) error {
	query := `UPDATE agents SET
		name = ?, description = ?, type = ?, provider_id = ?, model_id = ?, max_iteration = ?,
		temperature = ?, system_prompt = ?, color = ?, icon = ?
		WHERE id = ?`

	_, err := s.db.ExecContext(ctx, query,
		agent.Name, agent.Description, string(agent.Type), agent.ProviderID,
		agent.ModelID, agent.MaxIteration, agent.Temperature, agent.SystemPrompt,
		agent.Color, agent.Icon, agent.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update agent: %w", err)
	}
	return nil
}

// DeleteAgent deletes an agent.
func (s *AgentStore) DeleteAgent(ctx context.Context, id int64) error {
	query := `DELETE FROM agents WHERE id = ?`
	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete agent: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return ErrAgentNotFound
	}

	return nil
}

// ListAgents retrieves all agents.
func (s *AgentStore) ListAgents(ctx context.Context) ([]Agent, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT 
			id, name, description, type, provider_id, model_id, max_iteration, temperature,
			system_prompt, color, icon
			FROM agents ORDER BY name ASC`,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query agents: %w", err)
	}
	defer rows.Close()

	var agents []Agent
	for rows.Next() {
		var agent Agent
		var typeStr string
		if err := rows.Scan(
			&agent.ID, &agent.Name, &agent.Description, &typeStr, &agent.ProviderID,
			&agent.ModelID, &agent.MaxIteration, &agent.Temperature, &agent.SystemPrompt,
			&agent.Color, &agent.Icon,
		); err != nil {
			return nil, fmt.Errorf("failed to scan agent: %w", err)
		}
		agent.Type = AgentType(typeStr)
		agents = append(agents, agent)
	}

	return agents, nil
}