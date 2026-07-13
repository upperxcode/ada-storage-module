package storage

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

// ============================================================
// UPDATE TESTS - Workers
// ============================================================

func TestWorkers_UpdateWorker(t *testing.T) {
	db := openRealDB(t)
	defer db.Close()

	ctx := context.Background()
	store := NewWorkerStore(db)

	// Insert first
	worker := &Worker{
		Name:             "UpdateWorker-" + timestamp(),
		Persona:          sql.NullString{String: "Original persona", Valid: true},
		ResponseLanguage: "portuguese",
		ConnectionType:   "cli_command",
		Command:          sql.NullString{String: "echo original", Valid: true},
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

	// Update
	updatedWorker := &Worker{
		ID:                  workerID,
		Name:                "UpdatedWorker-" + timestamp(),
		Persona:             sql.NullString{String: "Updated persona", Valid: true},
		ResponseLanguage:    "english",
		ConnectionType:      "url",
		Command:             sql.NullString{String: "curl http://api", Valid: true},
		InheritanceFolders:  true,
		InheritanceSkills:   true,
		InheritancePersona:  true,
		InheritanceTools:    true,
	}

	err = store.UpdateWorker(ctx, updatedWorker)
	if err != nil {
		t.Fatalf("Failed to update worker: %v", err)
	}

	// Verify update
	retrieved, err := store.GetWorker(ctx, workerID)
	if err != nil {
		t.Fatalf("Failed to get worker: %v", err)
	}

	if retrieved.Name != updatedWorker.Name {
		t.Errorf("Name not updated: expected %s, got %s", updatedWorker.Name, retrieved.Name)
	}

	if retrieved.ResponseLanguage != "english" {
		t.Errorf("ResponseLanguage not updated: got %s", retrieved.ResponseLanguage)
	}

	if !retrieved.InheritanceFolders {
		t.Error("InheritanceFolders should be true")
	}

	// Cleanup
	store.DeleteWorker(ctx, workerID)
}

// ============================================================
// UPDATE TESTS - Agents
// ============================================================

func TestAgents_UpdateAgent(t *testing.T) {
	db := openRealDB(t)
	defer db.Close()

	ctx := context.Background()
	store := NewAgentStore(db)

	// Insert first
	agent := &Agent{
		Name:         "UpdateAgent-" + timestamp(),
		Type:         AgentTypeExecutor,
		MaxIteration: 3,
		Temperature:  0.5,
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

	// Update
	updated := &Agent{
		ID:           agentID,
		Name:         "UpdatedAgent-" + timestamp(),
		Type:         AgentTypeReviewer,
		MaxIteration: 10,
		Temperature:  0.9,
	}

	err = store.UpdateAgent(ctx, updated)
	if err != nil {
		t.Fatalf("Failed to update agent: %v", err)
	}

	// Verify
	retrieved, _ := store.GetAgent(ctx, agentID)
	if retrieved.Name != updated.Name {
		t.Errorf("Name not updated")
	}
	if retrieved.Type != AgentTypeReviewer {
		t.Errorf("Type not updated")
	}

	// Cleanup
	store.DeleteAgent(ctx, agentID)
}

// ============================================================
// UPDATE TESTS - Skills
// ============================================================

func TestSkills_UpdateSkill(t *testing.T) {
	db := openRealDB(t)
	defer db.Close()

	ctx := context.Background()
	store := NewSkillStore(db)

	// Insert first
	skill := &Skill{
		Name:        "UpdateSkill-" + timestamp(),
		Description: sql.NullString{String: "Original", Valid: true},
		Content:     "original content",
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

	// Update
	updated := &Skill{
		ID:          skillID,
		Name:        "UpdatedSkill-" + timestamp(),
		Description: sql.NullString{String: "Updated description", Valid: true},
		Content:     "updated content",
	}

	err = store.UpdateSkill(ctx, updated)
	if err != nil {
		t.Fatalf("Failed to update skill: %v", err)
	}

	// Verify
	retrieved, _ := store.GetSkill(ctx, skillID)
	if retrieved.Name != updated.Name {
		t.Errorf("Name not updated")
	}

	// Cleanup
	store.DeleteSkill(ctx, skillID)
}

// ============================================================
// UPDATE TESTS - Sessions
// ============================================================

func TestSessions_UpdateSession(t *testing.T) {
	db := openRealDB(t)
	defer db.Close()

	ctx := context.Background()
	store := NewSessionStore(db)

	// Insert first
	session := &Session{
		ID:    "update-session-" + timestamp(),
		Title: sql.NullString{String: "Original Title", Valid: true},
	}

	err := store.CreateSession(ctx, session)
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	// Update
	updated := &Session{
		ID:    session.ID,
		Title: sql.NullString{String: "Updated Title", Valid: true},
	}

	err = store.UpdateSession(ctx, updated)
	if err != nil {
		t.Fatalf("Failed to update session: %v", err)
	}

	// Verify
	retrieved, _ := store.GetSession(ctx, session.ID)
	if retrieved.Title.String != "Updated Title" {
		t.Errorf("Title not updated")
	}

// Cleanup
		store.DeleteSession(ctx, session.ID)
	}