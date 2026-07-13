package storage

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

// ============================================================
// DELETE TESTS - Workers
// ============================================================

func TestWorkers_DeleteWorker(t *testing.T) {
	db := openRealDB(t)
	defer db.Close()

	ctx := context.Background()
	store := NewWorkerStore(db)

	// Insert first
	worker := &Worker{
		Name:             "DeleteWorker-" + timestamp(),
		ResponseLanguage: "portuguese",
		ConnectionType:   "cli_command",
	}

	err := store.CreateWorker(ctx, worker)
	if err != nil {
		t.Fatalf("Failed to create worker: %v", err)
	}

	// Get ID
	workers, _ := store.ListWorkers(ctx)
	var workerID int64
	for _, w := range workers {
		if w.Name == worker.Name {
			workerID = w.ID
			break
		}
	}

	// Delete
	err = store.DeleteWorker(ctx, workerID)
	if err != nil {
		t.Fatalf("Failed to delete worker: %v", err)
	}

	// Verify deleted
	_, err = store.GetWorker(ctx, workerID)
	if err == nil {
		t.Error("Worker should be deleted")
	}
	if !errors.Is(err, ErrWorkerNotFound) {
		t.Errorf("Expected ErrWorkerNotFound, got: %v", err)
	}
}

// ============================================================
// DELETE TESTS - Agents
// ============================================================

func TestAgents_DeleteAgent(t *testing.T) {
	db := openRealDB(t)
	defer db.Close()

	ctx := context.Background()
	store := NewAgentStore(db)

	// Insert first
	agent := &Agent{
		Name:         "DeleteAgent-" + timestamp(),
		Type:         AgentTypeResearch,
	}

	err := store.CreateAgent(ctx, agent)
	if err != nil {
		t.Fatalf("Failed to create agent: %v", err)
	}

	// Get ID
	agents, _ := store.ListAgents(ctx)
	var agentID int64
	for _, a := range agents {
		if a.Name == agent.Name {
			agentID = a.ID
			break
		}
	}

	// Delete
	err = store.DeleteAgent(ctx, agentID)
	if err != nil {
		t.Fatalf("Failed to delete agent: %v", err)
	}

	// Verify deleted
	_, err = store.GetAgent(ctx, agentID)
	if err == nil {
		t.Error("Agent should be deleted")
	}
	if !errors.Is(err, ErrAgentNotFound) {
		t.Errorf("Expected ErrAgentNotFound, got: %v", err)
	}
}

// ============================================================
// DELETE TESTS - Skills
// ============================================================

func TestSkills_DeleteSkill(t *testing.T) {
	db := openRealDB(t)
	defer db.Close()

	ctx := context.Background()
	store := NewSkillStore(db)

	// Insert first
	skill := &Skill{
		Name: "DeleteSkill-" + timestamp(),
	}

	err := store.CreateSkill(ctx, skill)
	if err != nil {
		t.Fatalf("Failed to create skill: %v", err)
	}

	// Get ID
	skills, _ := store.ListSkills(ctx)
	var skillID int64
	for _, s := range skills {
		if s.Name == skill.Name {
			skillID = s.ID
			break
		}
	}

	// Delete
	err = store.DeleteSkill(ctx, skillID)
	if err != nil {
		t.Fatalf("Failed to delete skill: %v", err)
	}

	// Verify deleted
	_, err = store.GetSkill(ctx, skillID)
	if err == nil {
		t.Error("Skill should be deleted")
	}
	if !errors.Is(err, ErrSkillNotFound) {
		t.Errorf("Expected ErrSkillNotFound, got: %v", err)
	}
}

// ============================================================
// DELETE TESTS - Providers
// ============================================================

func TestProviders_DeleteProvider(t *testing.T) {
	db := openRealDB(t)
	defer db.Close()

	ctx := context.Background()
	store := NewProviderStore(db)

	// Insert first
	provider := &Provider{
		Name: "DeleteProvider-" + timestamp(),
	}

	err := store.CreateProvider(ctx, provider)
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	// Get ID
	providers, _ := store.ListProviders(ctx)
	var providerID int64
	for _, p := range providers {
		if p.Name == provider.Name {
			providerID = p.ID
			break
		}
	}

	// Delete
	err = store.DeleteProvider(ctx, providerID)
	if err != nil {
		t.Fatalf("Failed to delete provider: %v", err)
	}

	// Verify deleted
	_, err = store.GetProvider(ctx, providerID)
	if err == nil {
		t.Error("Provider should be deleted")
	}
	if !errors.Is(err, ErrProviderNotFound) {
		t.Errorf("Expected ErrProviderNotFound, got: %v", err)
	}
}

// ============================================================
// DELETE TESTS - Sessions
// ============================================================

func TestSessions_DeleteSession(t *testing.T) {
	db := openRealDB(t)
	defer db.Close()

	ctx := context.Background()
	store := NewSessionStore(db)

	// Insert first
	sessionID := "delete-session-" + timestamp()
	err := store.CreateSession(ctx, &Session{
		ID:    sessionID,
		Title: sql.NullString{String: "To Delete", Valid: true},
	})
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	// Add a message
	store.SaveMessage(ctx, &Message{
		SessionID: sessionID,
		Role:      "user",
		Content:   "test message",
	})

	// Delete session (should cascade delete messages)
	err = store.DeleteSession(ctx, sessionID)
	if err != nil {
		t.Fatalf("Failed to delete session: %v", err)
	}

	// Verify session deleted
	_, err = store.GetSession(ctx, sessionID)
	if err == nil {
		t.Error("Session should be deleted")
	}
	if !errors.Is(err, ErrSessionNotFound) {
		t.Errorf("Expected ErrSessionNotFound, got: %v", err)
	}
}

// ============================================================
// DELETE TESTS - Config
// ============================================================

func TestConfig_DeleteConfig(t *testing.T) {
	db := openRealDB(t)
	defer db.Close()

	ctx := context.Background()
	store := NewConfigStore(db)

	// Insert first
	key := "delete_key_" + timestamp()
	store.SetConfig(ctx, key, "value")

	// Delete
	err := store.DeleteConfig(ctx, key)
	if err != nil {
		t.Fatalf("Failed to delete config: %v", err)
	}

	// Verify deleted
	_, err = store.GetConfig(ctx, key)
	if err == nil {
		t.Error("Config should be deleted")
	}
if !errors.Is(err, ErrConfigNotFound) {
			t.Errorf("Expected ErrConfigNotFound, got: %v", err)
		}
	}