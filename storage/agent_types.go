package storage

import (
	"database/sql"
	"errors"
)

// Agent-related errors
var (
	ErrAgentNotFound = errors.New("agent not found")
	ErrWorkerNotFound = errors.New("worker not found")
	ErrSkillNotFound = errors.New("skill not found")
)

// AgentType represents the type of an agent.
type AgentType string

const (
	AgentTypeExecutor  AgentType = "executor"
	AgentTypeDelegator AgentType = "delegator"
	AgentTypeReviewer  AgentType = "reviewer"
	AgentTypeResearch  AgentType = "research"
)

// ConnectionType represents the type of worker connection.
type ConnectionType string

const (
	ConnectionTypeWebsocket ConnectionType = "websocket"
	ConnectionTypeURL       ConnectionType = "url"
	ConnectionTypeCLI       ConnectionType = "cli_command"
)

// Agent represents an agent entity.
type Agent struct {
	ID           int64          `json:"id"`
	Name         string         `json:"name"`
	Description  sql.NullString `json:"description"`
	Type         AgentType      `json:"type"`
	ProviderID   sql.NullInt64  `json:"provider_id"`
	ModelID      sql.NullInt64  `json:"model_id"`
	MaxIteration int            `json:"max_iteration"`
	Temperature  float64        `json:"temperature"`
	SystemPrompt sql.NullString `json:"system_prompt"`
	Color        string         `json:"color"`
	Icon         string         `json:"icon"`
}

// Worker represents a worker entity.
type Worker struct {
	ID                   int64          `json:"id"`
	Name                 string         `json:"name"`
	Persona              sql.NullString `json:"persona"`
	ResponseLanguage     string         `json:"response_language"`
	ConnectionType       string         `json:"connection_type"`
	Command              sql.NullString `json:"command"`
	Arguments            sql.NullString `json:"arguments"`
	Environment          sql.NullString `json:"environment"`
	InheritanceFolders   bool           `json:"inheritance_folders"`
	InheritanceSkills    bool           `json:"inheritance_skills"`
	InheritancePersona   bool           `json:"inheritance_persona"`
	InheritanceKnowledge bool           `json:"inheritance_knowledge"`
	InheritanceTools     bool           `json:"inheritance_tools"`
	Color                string         `json:"color"`
	Icon                 string         `json:"icon"`
}

// Skill represents a skill entity.
type Skill struct {
	ID          int64          `json:"id"`
	Name        string         `json:"name"`
	Description sql.NullString `json:"description"`
	Tags        sql.NullString `json:"tags"`
	Content     string         `json:"content"`
	Color       string         `json:"color"`
	Icon        string         `json:"icon"`
}