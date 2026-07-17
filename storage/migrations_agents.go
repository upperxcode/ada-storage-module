package storage

// Agent-related schema definitions.
const (
	// agentsTable creates the agents table
	agentsTable = `
	CREATE TABLE IF NOT EXISTS agents (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT,
		type TEXT NOT NULL CHECK(type IN ('executor','delegator','reviewer','research')),
		provider_id INTEGER,
		model_id INTEGER,
		max_iteration INTEGER NOT NULL DEFAULT 10,
		temperature REAL NOT NULL DEFAULT 0.7,
		system_prompt TEXT,
		color TEXT DEFAULT '',
		icon TEXT DEFAULT ''
	);
	`

	// workersTable creates the workers table
	workersTable = `
	CREATE TABLE IF NOT EXISTS workers (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		persona TEXT,
		response_language TEXT DEFAULT 'portuguese',
		connection_type TEXT NOT NULL,
		command TEXT,
		arguments TEXT,
		environment TEXT,
		inheritance_folders BOOL NOT NULL DEFAULT 0,
		inheritance_skills BOOL NOT NULL DEFAULT 0,
		inheritance_persona BOOL NOT NULL DEFAULT 0,
		inheritance_knowledge BOOL NOT NULL DEFAULT 0,
		inheritance_tools BOOL NOT NULL DEFAULT 0,
		color TEXT DEFAULT '',
		icon TEXT DEFAULT ''
	);
	`

	// skillsTable creates the skills table
	skillsTable = `
	CREATE TABLE IF NOT EXISTS skills (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT,
		tags TEXT,
		content TEXT NOT NULL,
		color TEXT DEFAULT '',
		icon TEXT DEFAULT ''
	);
	`

	// skillAddColor adds color column to skills for existing DBs
	skillAddColor = `
	ALTER TABLE skills ADD COLUMN color TEXT DEFAULT '';
	`

	// skillAddIcon adds icon column to skills for existing DBs
	skillAddIcon = `
	ALTER TABLE skills ADD COLUMN icon TEXT DEFAULT '';
	`
)