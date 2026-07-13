package storage

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

// StorageEngine manages the SQLite database connection and pool.
type StorageEngine struct {
	db *sql.DB
	mu sync.Mutex
}

// NewStorageEngine creates a new SQLite storage engine with WAL mode and foreign keys enabled.
func NewStorageEngine(dsn string) (*StorageEngine, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Configure connection pool for local SQLite
	db.SetMaxOpenConns(1)    // SQLite handles concurrency internally
	db.SetMaxIdleConns(1)    // Keep one idle connection
	db.SetConnMaxLifetime(0) // Never recycle connections

	engine := &StorageEngine{db: db}

	// Enable Write-Ahead Logging for concurrent reads during writes
	if err := engine.enableWAL(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to enable WAL mode: %w", err)
	}

	// Enable foreign key constraints
	if err := engine.enableForeignKeys(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	return engine, nil
}

// enableWAL enables Write-Ahead Logging for better concurrency.
func (e *StorageEngine) enableWAL() error {
	_, err := e.db.Exec("PRAGMA journal_mode=WAL;")
	return err
}

// enableForeignKeys enables foreign key constraint enforcement.
func (e *StorageEngine) enableForeignKeys() error {
	_, err := e.db.Exec("PRAGMA foreign_keys=ON;")
	return err
}

// DB returns the underlying database connection.
func (e *StorageEngine) DB() *sql.DB {
	return e.db
}

// Close safely closes the database connection.
func (e *StorageEngine) Close() error {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.db == nil {
		return nil
	}
	return e.db.Close()
}

// Ping verifies the database connection is still alive.
func (e *StorageEngine) Ping() error {
	return e.db.Ping()
}

// BeginTx starts a new transaction with context.
func (e *StorageEngine) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return e.db.BeginTx(ctx, nil)
}
