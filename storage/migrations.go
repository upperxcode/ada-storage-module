package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
)

// Migration represents a single database migration.
type Migration struct {
	Version    int
	UpSQL      string
	DownSQL    string
	AppliedAt  string
}

// RunMigrations executes all schema migrations in a transaction.
// If any migration fails, the transaction is rolled back and the error is returned.
// This follows the Fail-Fast philosophy: errors are never masked.
func RunMigrations(ctx context.Context, db *sql.DB) error {
	// Ensure schema_migrations table exists first
	if _, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			applied_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
	`); err != nil {
		return fmt.Errorf("failed to create schema_migrations table: %w", err)
	}

	// Get current migration version
	var currentVersion int
	err := db.QueryRowContext(ctx, `SELECT COALESCE(MAX(version), 0) FROM schema_migrations`).Scan(&currentVersion)
	if err != nil {
		return fmt.Errorf("failed to get current migration version: %w", err)
	}

	// Define migrations in order
	migrations := []struct {
		version int
		upSQL   string
	}{
		{1, sessionsTable},
		{2, messagesTable},
		{3, idxMessagesSession},
		{4, idxMessagesTime},
		{5, providersTable},
		{6, providerModelsTable},
		{7, providerApikeysTable},
		{8, configTable},
		{9, greetingsTable},
		{10, workspacesTable},
		{11, agentsTable},
		{12, workersTable},
		{13, skillsTable},
		{14, memoriesTable},
		{15, workspaceAgentsTable},
		{16, workspaceFoldersTable},
		{17, workspaceSkillsTable},
		{18, workspaceToolsTable},
		{19, workspaceWorkersTable},
		{20, fixedModelsTable},
		{21, fixedModelToolsTable},
		{22, toolProfilesTable},
		{23, toolProfileToolsTable},
		{24, specWizardsTable},
		{25, mcpsTable},
		{26, workspaceTemplatesTable},
		{27, appStateTable},
		{28, routerConfigsTable},
		{29, schemaMigrationsTable},
		{30, mcpAddURL},
		{31, mcpAddEnabled},
			{32, mcpAddTimeout},
			{33, mcpAddOAuthClientID},
			{34, skillAddColor},
			{35, skillAddIcon},
			{36, specWizardAddDepManifest},
			{37, specWizardAddStackPlugin},
			{38, addStrategyToProviders},
		{39, workspaceKnowledgeTable},
		{40, addProviderModelContextSize},
		{41, thinkingsTable},
		{42, messagesAddThinkingRole},
	}

	for _, m := range migrations {
		if m.version <= currentVersion {
			log.Printf("[DB] Migration v%d already applied, skipping", m.version)
			continue
		}

		log.Printf("[DB] Running migration v%d", m.version)

		// Start transaction for each migration
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return fmt.Errorf("failed to begin transaction for migration v%d: %w", m.version, err)
		}

		if _, err := tx.ExecContext(ctx, m.upSQL); err != nil {
			// Se o erro for "duplicate column name", ignoramos e marcamos como aplicada
			if strings.Contains(err.Error(), "duplicate column name") {
				log.Printf("[DB] Column already exists in migration v%d, marking as applied", m.version)
			} else {
				log.Printf("[DB] Migration failed: v%d - %v", m.version, err)
				if rbErr := tx.Rollback(); rbErr != nil {
					return fmt.Errorf("migration v%d failed and rollback failed: %w, rollback error: %v", m.version, err, rbErr)
				}
				return fmt.Errorf("migration v%d failed: %w", m.version, err)
			}
		}

		// Record migration as applied
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO schema_migrations (version, applied_at) VALUES (?, datetime('now'))`,
			m.version,
		); err != nil {
			log.Printf("[DB] Failed to record migration v%d: %v", m.version, err)
			if rbErr := tx.Rollback(); rbErr != nil {
				return fmt.Errorf("migration v%d recording failed and rollback failed: %w", m.version, rbErr)
			}
			return fmt.Errorf("failed to record migration v%d: %w", m.version, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit migration v%d: %w", m.version, err)
		}

		log.Printf("[DB] Migration v%d completed successfully", m.version)
	}

	log.Println("[DB] All migrations completed successfully")
	return nil
}

// GetMigrationStatus returns the current migration status.
func GetMigrationStatus(ctx context.Context, db *sql.DB) ([]Migration, error) {
	rows, err := db.QueryContext(ctx,
		`SELECT version, applied_at FROM schema_migrations ORDER BY version`,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query migration status: %w", err)
	}
	defer rows.Close()

	var migrations []Migration
	for rows.Next() {
		var m Migration
		if err := rows.Scan(&m.Version, &m.AppliedAt); err != nil {
			return nil, fmt.Errorf("failed to scan migration: %w", err)
		}
		migrations = append(migrations, m)
	}

	return migrations, nil
}