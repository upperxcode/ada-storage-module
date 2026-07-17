package storage

import (
	"context"
	"testing"
)

func TestWorkspaceSkills_Create(t *testing.T) {
	db := openRealDB(t)
	defer db.Close()

	ctx := context.Background()
	store := NewWorkspaceSkillsStore(db)
	wsStore := NewWorkspaceStore(db)

	workspace := &Workspace{
		Nome:  "SkillsTest-" + timestamp(),
		Path:  mustNullString("/tmp/skills-test"),
		Enabled: true,
	}
	wsStore.CreateWorkspace(ctx, workspace)

	workspaces, _ := wsStore.ListWorkspaces(ctx)
	var workspaceID int64
	for _, w := range workspaces {
		if w.Path.String == "/tmp/skills-test" {
			workspaceID = w.ID
			break
		}
	}

	if workspaceID == 0 {
		t.Fatal("Workspace not created")
	}

	skill := &WorkspaceSkill{
		WorkspaceID: workspaceID,
		SkillID:     1,
		Enabled:     true,
	}

	err := store.Create(ctx, skill)
	if err != nil {
		t.Fatalf("Failed to create workspace skill: %v", err)
	}

	skills, _ := store.ListByWorkspace(ctx, workspaceID)
	if len(skills) != 1 {
		t.Errorf("Expected 1 skill, got %d", len(skills))
	}

	store.DeleteByWorkspace(ctx, workspaceID)
	wsStore.DeleteWorkspace(ctx, workspaceID)
}

func TestWorkspaceSkills_DeleteAll(t *testing.T) {
	db := openRealDB(t)
	defer db.Close()

	ctx := context.Background()
	store := NewWorkspaceSkillsStore(db)
	wsStore := NewWorkspaceStore(db)

	workspace := &Workspace{
		Nome:  "SkillsDeleteTest-" + timestamp(),
		Path:  mustNullString("/tmp/skills-delete-test"),
		Enabled: true,
	}
	wsStore.CreateWorkspace(ctx, workspace)

	workspaces, _ := wsStore.ListWorkspaces(ctx)
	var workspaceID int64
	for _, w := range workspaces {
		if w.Path.String == "/tmp/skills-delete-test" {
			workspaceID = w.ID
			break
		}
	}

	store.Create(ctx, &WorkspaceSkill{WorkspaceID: workspaceID, SkillID: 1, Enabled: true})
	store.Create(ctx, &WorkspaceSkill{WorkspaceID: workspaceID, SkillID: 2, Enabled: true})

	store.DeleteByWorkspace(ctx, workspaceID)

	skills, _ := store.ListByWorkspace(ctx, workspaceID)
	if len(skills) != 0 {
		t.Errorf("Expected 0 skills after delete all, got %d", len(skills))
	}

	wsStore.DeleteWorkspace(ctx, workspaceID)
}