package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// ErrThinkingNotFound is returned when a thinking record does not exist.
var ErrThinkingNotFound = errors.New("thinking not found")

// Thinking represents a thinking record linked to a message.
type Thinking struct {
	ID        int64     `json:"id"`
	MessageID string    `json:"message_id"` // ID of the associated message (as string/UUID)
	SessionID string    `json:"session_id"`
	Content   string    `json:"content"`
	Duration  int       `json:"duration"` // thinking duration in seconds
	CreatedAt time.Time `json:"created_at"`
}

// ThinkingStore handles persistence operations for thinking records.
type ThinkingStore struct {
	db *sql.DB
}

// NewThinkingStore creates a new ThinkingStore instance.
func NewThinkingStore(db *sql.DB) *ThinkingStore {
	return &ThinkingStore{db: db}
}

// SaveThinking saves a thinking record for a message.
func (s *ThinkingStore) SaveThinking(ctx context.Context, sessionID, messageID string, content string, duration int) error {
	query := `INSERT INTO thinkings (session_id, message_id, content, duration)
			  VALUES (?, ?, ?, ?)`

	_, err := s.db.ExecContext(ctx, query, sessionID, messageID, content, duration)
	if err != nil {
		return fmt.Errorf("failed to save thinking: %w", err)
	}
	return nil
}

// GetThinkingByMessage retrieves the thinking record associated with a message.
// Returns ErrThinkingNotFound if no record exists.
func (s *ThinkingStore) GetThinkingByMessage(ctx context.Context, messageID string) (Thinking, error) {
	query := `SELECT id, message_id, session_id, content, duration, created_at
			  FROM thinkings WHERE message_id = ?`

	var t Thinking
	err := s.db.QueryRowContext(ctx, query, messageID).Scan(
		&t.ID, &t.MessageID, &t.SessionID, &t.Content, &t.Duration, &t.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Thinking{}, ErrThinkingNotFound
		}
		return Thinking{}, fmt.Errorf("failed to get thinking: %w", err)
	}
	return t, nil
}

// DeleteThinking deletes the thinking record associated with a message.
// Returns ErrThinkingNotFound if no record exists.
func (s *ThinkingStore) DeleteThinking(ctx context.Context, messageID string) error {
	query := `DELETE FROM thinkings WHERE message_id = ?`

	result, err := s.db.ExecContext(ctx, query, messageID)
	if err != nil {
		return fmt.Errorf("failed to delete thinking: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return ErrThinkingNotFound
	}
	return nil
}

// ListThinkingsBySession retrieves all thinking records for a session,
// ordered by creation time ascending.
func (s *ThinkingStore) ListThinkingsBySession(ctx context.Context, sessionID string) ([]Thinking, error) {
	query := `SELECT id, message_id, session_id, content, duration, created_at
			  FROM thinkings WHERE session_id = ? ORDER BY created_at ASC`

	rows, err := s.db.QueryContext(ctx, query, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to list thinkings: %w", err)
	}
	defer rows.Close()

	var thinkings []Thinking
	for rows.Next() {
		var t Thinking
		if err := rows.Scan(&t.ID, &t.MessageID, &t.SessionID, &t.Content, &t.Duration, &t.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan thinking: %w", err)
		}
		thinkings = append(thinkings, t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating thinkings: %w", err)
	}

	return thinkings, nil
}
