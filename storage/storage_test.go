package storage

import (
	"context"
	"database/sql"
	"errors"
	"testing"
)

func TestInit(t *testing.T) {
	ctx := context.Background()
	engine, chatStore, configStore, err := Init(ctx, "file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	defer engine.Close()

	if engine == nil {
		t.Error("Engine should not be nil")
	}
	if chatStore == nil {
		t.Error("ChatStore should not be nil")
	}
	if configStore == nil {
		t.Error("ConfigStore should not be nil")
	}
}

func TestSessionStore_CreateAndGetSession(t *testing.T) {
	ctx := context.Background()
	engine, _, _, err := Init(ctx, "file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	defer engine.Close()

	store := NewSessionStore(engine.DB())
	sessionID := GenerateSessionID()

	err = store.CreateSession(ctx, &Session{
		ID:    sessionID,
		Title: sql.NullString{String: "Test Session", Valid: true},
	})
	if err != nil {
		t.Fatalf("CreateSession failed: %v", err)
	}

	session, err := store.GetSession(ctx, sessionID)
	if err != nil {
		t.Fatalf("GetSession failed: %v", err)
	}
	if session.ID != sessionID {
		t.Errorf("Expected ID %s, got %s", sessionID, session.ID)
	}
}

func TestSessionStore_GetSession_NotFound(t *testing.T) {
	ctx := context.Background()
	engine, _, _, err := Init(ctx, "file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	defer engine.Close()

	store := NewSessionStore(engine.DB())

	_, err = store.GetSession(ctx, "nonexistent-id")
	if err == nil {
		t.Error("Expected error for non-existent session")
	}
	if !errors.Is(err, ErrSessionNotFound) {
		t.Errorf("Expected ErrSessionNotFound, got: %v", err)
	}
}

func TestSessionStore_SaveAndRetrieveMessages(t *testing.T) {
	ctx := context.Background()
	engine, _, _, err := Init(ctx, "file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	defer engine.Close()

	store := NewSessionStore(engine.DB())
	sessionID := GenerateSessionID()

	err = store.CreateSession(ctx, &Session{
		ID:    sessionID,
		Title: sql.NullString{String: "Test Session", Valid: true},
	})
	if err != nil {
		t.Fatalf("CreateSession failed: %v", err)
	}

	// Add messages
	err = store.SaveMessage(ctx, &Message{
		SessionID: sessionID,
		Role:      "user",
		Content:   "Hello",
	})
	if err != nil {
		t.Fatalf("SaveMessage failed: %v", err)
	}

	err = store.SaveMessage(ctx, &Message{
		SessionID: sessionID,
		Role:      "assistant",
		Content:   "Hi there!",
	})
	if err != nil {
		t.Fatalf("SaveMessage failed: %v", err)
	}

	// Retrieve messages
	messages, err := store.GetMessages(ctx, sessionID)
	if err != nil {
		t.Fatalf("GetMessages failed: %v", err)
	}
	if len(messages) != 2 {
		t.Errorf("Expected 2 messages, got %d", len(messages))
	}
}

func TestProviderStore_CreateAndGetProvider(t *testing.T) {
	ctx := context.Background()
	engine, _, _, err := Init(ctx, "file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	defer engine.Close()

	store := NewProviderStore(engine.DB())
	provider := &Provider{
		Name: "openai",
	}

	err = store.CreateProvider(ctx, provider)
	if err != nil {
		t.Fatalf("CreateProvider failed: %v", err)
	}

	retrieved, err := store.GetProviderByName(ctx, "openai")
	if err != nil {
		t.Fatalf("GetProviderByName failed: %v", err)
	}
	if retrieved.Name != "openai" {
		t.Errorf("Expected name 'openai', got '%s'", retrieved.Name)
	}
}

func TestProviderStore_GetProvider_NotFound(t *testing.T) {
	ctx := context.Background()
	engine, _, _, err := Init(ctx, "file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	defer engine.Close()

	store := NewProviderStore(engine.DB())

	_, err = store.GetProvider(ctx, 999)
	if err == nil {
		t.Error("Expected error for non-existent provider")
	}
	if !errors.Is(err, ErrProviderNotFound) {
		t.Errorf("Expected ErrProviderNotFound, got: %v", err)
	}
}

func TestConfigStore_SetAndGetConfig(t *testing.T) {
	ctx := context.Background()
	engine, _, _, err := Init(ctx, "file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	defer engine.Close()

	store := NewConfigStore(engine.DB())

	err = store.SetConfig(ctx, "test-key", "test-value")
	if err != nil {
		t.Fatalf("SetConfig failed: %v", err)
	}

	value, err := store.GetConfig(ctx, "test-key")
	if err != nil {
		t.Fatalf("GetConfig failed: %v", err)
	}
	if value != "test-value" {
		t.Errorf("Expected 'test-value', got '%s'", value)
	}
}

func TestConfigStore_GetConfig_NotFound(t *testing.T) {
	ctx := context.Background()
	engine, _, _, err := Init(ctx, "file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	defer engine.Close()

	store := NewConfigStore(engine.DB())

	_, err = store.GetConfig(ctx, "nonexistent-key")
	if err == nil {
		t.Error("Expected error for non-existent key")
	}
	if !errors.Is(err, ErrConfigNotFound) {
		t.Errorf("Expected ErrConfigNotFound, got: %v", err)
	}
}

func TestGreetingStore_CreateAndGetGreeting(t *testing.T) {
	ctx := context.Background()
	engine, _, _, err := Init(ctx, "file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	defer engine.Close()

	store := NewGreetingStore(engine.DB())

	err = store.CreateGreeting(ctx, &Greeting{
		Keyword:  "hello",
		Language: "en",
		Response: "Hello! How can I help you?",
	})
	if err != nil {
		t.Fatalf("CreateGreeting failed: %v", err)
	}

	greeting, err := store.GetGreeting(ctx, "hello", "en")
	if err != nil {
		t.Fatalf("GetGreeting failed: %v", err)
	}
	if greeting.Response != "Hello! How can I help you?" {
		t.Errorf("Unexpected response: %s", greeting.Response)
	}
}

func TestHealthCheck(t *testing.T) {
	ctx := context.Background()
	engine, _, _, err := Init(ctx, "file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	defer engine.Close()

	err = HealthCheck(ctx, engine.DB())
	if err != nil {
		t.Errorf("HealthCheck failed: %v", err)
	}
}

func TestMigrationVersionTracking(t *testing.T) {
	ctx := context.Background()
	engine, _, _, err := Init(ctx, "file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	defer engine.Close()

	// Check migration status
	status, err := GetMigrationStatus(ctx, engine.DB())
	if err != nil {
		t.Fatalf("GetMigrationStatus failed: %v", err)
	}

	if len(status) == 0 {
		t.Error("Expected at least one migration to be recorded")
	}
}