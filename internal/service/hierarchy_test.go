package service

import (
	"testing"
	"time"

	"github.com/lucasefe/ontop/internal/models"
)

// Helper to create test task
func makeTask(id, title string, priority int, parentID *string) *models.Task {
	return &models.Task{
		ID:        id,
		Title:     title,
		Priority:  priority,
		Column:    "inbox",
		Progress:  0,
		ParentID:  parentID,
		Archived:  false,
		Tags:      []string{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// Helper to create parent ID pointer
func ptr(s string) *string {
	return &s
}

// TestBuildFlatHierarchy_EmptyInput tests empty slice handling
func TestBuildFlatHierarchy_EmptyInput(t *testing.T) {
	result := BuildFlatHierarchy([]*models.Task{}, SortByPriority)
	if len(result) != 0 {
		t.Errorf("Expected empty result, got %d items", len(result))
	}
}

// TestBuildFlatHierarchy_OnlyParents tests no subtasks case
func TestBuildFlatHierarchy_OnlyParents(t *testing.T) {
	tasks := []*models.Task{
		makeTask("P1", "Parent 1", 2, nil),
		makeTask("P2", "Parent 2", 1, nil),
	}

	result := BuildFlatHierarchy(tasks, SortByPriority)

	if len(result) != 2 {
		t.Fatalf("Expected 2 items, got %d", len(result))
	}

	// Should be sorted by priority (P2 priority 1, then P1 priority 2)
	if result[0].Task.ID != "P2" {
		t.Errorf("Expected P2 first (priority 1), got %s", result[0].Task.ID)
	}
	if result[1].Task.ID != "P1" {
		t.Errorf("Expected P1 second (priority 2), got %s", result[1].Task.ID)
	}

	// Both should be top-level
	if result[0].IsSubtask || result[1].IsSubtask {
		t.Error("Expected no subtasks")
	}
	if result[0].Indentation != 0 || result[1].Indentation != 0 {
		t.Error("Expected zero indentation for parents")
	}
}

// TestBuildFlatHierarchy_OnlySubtasks tests orphaned subtasks (all treated as top-level)
func TestBuildFlatHierarchy_OnlySubtasks(t *testing.T) {
	tasks := []*models.Task{
		makeTask("S1", "Subtask 1", 1, ptr("MISSING")),
		makeTask("S2", "Subtask 2", 2, ptr("MISSING")),
	}

	result := BuildFlatHierarchy(tasks, SortByPriority)

	if len(result) != 2 {
		t.Fatalf("Expected 2 items, got %d", len(result))
	}

	// Orphaned subtasks should be treated as top-level
	if result[0].IsSubtask || result[1].IsSubtask {
		t.Error("Expected orphaned subtasks treated as top-level")
	}
	if result[0].Indentation != 0 || result[1].Indentation != 0 {
		t.Error("Expected zero indentation for orphaned subtasks")
	}
}

// TestBuildFlatHierarchy_MixedParentsAndSubtasks tests standard hierarchy case
func TestBuildFlatHierarchy_MixedParentsAndSubtasks(t *testing.T) {
	tasks := []*models.Task{
		makeTask("S1", "Subtask 1", 1, ptr("P1")),
		makeTask("P1", "Parent 1", 2, nil),
		makeTask("S2", "Subtask 2", 3, ptr("P1")),
		makeTask("P2", "Parent 2", 1, nil),
	}

	result := BuildFlatHierarchy(tasks, SortByPriority)

	if len(result) != 4 {
		t.Fatalf("Expected 4 items, got %d", len(result))
	}

	// Expected order: P2 (priority 1), P1 (priority 2), S1 (subtask priority 1), S2 (subtask priority 3)
	expected := []string{"P2", "P1", "S1", "S2"}
	for i, exp := range expected {
		if result[i].Task.ID != exp {
			t.Errorf("Position %d: expected %s, got %s", i, exp, result[i].Task.ID)
		}
	}

	// Check subtask flags
	if result[0].IsSubtask || result[1].IsSubtask {
		t.Error("Parents should not be marked as subtasks")
	}
	if !result[2].IsSubtask || !result[3].IsSubtask {
		t.Error("Subtasks should be marked as subtasks")
	}

	// Check indentation
	if result[0].Indentation != 0 || result[1].Indentation != 0 {
		t.Error("Parents should have zero indentation")
	}
	if result[2].Indentation != 2 || result[3].Indentation != 2 {
		t.Error("Subtasks should have 2-space indentation")
	}

	// Check parent title
	if result[2].ParentTitle != "Parent 1" {
		t.Errorf("Expected parent title 'Parent 1', got '%s'", result[2].ParentTitle)
	}
	if result[3].ParentTitle != "Parent 1" {
		t.Errorf("Expected parent title 'Parent 1', got '%s'", result[3].ParentTitle)
	}
}

// TestBuildFlatHierarchy_OrphanedSubtask tests mixed orphans and valid subtasks
func TestBuildFlatHierarchy_OrphanedSubtask(t *testing.T) {
	tasks := []*models.Task{
		makeTask("P1", "Parent 1", 1, nil),
		makeTask("S1", "Valid subtask", 2, ptr("P1")),
		makeTask("S2", "Orphaned subtask", 3, ptr("MISSING")),
	}

	result := BuildFlatHierarchy(tasks, SortByPriority)

	if len(result) != 3 {
		t.Fatalf("Expected 3 items, got %d", len(result))
	}

	// Expected order: P1, S1 (under P1), S2 (orphaned, at end)
	expected := []string{"P1", "S1", "S2"}
	for i, exp := range expected {
		if result[i].Task.ID != exp {
			t.Errorf("Position %d: expected %s, got %s", i, exp, result[i].Task.ID)
		}
	}

	// S2 should be treated as top-level (not subtask)
	if result[2].IsSubtask {
		t.Error("Orphaned subtask should be treated as top-level")
	}
	if result[2].Indentation != 0 {
		t.Error("Orphaned subtask should have zero indentation")
	}
}

// TestBuildFlatHierarchy_SortByDescription tests title sorting
func TestBuildFlatHierarchy_SortByDescription(t *testing.T) {
	tasks := []*models.Task{
		makeTask("P1", "Zebra", 1, nil),
		makeTask("P2", "Apple", 2, nil),
	}

	result := BuildFlatHierarchy(tasks, SortByDescription)

	// Should be sorted alphabetically
	if result[0].Task.Title != "Apple" {
		t.Errorf("Expected Apple first, got %s", result[0].Task.Title)
	}
	if result[1].Task.Title != "Zebra" {
		t.Errorf("Expected Zebra second, got %s", result[1].Task.Title)
	}
}

// TestBuildFlatHierarchy_SortByCreated tests date sorting
func TestBuildFlatHierarchy_SortByCreated(t *testing.T) {
	older := time.Now().Add(-1 * time.Hour)
	newer := time.Now()

	tasks := []*models.Task{
		{ID: "P1", Title: "Newer", Priority: 1, ParentID: nil, CreatedAt: newer, UpdatedAt: time.Now()},
		{ID: "P2", Title: "Older", Priority: 2, ParentID: nil, CreatedAt: older, UpdatedAt: time.Now()},
	}

	result := BuildFlatHierarchy(tasks, SortByCreated)

	// Should be sorted oldest first
	if result[0].Task.ID != "P2" {
		t.Errorf("Expected P2 (older) first, got %s", result[0].Task.ID)
	}
}

// TestGetSubtasksForParent_NoSubtasks tests empty result
func TestGetSubtasksForParent_NoSubtasks(t *testing.T) {
	tasks := []*models.Task{
		makeTask("P1", "Parent", 1, nil),
		makeTask("P2", "Another parent", 2, nil),
	}

	subtasks := GetSubtasksForParent(tasks, "P1")

	if len(subtasks) != 0 {
		t.Errorf("Expected no subtasks, got %d", len(subtasks))
	}
}

// TestGetSubtasksForParent_MultipleSubtasks tests finding all subtasks
func TestGetSubtasksForParent_MultipleSubtasks(t *testing.T) {
	tasks := []*models.Task{
		makeTask("P1", "Parent", 1, nil),
		makeTask("S1", "Subtask 1", 2, ptr("P1")),
		makeTask("S2", "Subtask 2", 3, ptr("P1")),
		makeTask("P2", "Other parent", 4, nil),
		makeTask("S3", "Other subtask", 5, ptr("P2")),
	}

	subtasks := GetSubtasksForParent(tasks, "P1")

	if len(subtasks) != 2 {
		t.Fatalf("Expected 2 subtasks for P1, got %d", len(subtasks))
	}

	// Check we got the right subtasks
	found := make(map[string]bool)
	for _, st := range subtasks {
		found[st.ID] = true
	}

	if !found["S1"] || !found["S2"] {
		t.Error("Expected to find S1 and S2 as subtasks of P1")
	}
	if found["S3"] {
		t.Error("Should not find S3 (belongs to P2)")
	}
}

// TestFormatWithIndentation_Parent tests parent formatting
func TestFormatWithIndentation_Parent(t *testing.T) {
	ht := &HierarchicalTask{
		Task:        makeTask("P1", "Parent task", 1, nil),
		IsSubtask:   false,
		Indentation: 0,
	}

	formatted := FormatWithIndentation(ht)
	expected := "Parent task"

	if formatted != expected {
		t.Errorf("Expected '%s', got '%s'", expected, formatted)
	}
}

// TestFormatWithIndentation_Subtask tests subtask formatting
func TestFormatWithIndentation_Subtask(t *testing.T) {
	ht := &HierarchicalTask{
		Task:        makeTask("S1", "Subtask", 2, ptr("P1")),
		IsSubtask:   true,
		Indentation: 2,
		ParentTitle: "Parent task",
	}

	formatted := FormatWithIndentation(ht)
	expected := "  - Subtask" // 2 spaces + dash + space + title

	if formatted != expected {
		t.Errorf("Expected '%s', got '%s'", expected, formatted)
	}
}

// TestIsOrphanedSubtask_ValidParent tests valid parent case
func TestIsOrphanedSubtask_ValidParent(t *testing.T) {
	tasks := []*models.Task{
		makeTask("P1", "Parent", 1, nil),
		makeTask("S1", "Subtask", 2, ptr("P1")),
	}

	isOrphan := IsOrphanedSubtask(tasks[1], tasks)

	if isOrphan {
		t.Error("Expected subtask to have valid parent")
	}
}

// TestIsOrphanedSubtask_MissingParent tests orphan case
func TestIsOrphanedSubtask_MissingParent(t *testing.T) {
	tasks := []*models.Task{
		makeTask("P1", "Parent", 1, nil),
		makeTask("S1", "Orphaned subtask", 2, ptr("MISSING")),
	}

	isOrphan := IsOrphanedSubtask(tasks[1], tasks)

	if !isOrphan {
		t.Error("Expected subtask to be orphaned (parent not in list)")
	}
}

// TestIsOrphanedSubtask_TopLevelTask tests parent task returns false
func TestIsOrphanedSubtask_TopLevelTask(t *testing.T) {
	tasks := []*models.Task{
		makeTask("P1", "Parent", 1, nil),
	}

	isOrphan := IsOrphanedSubtask(tasks[0], tasks)

	if isOrphan {
		t.Error("Expected top-level task to not be orphaned")
	}
}
