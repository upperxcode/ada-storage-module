package storage

import (
	"context"
	"testing"
)

func TestWorkspaceKnowledge_Create(t *testing.T) {
	db := openRealDB(t)
	defer db.Close()

	ctx := context.Background()
	store := NewWorkspaceKnowledgeStore(db)
	wsStore := NewWorkspaceStore(db)

	workspace := &Workspace{
		Nome:  "KnowledgeTest-" + timestamp(),
		Path:  mustNullString("/tmp/knowledge-test"),
		Enabled: true,
	}
	wsStore.CreateWorkspace(ctx, workspace)

	workspaces, _ := wsStore.ListWorkspaces(ctx)
	var workspaceID int64
	for _, w := range workspaces {
		if w.Path.String == "/tmp/knowledge-test" {
			workspaceID = w.ID
			break
		}
	}

	if workspaceID == 0 {
		t.Fatal("Workspace not created")
	}

	knowledge := &WorkspaceKnowledge{
		WorkspaceID:   workspaceID,
		KnowledgeItem: "/tmp/knowledge-test/doc1.md",
	}

	err := store.AddKnowledge(ctx, knowledge)
	if err != nil {
		t.Fatalf("Failed to create workspace knowledge: %v", err)
	}

	knowledgeItems, _ := store.ListKnowledge(ctx, workspaceID)
	if len(knowledgeItems) != 1 {
		t.Errorf("Expected 1 knowledge item, got %d", len(knowledgeItems))
	}

	store.DeleteAllKnowledge(ctx, workspaceID)
	wsStore.DeleteWorkspace(ctx, workspaceID)
}

func TestWorkspaceKnowledge_DeleteAll(t *testing.T) {
	db := openRealDB(t)
	defer db.Close()

	ctx := context.Background()
	store := NewWorkspaceKnowledgeStore(db)
	wsStore := NewWorkspaceStore(db)

	workspace := &Workspace{
		Nome:  "KnowledgeDeleteTest-" + timestamp(),
		Path:  mustNullString("/tmp/knowledge-delete-test"),
		Enabled: true,
	}
	wsStore.CreateWorkspace(ctx, workspace)

	workspaces, _ := wsStore.ListWorkspaces(ctx)
	var workspaceID int64
	for _, w := range workspaces {
		if w.Path.String == "/tmp/knowledge-delete-test" {
			workspaceID = w.ID
			break
		}
	}

	store.AddKnowledge(ctx, &WorkspaceKnowledge{WorkspaceID: workspaceID, KnowledgeItem: "item1"})
	store.AddKnowledge(ctx, &WorkspaceKnowledge{WorkspaceID: workspaceID, KnowledgeItem: "item2"})

	store.DeleteAllKnowledge(ctx, workspaceID)

	items, _ := store.ListKnowledge(ctx, workspaceID)
	if len(items) != 0 {
		t.Errorf("Expected 0 items after delete all, got %d", len(items))
	}

	wsStore.DeleteWorkspace(ctx, workspaceID)
}