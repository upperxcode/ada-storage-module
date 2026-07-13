package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ErrSessionNotFound is returned when a session does not exist.
var ErrSessionNotFound = errors.New("session not found")

// ErrMessageNotFound is returned when a message does not exist.
var ErrMessageNotFound = errors.New("message not found")

// Session represents a chat session.
type Session struct {
	ID                   string         `json:"id"`
	WorkspacePath        sql.NullString `json:"workspace_path"`
	Title                sql.NullString `json:"title"`
	Pinned               bool           `json:"pinned"`
	Embedding            []byte         `json:"embedding"`
	CreatedAt            time.Time      `json:"created_at"`
	UpdatedAt            time.Time      `json:"updated_at"`
	WorkerName           string         `json:"worker_name"`
	ParentSessionID      string         `json:"parent_session_id"`
	Model                string         `json:"model"`
	Provider             string         `json:"provider"`
	Mode                 string         `json:"mode"`
	Thinking             sql.NullString `json:"thinking"`
	Summary              sql.NullString `json:"summary"`
	SummarizedContext    sql.NullString `json:"summarized_context"`
	SummaryTokenCount    int            `json:"summary_token_count"`
	SummarizedAt         sql.NullTime   `json:"summarized_at"`
	LastSummarizedMsgID  int64          `json:"last_summarized_msg_id"`
}

// Message represents a chat message.
type Message struct {
	ID         int64          `json:"id"`
	SessionID  string         `json:"session_id"`
	Role       string         `json:"role"`
	Content    string         `json:"content"`
	Tokens     int            `json:"tokens"`
	Time       time.Time      `json:"time"`
	ServedBy   sql.NullString `json:"served_by"`
}

// SessionStore handles persistence operations for sessions and messages.
type SessionStore struct {
	db *sql.DB
}

// NewSessionStore creates a new SessionStore instance.
func NewSessionStore(db *sql.DB) *SessionStore {
	return &SessionStore{db: db}
}

// GenerateSessionID generates a new unique session ID.
func GenerateSessionID() string {
	return uuid.NewString()
}

// GetSession retrieves a session by ID.
// Returns ErrSessionNotFound if the session does not exist.
func (s *SessionStore) GetSession(ctx context.Context, id string) (*Session, error) {
	query := `SELECT 
		id, workspace_path, title, pinned, embedding, created_at, updated_at,
		worker_name, parent_session_id, model, provider, mode, thinking, summary,
		summarized_context, summary_token_count, summarized_at, last_summarized_msg_id
		FROM sessions WHERE id = ?`

	var session Session
	var pinnedInt int
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&session.ID, &session.WorkspacePath, &session.Title, &pinnedInt, &session.Embedding,
		&session.CreatedAt, &session.UpdatedAt, &session.WorkerName, &session.ParentSessionID,
		&session.Model, &session.Provider, &session.Mode, &session.Thinking, &session.Summary,
		&session.SummarizedContext, &session.SummaryTokenCount, &session.SummarizedAt,
		&session.LastSummarizedMsgID,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrSessionNotFound
		}
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	session.Pinned = pinnedInt == 1
	return &session, nil
}

// CreateSession creates a new session with the given ID.
func (s *SessionStore) CreateSession(ctx context.Context, session *Session) error {
	query := `INSERT INTO sessions 
		(id, workspace_path, title, pinned, embedding, worker_name, parent_session_id,
		 model, provider, mode, thinking, summary, summarized_context, summary_token_count,
		 summarized_at, last_summarized_msg_id)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := s.db.ExecContext(ctx, query,
		session.ID, session.WorkspacePath, session.Title, session.Pinned, session.Embedding,
		session.WorkerName, session.ParentSessionID, session.Model, session.Provider, session.Mode,
		session.Thinking, session.Summary, session.SummarizedContext, session.SummaryTokenCount,
		session.SummarizedAt, session.LastSummarizedMsgID,
	)

	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	return nil
}

// UpdateSession updates an existing session.
func (s *SessionStore) UpdateSession(ctx context.Context, session *Session) error {
	query := `UPDATE sessions SET
		workspace_path = ?, title = ?, pinned = ?, embedding = ?, worker_name = ?,
		parent_session_id = ?, model = ?, provider = ?, mode = ?, thinking = ?, summary = ?,
		summarized_context = ?, summary_token_count = ?, summarized_at = ?, last_summarized_msg_id = ?,
		updated_at = CURRENT_TIMESTAMP
		WHERE id = ?`

	_, err := s.db.ExecContext(ctx, query,
		session.WorkspacePath, session.Title, session.Pinned, session.Embedding,
		session.WorkerName, session.ParentSessionID, session.Model, session.Provider,
		session.Mode, session.Thinking, session.Summary, session.SummarizedContext,
		session.SummaryTokenCount, session.SummarizedAt, session.LastSummarizedMsgID,
		session.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update session: %w", err)
	}
	return nil
}

// DeleteSession deletes a session and all its messages (due to CASCADE).
func (s *SessionStore) DeleteSession(ctx context.Context, id string) error {
	query := `DELETE FROM sessions WHERE id = ?`
	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return ErrSessionNotFound
	}

	return nil
}

// SaveMessage saves a message to the database.
func (s *SessionStore) SaveMessage(ctx context.Context, msg *Message) error {
	query := `INSERT INTO messages 
		(session_id, role, content, tokens, time, served_by)
		VALUES (?, ?, ?, ?, ?, ?)`

	_, err := s.db.ExecContext(ctx, query,
		msg.SessionID, msg.Role, msg.Content, msg.Tokens, msg.Time, msg.ServedBy,
	)

	if err != nil {
		return fmt.Errorf("failed to save message: %w", err)
	}
	return nil
}

// GetMessages retrieves all messages for a session in chronological order.
// Returns ErrSessionNotFound if the session does not exist.
func (s *SessionStore) GetMessages(ctx context.Context, sessionID string) ([]Message, error) {
	// Verify session exists
	var exists bool
	err := s.db.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM sessions WHERE id = ?)`,
		sessionID,
	).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("failed to check session existence: %w", err)
	}
	if !exists {
		return nil, ErrSessionNotFound
	}

	// Retrieve messages
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, session_id, role, content, tokens, time, served_by
		 FROM messages WHERE session_id = ? ORDER BY time ASC`,
		sessionID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query messages: %w", err)
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.ID, &msg.SessionID, &msg.Role, &msg.Content,
			&msg.Tokens, &msg.Time, &msg.ServedBy); err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}
		messages = append(messages, msg)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating messages: %w", err)
	}

	return messages, nil
}

// DeleteMessages deletes all messages for a session.
func (s *SessionStore) DeleteMessages(ctx context.Context, sessionID string) error {
	query := `DELETE FROM messages WHERE session_id = ?`
	_, err := s.db.ExecContext(ctx, query, sessionID)
	if err != nil {
		return fmt.Errorf("failed to delete messages: %w", err)
	}
	return nil
}

// ListSessions retrieves all sessions, optionally filtered by workspace.
func (s *SessionStore) ListSessions(ctx context.Context, workspacePath string) ([]Session, error) {
	var rows *sql.Rows
	var err error

	if workspacePath != "" {
		rows, err = s.db.QueryContext(ctx,
			`SELECT 
				id, workspace_path, title, pinned, embedding, created_at, updated_at,
				worker_name, parent_session_id, model, provider, mode, thinking, summary,
				summarized_context, summary_token_count, summarized_at, last_summarized_msg_id
				FROM sessions WHERE workspace_path = ? ORDER BY updated_at DESC`,
			workspacePath,
		)
	} else {
		rows, err = s.db.QueryContext(ctx,
			`SELECT 
				id, workspace_path, title, pinned, embedding, created_at, updated_at,
				worker_name, parent_session_id, model, provider, mode, thinking, summary,
				summarized_context, summary_token_count, summarized_at, last_summarized_msg_id
				FROM sessions ORDER BY updated_at DESC`,
		)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to query sessions: %w", err)
	}
	defer rows.Close()

	var sessions []Session
	for rows.Next() {
		var session Session
		var pinnedInt int
		if err := rows.Scan(
			&session.ID, &session.WorkspacePath, &session.Title, &pinnedInt, &session.Embedding,
			&session.CreatedAt, &session.UpdatedAt, &session.WorkerName, &session.ParentSessionID,
			&session.Model, &session.Provider, &session.Mode, &session.Thinking, &session.Summary,
			&session.SummarizedContext, &session.SummaryTokenCount, &session.SummarizedAt,
			&session.LastSummarizedMsgID,
		); err != nil {
			return nil, fmt.Errorf("failed to scan session: %w", err)
		}
		session.Pinned = pinnedInt == 1
		sessions = append(sessions, session)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating sessions: %w", err)
	}

	return sessions, nil
}