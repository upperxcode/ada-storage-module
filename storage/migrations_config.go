package storage

// Config-related schema definitions.
const (
	// configTable creates the config table
	configTable = `
	CREATE TABLE IF NOT EXISTS config (
		key TEXT PRIMARY KEY,
		value TEXT NOT NULL
	);
	`

	// greetingsTable creates the greetings table
	greetingsTable = `
	CREATE TABLE IF NOT EXISTS greetings (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		keyword TEXT UNIQUE NOT NULL,
		language TEXT NOT NULL,
		response TEXT NOT NULL
	);
	`
)