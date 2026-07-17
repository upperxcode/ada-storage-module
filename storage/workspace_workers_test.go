package storage

import (
	"context"
	"testing"
)

func TestWorkspaceWorkers_Create(t *testing.T) {
	db := openRealDB(t)
	defer db.Close()

	ctx := context.Background()
	store := NewWorkspaceWorkersStore(db)
	wsStore := NewWorkspaceStore(db)

	workspace := &Workspace{
		Nome:  "WorkersTest-" + timestamp(),
		Path:  mustNullString("/tmp/workers-test"),
		Enabled: true,
	}
	wsStore.CreateWorkspace(ctx, workspace)

	workspaces, _ := wsStore.ListWorkspaces(ctx)
	var workspaceID int64
	for _, w := range workspaces {
		if w.Path.String == "/tmp/workers-test" {
			workspaceID = w.ID
			break
		}
	}

	if workspaceID == 0 {
		t.Fatal("Workspace not created")
	}

	worker := &WorkspaceWorker{
		WorkspaceID: workspaceID,
		WorkerID:    1,
		Enabled:     true,
	}

	err := store.AddWorker(ctx, worker)
	if err != nil {
		t.Fatalf("Failed to add workspace worker: %v", err)
	}

	workers, _ := store.ListWorkers(ctx, workspaceID)
	if len(workers) != 1 {
		t.Errorf("Expected 1 worker, got %d", len(workers))
	}

	store.DeleteAllWorkers(ctx, workspaceID)
	wsStore.DeleteWorkspace(ctx, workspaceID)
}

func TestWorkspaceWorkers_DeleteAll(t *testing.T) {
	db := openRealDB(t)
	defer db.Close()

	ctx := context.Background()
	store := NewWorkspaceWorkersStore(db)
	wsStore := NewWorkspaceStore(db)

	workspace := &Workspace{
		Nome:  "WorkersDeleteTest-" + timestamp(),
		Path:  mustNullString("/tmp/workers-delete-test"),
		Enabled: true,
	}
	wsStore.CreateWorkspace(ctx, workspace)

	workspaces, _ := wsStore.ListWorkspaces(ctx)
	var workspaceID int64
	for _, w := range workspaces {
		if w.Path.String == "/tmp/workers-delete-test" {
			workspaceID = w.ID
			break
		}
	}

	store.AddWorker(ctx, &WorkspaceWorker{WorkspaceID: workspaceID, WorkerID: 1, Enabled: true})
	store.AddWorker(ctx, &WorkspaceWorker{WorkspaceID: workspaceID, WorkerID: 2, Enabled: true})

	store.DeleteAllWorkers(ctx, workspaceID)

	workers, _ := store.ListWorkers(ctx, workspaceID)
	if len(workers) != 0 {
		t.Errorf("Expected 0 workers after delete all, got %d", len(workers))
	}

	wsStore.DeleteWorkspace(ctx, workspaceID)
}