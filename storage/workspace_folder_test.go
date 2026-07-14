package storage

import (
	"context"
	"database/sql"
	"errors"
	"testing"
)

// ============================================================
// WORKSPACE FOLDER TESTS
// ============================================================

func TestWorkspaceFolder_Create(t *testing.T) {
	db := openRealDB(t)
	defer db.Close()

	ctx := context.Background()
	store := NewWorkspaceFolderStore(db)

	// First create a workspace to get an ID
	wsStore := NewWorkspaceStore(db)
	workspace := &Workspace{
		Nome:  "TestWorkspace-" + timestamp(),
		Path:  mustNullString("/tmp/test-workspace"),
		Enabled: true,
	}
	wsStore.CreateWorkspace(ctx, workspace)

	// Get the workspace ID
	workspaces, _ := wsStore.ListWorkspaces(ctx)
	var workspaceID int64
	for _, w := range workspaces {
		if w.Path.String == "/tmp/test-workspace" {
			workspaceID = w.ID
			break
		}
	}

	if workspaceID == 0 {
		t.Fatal("Workspace not created")
	}

	// Create folder
	folder := &WorkspaceFolder{
		WorkspaceID: workspaceID,
		FolderPath:  "/tmp/test-workspace/src",
	}

	err := store.Create(ctx, folder)
	if err != nil {
		t.Fatalf("Failed to create workspace folder: %v", err)
	}

	// Verify
	folders, _ := store.ListByWorkspace(ctx, workspaceID)
	found := false
	for _, f := range folders {
		if f.FolderPath == folder.FolderPath {
			found = true
			break
		}
	}

	if !found {
		t.Error("Created folder not found")
	}

	// Cleanup
	store.Delete(ctx, workspaceID, folder.FolderPath)
	wsStore.DeleteWorkspace(ctx, workspaceID)
}

func TestWorkspaceFolder_Delete(t *testing.T) {
	db := openRealDB(t)
	defer db.Close()

	ctx := context.Background()
	store := NewWorkspaceFolderStore(db)
	wsStore := NewWorkspaceStore(db)

	// Create workspace
	workspace := &Workspace{
		Nome:  "DeleteTest-" + timestamp(),
		Path:  mustNullString("/tmp/delete-test"),
		Enabled: true,
	}
	wsStore.CreateWorkspace(ctx, workspace)

	workspaces, _ := wsStore.ListWorkspaces(ctx)
	var workspaceID int64
	for _, w := range workspaces {
		if w.Path.String == "/tmp/delete-test" {
			workspaceID = w.ID
			break
		}
	}

	// Create folder
	folder := &WorkspaceFolder{
		WorkspaceID: workspaceID,
		FolderPath:  "/tmp/delete-test/folder",
	}
	store.Create(ctx, folder)

	// Delete folder
	err := store.Delete(ctx, workspaceID, folder.FolderPath)
	if err != nil {
		t.Fatalf("Failed to delete workspace folder: %v", err)
	}

	// Verify deleted
	err = store.Delete(ctx, workspaceID, folder.FolderPath)
	if err == nil {
		t.Error("Folder should be deleted")
	}
	if !errors.Is(err, ErrWorkspaceFolderNotFound) {
		t.Errorf("Expected ErrWorkspaceFolderNotFound, got: %v", err)
	}

	// Cleanup
	wsStore.DeleteWorkspace(ctx, workspaceID)
}

func TestWorkspaceFolder_ListByWorkspace(t *testing.T) {
	db := openRealDB(t)
	defer db.Close()

	ctx := context.Background()
	store := NewWorkspaceFolderStore(db)
	wsStore := NewWorkspaceStore(db)

	// Create workspace
	workspace := &Workspace{
		Nome:  "ListTest-" + timestamp(),
		Path:  mustNullString("/tmp/list-test"),
		Enabled: true,
	}
	wsStore.CreateWorkspace(ctx, workspace)

	workspaces, _ := wsStore.ListWorkspaces(ctx)
	var workspaceID int64
	for _, w := range workspaces {
		if w.Path.String == "/tmp/list-test" {
			workspaceID = w.ID
			break
		}
	}

	// Create folders
	folders := []string{"/tmp/list-test/folder1", "/tmp/list-test/folder2"}
	for _, path := range folders {
		store.Create(ctx, &WorkspaceFolder{
			WorkspaceID: workspaceID,
			FolderPath:  path,
		})
	}

	// List folders
	retrieved, err := store.ListByWorkspace(ctx, workspaceID)
	if err != nil {
		t.Fatalf("Failed to list workspace folders: %v", err)
	}

	if len(retrieved) != 2 {
		t.Errorf("Expected 2 folders, got %d", len(retrieved))
	}

	// Cleanup
	store.DeleteAllByWorkspace(ctx, workspaceID)
	wsStore.DeleteWorkspace(ctx, workspaceID)
}

func TestWorkspaceFolder_DeleteAllByWorkspace(t *testing.T) {
	db := openRealDB(t)
	defer db.Close()

	ctx := context.Background()
	store := NewWorkspaceFolderStore(db)
	wsStore := NewWorkspaceStore(db)

	// Create workspace
	workspace := &Workspace{
		Nome:  "DeleteAllTest-" + timestamp(),
		Path:  mustNullString("/tmp/delete-all-test"),
		Enabled: true,
	}
	wsStore.CreateWorkspace(ctx, workspace)

	workspaces, _ := wsStore.ListWorkspaces(ctx)
	var workspaceID int64
	for _, w := range workspaces {
		if w.Path.String == "/tmp/delete-all-test" {
			workspaceID = w.ID
			break
		}
	}

	// Create folders
	store.Create(ctx, &WorkspaceFolder{WorkspaceID: workspaceID, FolderPath: "/tmp/test1"})
	store.Create(ctx, &WorkspaceFolder{WorkspaceID: workspaceID, FolderPath: "/tmp/test2"})

	// Delete all
	err := store.DeleteAllByWorkspace(ctx, workspaceID)
	if err != nil {
		t.Fatalf("Failed to delete all workspace folders: %v", err)
	}

	// Verify empty
	folders, _ := store.ListByWorkspace(ctx, workspaceID)
	if len(folders) != 0 {
		t.Errorf("Expected 0 folders after delete all, got %d", len(folders))
	}

	// Cleanup
	wsStore.DeleteWorkspace(ctx, workspaceID)
}

// Helper function
func mustNullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: true}
}