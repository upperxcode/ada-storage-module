package storage

// Workspace-related schema definitions.
const (
	// workspacesTable creates the workspaces table
	workspacesTable = `
	CREATE TABLE IF NOT EXISTS workspaces (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		nome TEXT NOT NULL,
		description TEXT,
		path TEXT UNIQUE,
		max_prompt INTEGER NOT NULL DEFAULT 4096,
		max_content INTEGER NOT NULL DEFAULT 8192,
		"commit" BOOL NOT NULL DEFAULT 1,
		spec_provider TEXT,
		spec_wizard_id TEXT REFERENCES spec_wizards(id) ON DELETE SET NULL,
		personality TEXT,
		color TEXT DEFAULT '',
		icon TEXT DEFAULT '',
		summary TEXT DEFAULT '',
		enabled BOOL NOT NULL DEFAULT 1,
		max_prompt_send INTEGER NOT NULL DEFAULT 0,
		commit_changes BOOL NOT NULL DEFAULT 0,
		max_context_length INTEGER NOT NULL DEFAULT 0,
		embedding_model TEXT DEFAULT '',
		embedding_provider TEXT DEFAULT '',
		routing_rules TEXT DEFAULT ''
	);
	`

	// memoriesTable creates the memories table
	memoriesTable = `
	CREATE TABLE IF NOT EXISTS memories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		workspace_path TEXT,
		content TEXT NOT NULL,
		importance INTEGER DEFAULT 1,
		embedding BLOB NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`

	// workspaceAgentsTable creates the workspace_agents junction table
	workspaceAgentsTable = `
	CREATE TABLE IF NOT EXISTS workspace_agents (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		workspace_id INTEGER NOT NULL,
		agent_id INTEGER NOT NULL,
		enabled BOOL NOT NULL DEFAULT 1,
		FOREIGN KEY(workspace_id) REFERENCES workspaces(id) ON DELETE CASCADE,
		FOREIGN KEY(agent_id) REFERENCES agents(id) ON DELETE CASCADE
	);
	`

	// workspaceFoldersTable creates the workspace_folders junction table
	workspaceFoldersTable = `
	CREATE TABLE IF NOT EXISTS workspace_folders (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		workspace_id INTEGER NOT NULL,
		folder_path TEXT NOT NULL,
		FOREIGN KEY(workspace_id) REFERENCES workspaces(id) ON DELETE CASCADE
	);
	`

	// workspaceSkillsTable creates the workspace_skills junction table
	workspaceSkillsTable = `
	CREATE TABLE IF NOT EXISTS workspace_skills (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		workspace_id INTEGER NOT NULL,
		skill_id INTEGER NOT NULL,
		enabled BOOL NOT NULL DEFAULT 1,
		FOREIGN KEY(workspace_id) REFERENCES workspaces(id) ON DELETE CASCADE,
		FOREIGN KEY(skill_id) REFERENCES skills(id) ON DELETE CASCADE
	);
	`

	// workspaceToolsTable creates the workspace_tools junction table
	workspaceToolsTable = `
	CREATE TABLE IF NOT EXISTS workspace_tools (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		workspace_id INTEGER NOT NULL,
		tool_name TEXT NOT NULL,
		enabled BOOL NOT NULL DEFAULT 1,
		FOREIGN KEY(workspace_id) REFERENCES workspaces(id) ON DELETE CASCADE
	);
	`

	// workspaceWorkersTable creates the workspace_workers junction table
	workspaceWorkersTable = `
	CREATE TABLE IF NOT EXISTS workspace_workers (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		workspace_id INTEGER NOT NULL,
		worker_id INTEGER NOT NULL,
		enabled BOOL NOT NULL DEFAULT 1,
		FOREIGN KEY(workspace_id) REFERENCES workspaces(id) ON DELETE CASCADE,
		FOREIGN KEY(worker_id) REFERENCES workers(id) ON DELETE CASCADE
	);
	`
)