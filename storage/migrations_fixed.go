package storage

// Fixed models and tool profiles schema definitions.
const (
	// fixedModelsTable creates the fixed_models table
	fixedModelsTable = `
	CREATE TABLE IF NOT EXISTS fixed_models (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE NOT NULL,
		provider TEXT,
		model TEXT
	);
	`

	// fixedModelToolsTable creates the fixed_model_tools junction table
	fixedModelToolsTable = `
	CREATE TABLE IF NOT EXISTS fixed_model_tools (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		fixed_model_id INTEGER NOT NULL REFERENCES fixed_models(id) ON DELETE CASCADE,
		tool TEXT NOT NULL
	);
	`

	// toolProfilesTable creates the tool_profiles table
	toolProfilesTable = `
	CREATE TABLE IF NOT EXISTS tool_profiles (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		description TEXT,
		color TEXT DEFAULT '',
		icon TEXT DEFAULT ''
	);
	`

	// toolProfileToolsTable creates the tool_profile_tools junction table
	toolProfileToolsTable = `
	CREATE TABLE IF NOT EXISTS tool_profile_tools (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		profile_id INTEGER NOT NULL REFERENCES tool_profiles(id) ON DELETE CASCADE,
		tool_name TEXT NOT NULL
	);
	`
)