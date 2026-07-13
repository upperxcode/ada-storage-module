package storage

// Core schema definitions for sessions and messages.
const (
	// sessionsTable creates the sessions table
	sessionsTable = `
	CREATE TABLE IF NOT EXISTS sessions (
		id TEXT PRIMARY KEY,
		workspace_path TEXT,
		title TEXT,
		pinned INTEGER DEFAULT 0,
		embedding BLOB,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		worker_name TEXT DEFAULT '',
		parent_session_id TEXT DEFAULT '',
		model TEXT DEFAULT '',
		provider TEXT DEFAULT '',
		mode TEXT DEFAULT 'ask',
		thinking TEXT DEFAULT '',
		summary TEXT,
		summarized_context TEXT DEFAULT '',
		summary_token_count INTEGER DEFAULT 0,
		summarized_at DATETIME,
		last_summarized_msg_id INTEGER DEFAULT 0
	);
	`

	// messagesTable creates the messages table
	messagesTable = `
	CREATE TABLE IF NOT EXISTS messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		session_id TEXT NOT NULL,
		role TEXT NOT NULL CHECK(role IN ('user','assistant','system','tool')),
		content TEXT NOT NULL,
		tokens INTEGER DEFAULT 0,
		time DATETIME DEFAULT CURRENT_TIMESTAMP,
		served_by TEXT,
		FOREIGN KEY(session_id) REFERENCES sessions(id) ON DELETE CASCADE
	);
	`

	// Indexes for messages
	idxMessagesSession = `CREATE INDEX IF NOT EXISTS idx_messages_session ON messages(session_id, time);`
	idxMessagesTime    = `CREATE INDEX IF NOT EXISTS idx_messages_time ON messages(time);`

	// appStateTable creates the app_state table
	appStateTable = `
	CREATE TABLE IF NOT EXISTS app_state (
		id INTEGER PRIMARY KEY CHECK (id=1),
		active_workspace_path TEXT,
		active_workspace_index INTEGER DEFAULT 0
	);
	`

	// routerConfigsTable creates the router_configs table
	routerConfigsTable = `
	CREATE TABLE IF NOT EXISTS router_configs (
		name TEXT PRIMARY KEY,
		endpoint TEXT NOT NULL,
		labels_json TEXT NOT NULL,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		router_type TEXT DEFAULT 'http-classifier',
		backend_model TEXT DEFAULT ''
	);
	`

	// schemaMigrationsTable tracks migration versions
	schemaMigrationsTable = `
	CREATE TABLE IF NOT EXISTS schema_migrations (
		version INTEGER PRIMARY KEY,
		applied_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);
	`
)