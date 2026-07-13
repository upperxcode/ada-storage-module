package storage

// Provider-related schema definitions.
const (
	// providersTable creates the providers table
	providersTable = `
	CREATE TABLE IF NOT EXISTS providers (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		api_url TEXT,
		connection_types TEXT,
		color TEXT DEFAULT '',
		icon TEXT DEFAULT ''
	);
	`

	// providerModelsTable creates the provider_models table
	providerModelsTable = `
	CREATE TABLE IF NOT EXISTS provider_models (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		provider_id INTEGER NOT NULL,
		model TEXT NOT NULL,
		free BOOL NOT NULL DEFAULT 0,
		thinking BOOL NOT NULL DEFAULT 0,
		tool BOOL NOT NULL DEFAULT 0,
		embedding BOOL NOT NULL DEFAULT 0,
		vision BOOL NOT NULL DEFAULT 0,
		health INTEGER NOT NULL DEFAULT 100,
		UNIQUE(provider_id, model),
		FOREIGN KEY(provider_id) REFERENCES providers(id) ON DELETE CASCADE
	);
	`

	// providerApikeysTable creates the provider_apikeys table
	providerApikeysTable = `
	CREATE TABLE IF NOT EXISTS provider_apikeys (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		provider_id INTEGER NOT NULL,
		apikey TEXT NOT NULL,
		UNIQUE(provider_id, apikey),
		FOREIGN KEY(provider_id) REFERENCES providers(id) ON DELETE CASCADE
	);
	`
)