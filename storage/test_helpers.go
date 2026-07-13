package storage

import (
	"database/sql"
	"os"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// openRealDB opens the real ada-love-ide database for integration tests
func openRealDB(t *testing.T) *sql.DB {
	t.Helper()
	dbPath := "/home/john/.config/ada-love-ide/ada_love.db"
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Skipf("Database not found at %s", dbPath)
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	return db
}

// timestamp returns a unique suffix for test data
func timestamp() string {
	return time.Now().Format("20060102150405")
}

// timestampSuffix returns a unique suffix for test data (alias)
func timestampSuffix() string {
	return time.Now().Format("20060102150405")
}