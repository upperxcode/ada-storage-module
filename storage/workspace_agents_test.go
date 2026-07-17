package storage

import (
	"context"
	"testing"
)

func TestWorkspaceAgents_Create(t *testing.T) {
	db := openRealDB(t)
	defer db.Close()

	ctx := context.Background()
	store := NewWorkspaceAgentsStore(db)
	wsStore := NewWorkspaceStore(db)

	workspace := &Workspace{
		Nome:  "AgentsTest-" + timestamp(),
		Path:  mustNullString("/tmp/agents-test"),
		Enabled: true,
	}
	wsStore.CreateWorkspace(ctx, workspace)

	workspaces, _ := wsStore.ListWorkspaces(ctx)
	var workspaceID int64
	for _, w := range workspaces {
		if w.Path.String == "/tmp/agents-test" {
			workspaceID = w.ID
			break
		}
	}

	if workspaceID == 0 {
		t.Fatal("Workspace not created")
	}

	agent := &WorkspaceAgent{
		WorkspaceID: workspaceID,
		AgentID:     1,
		Enabled:     true,
	}

	err := store.AddAgent(ctx, agent)
	if err != nil {
		t.Fatalf("Failed to add workspace agent: %v", err)
	}

	agents, _ := store.ListAgents(ctx, workspaceID)
	if len(agents) != 1 {
		t.Errorf("Expected 1 agent, got %d", len(agents))
	}

	store.DeleteAllAgents(ctx, workspaceID)
	wsStore.DeleteWorkspace(ctx, workspaceID)
}

func TestWorkspaceAgents_DeleteAll(t *testing.T) {
	db := openRealDB(t)
	defer db.Close()

	ctx := context.Background()
	store := NewWorkspaceAgentsStore(db)
	wsStore := NewWorkspaceStore(db)

	workspace := &Workspace{
		Nome:  "AgentsDeleteTest-" + timestamp(),
		Path:  mustNullString("/tmp/agents-delete-test"),
		Enabled: true,
	}
	wsStore.CreateWorkspace(ctx, workspace)

	workspaces, _ := wsStore.ListWorkspaces(ctx)
	var workspaceID int64
	for _, w := range workspaces {
		if w.Path.String == "/tmp/agents-delete-test" {
			workspaceID = w.ID
			break
		}
	}

	store.AddAgent(ctx, &WorkspaceAgent{WorkspaceID: workspaceID, AgentID: 1, Enabled: true})
	store.AddAgent(ctx, &WorkspaceAgent{WorkspaceID: workspaceID, AgentID: 2, Enabled: true})

	store.DeleteAllAgents(ctx, workspaceID)

	agents, _ := store.ListAgents(ctx, workspaceID)
	if len(agents) != 0 {
		t.Errorf("Expected 0 agents after delete all, got %d", len(agents))
	}

	wsStore.DeleteWorkspace(ctx, workspaceID)
}