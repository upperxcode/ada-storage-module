package storage

// Other schema definitions.
const (
	// specWizardsTable creates the spec_wizards table
	specWizardsTable = `
	CREATE TABLE IF NOT EXISTS spec_wizards (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT,
		expert_language_plugin TEXT,
		prd TEXT,
		functional_requirements TEXT,
		non_functional_requirements TEXT,
		persistence TEXT,
		architecture TEXT,
		engineering_philosophies TEXT,
		design_patterns TEXT,
		data_patterns TEXT,
		stack_config TEXT,
		business_state_management TEXT,
		business_api_contract TEXT,
		business_customization_details TEXT,
		business_final_adjustments TEXT,
		business_architecture_recommendations TEXT,
		color TEXT DEFAULT '#3b82f6',
		icon TEXT DEFAULT '📝',
		architecture_health INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`

	// mcpsTable creates the mcps table
	mcpsTable = `
	CREATE TABLE IF NOT EXISTS mcps (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		nome TEXT NOT NULL,
		connect_type TEXT NOT NULL CHECK(connect_type IN ('websocket','url','cli_command')),
		command TEXT,
		arguments TEXT,
		environment TEXT,
		color TEXT DEFAULT '',
		icon TEXT DEFAULT ''
	);
	`

	// workspaceTemplatesTable creates the workspace_templates table
	workspaceTemplatesTable = `
	CREATE TABLE IF NOT EXISTS workspace_templates (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		description TEXT,
		personality TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`
)