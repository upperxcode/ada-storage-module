package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// Init initializes the storage engine with migrations.
// This is the main entry point for the storage module.
func Init(ctx context.Context, dbPath string) (*StorageEngine, *SessionStore, *ConfigStore, error) {
	// Ensure the directory exists
	dir := filepath.Dir(dbPath)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, nil, nil, fmt.Errorf("failed to create database directory: %w", err)
		}
	}

	// Create storage engine
	engine, err := NewStorageEngine(dbPath)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create storage engine: %w", err)
	}

	// Run migrations (Fail-Fast: abort if migrations fail)
	if err := RunMigrations(ctx, engine.DB()); err != nil {
		engine.Close()
		return nil, nil, nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("[DB] Storage engine initialized successfully")

	// Create stores
	sessionStore := NewSessionStore(engine.DB())
	configStore := NewConfigStore(engine.DB())

	return engine, sessionStore, configStore, nil
}

// HealthCheck verifies the storage engine is healthy.
func HealthCheck(ctx context.Context, db *sql.DB) error {
	// Check database connectivity
	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	// Check that required tables exist
	tables := []string{
		"sessions", "messages", "providers", "provider_models",
		"config", "greetings", "workspaces", "agents", "workers", "skills",
	}
	for _, table := range tables {
		var exists bool
		if err := db.QueryRowContext(ctx,
			`SELECT EXISTS(SELECT 1 FROM sqlite_master WHERE type='table' AND name=?)`,
			table,
		).Scan(&exists); err != nil {
			return fmt.Errorf("failed to check table %s: %w", table, err)
		}
		if !exists {
			return fmt.Errorf("required table %s does not exist", table)
		}
	}

	return nil
}