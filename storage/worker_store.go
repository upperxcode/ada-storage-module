package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// WorkerStore handles persistence operations for workers.
type WorkerStore struct {
	db *sql.DB
}

// NewWorkerStore creates a new WorkerStore instance.
func NewWorkerStore(db *sql.DB) *WorkerStore {
	return &WorkerStore{db: db}
}

// CreateWorker creates a new worker.
func (s *WorkerStore) CreateWorker(ctx context.Context, worker *Worker) error {
	query := `INSERT INTO workers 
		(name, persona, response_language, connection_type, command, arguments, environment,
		 inheritance_folders, inheritance_skills, inheritance_persona, inheritance_knowledge,
		 inheritance_tools, color, icon)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := s.db.ExecContext(ctx, query,
		worker.Name, worker.Persona, worker.ResponseLanguage, worker.ConnectionType,
		worker.Command, worker.Arguments, worker.Environment, worker.InheritanceFolders,
		worker.InheritanceSkills, worker.InheritancePersona, worker.InheritanceKnowledge,
		worker.InheritanceTools, worker.Color, worker.Icon,
	)

	if err != nil {
		return fmt.Errorf("failed to create worker: %w", err)
	}
	return nil
}

// GetWorker retrieves a worker by ID.
func (s *WorkerStore) GetWorker(ctx context.Context, id int64) (*Worker, error) {
	query := `SELECT 
		id, name, persona, response_language, connection_type, command, arguments, environment,
		inheritance_folders, inheritance_skills, inheritance_persona, inheritance_knowledge,
		inheritance_tools, color, icon
		FROM workers WHERE id = ?`

	var worker Worker
	var connFolders, connSkills, connPersona, connKnowledge, connTools int
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&worker.ID, &worker.Name, &worker.Persona, &worker.ResponseLanguage,
		&worker.ConnectionType, &worker.Command, &worker.Arguments, &worker.Environment,
		&connFolders, &connSkills, &connPersona, &connKnowledge, &connTools,
		&worker.Color, &worker.Icon,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrWorkerNotFound
		}
		return nil, fmt.Errorf("failed to get worker: %w", err)
	}

	worker.InheritanceFolders = connFolders == 1
	worker.InheritanceSkills = connSkills == 1
	worker.InheritancePersona = connPersona == 1
	worker.InheritanceKnowledge = connKnowledge == 1
	worker.InheritanceTools = connTools == 1

	return &worker, nil
}

// GetWorkerByName retrieves a worker by name.
func (s *WorkerStore) GetWorkerByName(ctx context.Context, name string) (*Worker, error) {
	query := `SELECT 
		id, name, persona, response_language, connection_type, command, arguments, environment,
		inheritance_folders, inheritance_skills, inheritance_persona, inheritance_knowledge,
		inheritance_tools, color, icon
		FROM workers WHERE name = ?`

	var worker Worker
	var connFolders, connSkills, connPersona, connKnowledge, connTools int
	err := s.db.QueryRowContext(ctx, query, name).Scan(
		&worker.ID, &worker.Name, &worker.Persona, &worker.ResponseLanguage,
		&worker.ConnectionType, &worker.Command, &worker.Arguments, &worker.Environment,
		&connFolders, &connSkills, &connPersona, &connKnowledge, &connTools,
		&worker.Color, &worker.Icon,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrWorkerNotFound
		}
		return nil, fmt.Errorf("failed to get worker by name: %w", err)
	}

	worker.InheritanceFolders = connFolders == 1
	worker.InheritanceSkills = connSkills == 1
	worker.InheritancePersona = connPersona == 1
	worker.InheritanceKnowledge = connKnowledge == 1
	worker.InheritanceTools = connTools == 1

	return &worker, nil
}

// UpdateWorker updates an existing worker.
func (s *WorkerStore) UpdateWorker(ctx context.Context, worker *Worker) error {
	query := `UPDATE workers SET
		name = ?, persona = ?, response_language = ?, connection_type = ?, command = ?, arguments = ?,
		environment = ?, inheritance_folders = ?, inheritance_skills = ?, inheritance_persona = ?,
		inheritance_knowledge = ?, inheritance_tools = ?, color = ?, icon = ?
		WHERE id = ?`

	_, err := s.db.ExecContext(ctx, query,
		worker.Name, worker.Persona, worker.ResponseLanguage, worker.ConnectionType,
		worker.Command, worker.Arguments, worker.Environment, worker.InheritanceFolders,
		worker.InheritanceSkills, worker.InheritancePersona, worker.InheritanceKnowledge,
		worker.InheritanceTools, worker.Color, worker.Icon, worker.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update worker: %w", err)
	}
	return nil
}

// DeleteWorker deletes a worker.
func (s *WorkerStore) DeleteWorker(ctx context.Context, id int64) error {
	query := `DELETE FROM workers WHERE id = ?`
	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete worker: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return ErrWorkerNotFound
	}

	return nil
}

// ListWorkers retrieves all workers.
func (s *WorkerStore) ListWorkers(ctx context.Context) ([]Worker, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT 
			id, name, persona, response_language, connection_type, command, arguments, environment,
			inheritance_folders, inheritance_skills, inheritance_persona, inheritance_knowledge,
			inheritance_tools, color, icon
			FROM workers ORDER BY name ASC`,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query workers: %w", err)
	}
	defer rows.Close()

	var workers []Worker
	for rows.Next() {
		var worker Worker
		var connFolders, connSkills, connPersona, connKnowledge, connTools int
		if err := rows.Scan(
			&worker.ID, &worker.Name, &worker.Persona, &worker.ResponseLanguage,
			&worker.ConnectionType, &worker.Command, &worker.Arguments, &worker.Environment,
			&connFolders, &connSkills, &connPersona, &connKnowledge, &connTools,
			&worker.Color, &worker.Icon,
		); err != nil {
			return nil, fmt.Errorf("failed to scan worker: %w", err)
		}
		worker.InheritanceFolders = connFolders == 1
		worker.InheritanceSkills = connSkills == 1
		worker.InheritancePersona = connPersona == 1
		worker.InheritanceKnowledge = connKnowledge == 1
		worker.InheritanceTools = connTools == 1
		workers = append(workers, worker)
	}

	return workers, nil
}