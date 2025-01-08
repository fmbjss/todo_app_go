package store

import (
	"testing"

	"github.com/google/uuid"
)

func TestStore(t *testing.T) {
	t.Run("add task", func(t *testing.T) {
		memStore := &Store{}

		taskID := addTaskToStore(memStore, "Mow the lawn", Medium)

		items := GetAllItems(memStore)
		expectedTask := "Mow the lawn"
		expectedPriority := Medium

		if len(items) != 1 {
			t.Fatalf("expected 1 task, got %d", len(items))
		}
		if items[0].Title != expectedTask {
			t.Errorf("expected task title '%s', got '%s'", expectedTask, items[0].Title)
		}
		if items[0].Priority != expectedPriority {
			t.Errorf("expected priority '%s', got '%s'", expectedPriority, items[0].Priority)
		}
		if items[0].ID != taskID {
			t.Errorf("expected task ID '%v', got '%v'", taskID, items[0].ID)
		}
	})

	t.Run("delete task", func(t *testing.T) {
		memStore := &Store{}

		taskID := addTaskToStore(memStore, "Test Task to Delete", Low)

		err := DeleteItem(memStore, taskID)
		if err != nil {
			t.Fatalf("expected no error but got: %v", err)
		}

		items := GetAllItems(memStore)
		if len(items) != 0 {
			t.Errorf("expected 0 tasks, got %d", len(items))
		}
	})

	t.Run("edit task", func(t *testing.T) {
		memStore := &Store{}

		taskID := addTaskToStore(memStore, "Test Task", Low)

		err := EditTask(memStore, taskID, "Updated Test Task")
		if err != nil {
			t.Fatalf("expected no error but got: %v", err)
		}

		items := GetAllItems(memStore)
		if len(items) != 1 {
			t.Fatalf("expected 1 task, got %d", len(items))
		}
		if items[0].Title != "Updated Test Task" {
			t.Errorf("expected 'Updated Test Task', got '%s'", items[0].Title)
		}
	})

	t.Run("complete task", func(t *testing.T) {
		memStore := &Store{}

		taskID := addTaskToStore(memStore, "Test Task", Low)

		err := ToggleDone(memStore, taskID)
		if err != nil {
			t.Fatalf("expected no error but got: %v", err)
		}

		items := GetAllItems(memStore)
		if len(items) != 1 {
			t.Fatalf("expected 1 task, got %d", len(items))
		}
		if !items[0].Done {
			t.Errorf("expected task to be completed, but it is not")
		}

		// Toggle it back to false
		err = ToggleDone(memStore, taskID)
		if err != nil {
			t.Fatalf("expected no error but got: %v", err)
		}
		if items[0].Done {
			t.Errorf("expected task to be incomplete, but it is still completed")
		}
	})
}

func addTaskToStore(memStore *Store, title string, priority Priority) uuid.UUID {
	taskID := uuid.New()
	AddItem(memStore, taskID, title, priority)
	return taskID
}
