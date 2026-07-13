package storage

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

// ============================================================
// INSERT TESTS - Workers
// ============================================================

func TestWorkers_InsertWorker(t *testing.T) {
	db := openRealDB(t)
	defer db.Close()

	ctx := context.Background()
	store := NewWorkerStore(db)

	worker := &Worker{
		Name:             "InsertWorker-" + timestamp(),
		Persona:          sql.NullString{String: "Test persona", Valid: true},
		ResponseLanguage: "portuguese",
		ConnectionType:   "cli_command",
		Command:          sql.NullString{String: "echo test", Valid: true},
		InheritanceFolders: true,
	}

	err := store.CreateWorker(ctx, worker)
	if err != nil {
		t.Fatalf("Failed to insert worker: %v", err)
	}

	// Verify insertion
	workers, err := store.ListWorkers(ctx)
	if err != nil {
		t.Fatalf("Failed to list workers: %v", err)
	}

	found := false
	for _, w := range workers {
		if w.Name == worker.Name {
			found = true
			if w.Persona.String != worker.Persona.String {
				t.Errorf("Persona mismatch")
			}
			if w.InheritanceFolders != worker.InheritanceFolders {
				t.Errorf("InheritanceFolders mismatch")
			}
			break
		}
	}

	if !found {
		t.Error("Inserted worker not found")
	}

	// Cleanup
	for _, w := range workers {
		if w.Name == worker.Name {
			store.DeleteWorker(ctx, w.ID)
			break
		}
	}
}

// ============================================================
// INSERT TESTS - Agents
// ============================================================

func TestAgents_InsertAgent(t *testing.T) {
	db := openRealDB(t)
	defer db.Close()

	ctx := context.Background()
	store := NewAgentStore(db)

	agent := &Agent{
		Name:         "InsertAgent-" + timestamp(),
		Description:  sql.NullString{String: "Test agent", Valid: true},
		Type:         AgentTypeResearch,
		MaxIteration: 5,
		Temperature:  0.7,
	}

	err := store.CreateAgent(ctx, agent)
	if err != nil {
		t.Fatalf("Failed to insert agent: %v", err)
	}

	// Cleanup
	agents, _ := store.ListAgents(ctx)
	for _, a := range agents {
		if a.Name == agent.Name {
			store.DeleteAgent(ctx, a.ID)
			break
		}
	}
}

// ============================================================
// INSERT TESTS - Skills
// ============================================================

func TestSkills_InsertSkill(t *testing.T) {
	db := openRealDB(t)
	defer db.Close()

	ctx := context.Background()
	store := NewSkillStore(db)

	skill := &Skill{
		Name:        "InsertSkill-" + timestamp(),
		Description: sql.NullString{String: "Test skill", Valid: true},
		Content:     "test content",
	}

	err := store.CreateSkill(ctx, skill)
	if err != nil {
		t.Fatalf("Failed to insert skill: %v", err)
	}

	// Cleanup
	skills, _ := store.ListSkills(ctx)
	for _, s := range skills {
		if s.Name == skill.Name {
			store.DeleteSkill(ctx, s.ID)
			break
		}
	}
}

// ============================================================
// INSERT TESTS - Providers
// ============================================================

func TestProviders_InsertProvider(t *testing.T) {
	db := openRealDB(t)
	defer db.Close()

	ctx := context.Background()
	store := NewProviderStore(db)

	provider := &Provider{
		Name: "InsertProvider-" + timestamp(),
	}

	err := store.CreateProvider(ctx, provider)
	if err != nil {
		t.Fatalf("Failed to insert provider: %v", err)
	}

	// Cleanup
	providers, _ := store.ListProviders(ctx)
	for _, p := range providers {
		if p.Name == provider.Name {
			store.DeleteProvider(ctx, p.ID)
			break
		}
	}
}

// ============================================================
// INSERT TESTS - Sessions
// ============================================================

func TestSessions_InsertSession(t *testing.T) {
	db := openRealDB(t)
	defer db.Close()

	ctx := context.Background()
	store := NewSessionStore(db)

	session := &Session{
		ID:    "insert-session-" + timestamp(),
		Title: sql.NullString{String: "Insert Test Session", Valid: true},
	}

	err := store.CreateSession(ctx, session)
	if err != nil {
		t.Fatalf("Failed to insert session: %v", err)
	}

	// Cleanup
	store.DeleteSession(ctx, session.ID)
}

// ============================================================
// INSERT TESTS - Config
// ============================================================

func TestConfig_InsertConfig(t *testing.T) {
	db := openRealDB(t)
	defer db.Close()

	ctx := context.Background()
	store := NewConfigStore(db)

	key := "insert_key_" + timestamp()

	err := store.SetConfig(ctx, key, "insert_value")
	if err != nil {
		t.Fatalf("Failed to insert config: %v", err)
	}

	// Cleanup
	store.DeleteConfig(ctx, key)
}