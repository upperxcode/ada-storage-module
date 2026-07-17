package storage

import (
	"context"
	"testing"
)

func TestWorkspaceTools_Create(t *testing.T) {
	db := openRealDB(t)
	defer db.Close()

	ctx := context.Background()
	store := NewWorkspaceToolsStore(db)
	wsStore := NewWorkspaceStore(db)

	workspace := &Workspace{
		Nome:  "ToolsTest-" + timestamp(),
		Path:  mustNullString("/tmp/tools-test"),
		Enabled: true,
	}
	wsStore.CreateWorkspace(ctx, workspace)

	workspaces, _ := wsStore.ListWorkspaces(ctx)
	var workspaceID int64
	for _, w := range workspaces {
		if w.Path.String == "/tmp/tools-test" {
			workspaceID = w.ID
			break
		}
	}

	if workspaceID == 0 {
		t.Fatal("Workspace not created")
	}

	tool := &WorkspaceTool{
		WorkspaceID: workspaceID,
		ToolName:    "read",
		Enabled:     true,
	}

	err := store.AddTool(ctx, tool)
	if err != nil {
		t.Fatalf("Failed to add workspace tool: %v", err)
	}

	tools, _ := store.ListTools(ctx, workspaceID)
	if len(tools) != 1 {
		t.Errorf("Expected 1 tool, got %d", len(tools))
	}

	store.DeleteAllTools(ctx, workspaceID)
	wsStore.DeleteWorkspace(ctx, workspaceID)
}

func TestWorkspaceTools_DeleteAll(t *testing.T) {
	db := openRealDB(t)
	defer db.Close()

	ctx := context.Background()
	store := NewWorkspaceToolsStore(db)
	wsStore := NewWorkspaceStore(db)

	workspace := &Workspace{
		Nome:  "ToolsDeleteTest-" + timestamp(),
		Path:  mustNullString("/tmp/tools-delete-test"),
		Enabled: true,
	}
	wsStore.CreateWorkspace(ctx, workspace)

	workspaces, _ := wsStore.ListWorkspaces(ctx)
	var workspaceID int64
	for _, w := range workspaces {
		if w.Path.String == "/tmp/tools-delete-test" {
			workspaceID = w.ID
			break
		}
	}

	store.AddTool(ctx, &WorkspaceTool{WorkspaceID: workspaceID, ToolName: "read", Enabled: true})
	store.AddTool(ctx, &WorkspaceTool{WorkspaceID: workspaceID, ToolName: "write", Enabled: true})

	store.DeleteAllTools(ctx, workspaceID)

	tools, _ := store.ListTools(ctx, workspaceID)
	if len(tools) != 0 {
		t.Errorf("Expected 0 tools after delete all, got %d", len(tools))
	}

	wsStore.DeleteWorkspace(ctx, workspaceID)
}