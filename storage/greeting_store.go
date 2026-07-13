package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// GreetingStore handles persistence operations for greetings.
type GreetingStore struct {
	db *sql.DB
}

// NewGreetingStore creates a new GreetingStore instance.
func NewGreetingStore(db *sql.DB) *GreetingStore {
	return &GreetingStore{db: db}
}

// CreateGreeting creates a new greeting.
func (s *GreetingStore) CreateGreeting(ctx context.Context, greeting *Greeting) error {
	query := `INSERT INTO greetings (keyword, language, response)
			  VALUES (?, ?, ?)`

	_, err := s.db.ExecContext(ctx, query, greeting.Keyword, greeting.Language, greeting.Response)
	if err != nil {
		return fmt.Errorf("failed to create greeting: %w", err)
	}
	return nil
}

// GetGreeting retrieves a greeting by keyword and language.
func (s *GreetingStore) GetGreeting(ctx context.Context, keyword, language string) (*Greeting, error) {
	query := `SELECT id, keyword, language, response FROM greetings WHERE keyword = ? AND language = ?`

	var greeting Greeting
	err := s.db.QueryRowContext(ctx, query, keyword, language).Scan(
		&greeting.ID, &greeting.Keyword, &greeting.Language, &greeting.Response,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrGreetingNotFound
		}
		return nil, fmt.Errorf("failed to get greeting: %w", err)
	}

	return &greeting, nil
}

// LoadAllGreetings retrieves all greetings as a map keyed by keyword.
func (s *GreetingStore) LoadAllGreetings(ctx context.Context) (map[string]string, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT keyword, response FROM greetings`,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query greetings: %w", err)
	}
	defer rows.Close()

	greetings := make(map[string]string)
	for rows.Next() {
		var keyword, response string
		if err := rows.Scan(&keyword, &response); err != nil {
			return nil, fmt.Errorf("failed to scan greeting: %w", err)
		}
		greetings[keyword] = response
	}

	return greetings, nil
}

// DeleteGreeting removes a greeting.
func (s *GreetingStore) DeleteGreeting(ctx context.Context, keyword, language string) error {
	query := `DELETE FROM greetings WHERE keyword = ? AND language = ?`
	result, err := s.db.ExecContext(ctx, query, keyword, language)
	if err != nil {
		return fmt.Errorf("failed to delete greeting: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return ErrGreetingNotFound
	}

	return nil
}

// ListGreetings retrieves all greetings.
func (s *GreetingStore) ListGreetings(ctx context.Context) ([]Greeting, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, keyword, language, response FROM greetings ORDER BY keyword, language`,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query greetings: %w", err)
	}
	defer rows.Close()

	var greetings []Greeting
	for rows.Next() {
		var g Greeting
		if err := rows.Scan(&g.ID, &g.Keyword, &g.Language, &g.Response); err != nil {
			return nil, fmt.Errorf("failed to scan greeting: %w", err)
		}
		greetings = append(greetings, g)
	}

	return greetings, nil
}