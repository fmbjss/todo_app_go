package store

import (
	"fmt"
	"sync"
	"testing"

	"github.com/google/uuid"
)

func TestStore(t *testing.T) {
	c := Config{LoadFromFile: false}

	t.Run("add task", func(t *testing.T) {
		store := NewInMemoryStore(c)
		taskID := uuid.New()
		taskTitle := "Test Task"
		taskPriority := High

		err := store.AddItem(taskID, taskTitle, taskPriority)
		if err != nil {
			t.Errorf("Error adding task: %s", err)
		}

		tasks := store.GetAllItems()

		if len(tasks) != 1 {
			t.Errorf("expected 1 task, got %d", len(tasks))
		}
		fmt.Println(tasks[0].ID == taskID)

		if tasks[0].ID != taskID {
			t.Errorf("expected task ID %s, got %s", taskID, tasks[0].ID)
		}
		if tasks[0].Title != taskTitle {
			t.Errorf("expected task title '%s', got '%s'", taskTitle, tasks[0].Title)
		}
		if tasks[0].Priority != taskPriority {
			t.Errorf("expected task priority '%s', got '%s'", taskPriority, tasks[0].Priority)
		}

	})

	t.Run("edit task", func(t *testing.T) {
		store := NewInMemoryStore(c)
		taskID := uuid.New()

		err := store.AddItem(taskID, "Test Task", Low)
		if err != nil {
			t.Errorf("Error adding task: %s", err)
		}
		updatedTitle := "Updated Task"
		if err := store.EditTask(taskID, updatedTitle); err != nil {
			t.Fatalf("expected no error, got %s", err)
		}

		tasks := store.GetAllItems()

		if len(tasks) != 1 {
			t.Errorf("expected 1 task, got %d", len(tasks))
		}
		if tasks[0].Title != updatedTitle {
			t.Errorf("expected task title '%s', got '%s'", updatedTitle, tasks[0].Title)
		}
	})
	t.Run("delete task", func(t *testing.T) {
		store := NewInMemoryStore(c)
		taskID := uuid.New()
		err := store.AddItem(taskID, "Test Task", Low)
		if err != nil {
			t.Errorf("Error adding task: %s", err)
		}

		errDelete := store.DeleteItem(taskID)
		if errDelete != nil {
			t.Errorf("Error deleting task: %s", err)
		}

		tasks := store.GetAllItems()
		if len(tasks) != 0 {
			t.Errorf("expected 0 task, got %d", len(tasks))
		}

	})
	t.Run("toggle tasks", func(t *testing.T) {
		store := NewInMemoryStore(c)
		taskID := uuid.New()
		err := store.AddItem(taskID, "Test Task", Low)
		if err != nil {
			t.Errorf("Error adding task: %s", err)
		}

		errToggle := store.ToggleDone(taskID)
		if errToggle != nil {
			t.Errorf("Error toggling task: %s", err)
		}

		tasks := store.GetAllItems()
		if len(tasks) != 1 {
			t.Errorf("expected 1 task, got %d", len(tasks))
		}
		if tasks[0].Done != true {
			t.Errorf("expected task done true, got %t", tasks[0].Done)
		}

	})
}

func BenchmarkStoreOperations(b *testing.B) {
	c := Config{LoadFromFile: false}

	b.Run("AddItem", func(b *testing.B) {
		store := NewInMemoryStore(c)
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				taskID := uuid.New()
				taskTitle := "Benchmark Task"
				taskPriority := High
				err := store.AddItem(taskID, taskTitle, taskPriority)
				if err != nil {
					b.Errorf("Error adding task: %s", err)
				}
			}
		})
	})

	b.Run("EditItem", func(b *testing.B) {
		store := NewInMemoryStore(c)
		taskIDs := make([]uuid.UUID, b.N)
		for i := 0; i < b.N; i++ {
			taskID := uuid.New()
			taskIDs[i] = taskID
			err := store.AddItem(taskID, "Benchmark Task", Medium)
			if err != nil {
				return
			}
		}

		b.ResetTimer()

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				for _, taskID := range taskIDs {
					err := store.EditTask(taskID, "Edited Task")
					if err != nil {
						b.Errorf("Error editing task: %s", err)
					}
				}
			}
		})
	})

	b.Run("ToggleDone", func(b *testing.B) {
		store := NewInMemoryStore(c)
		taskIDs := make([]uuid.UUID, b.N)
		for i := 0; i < b.N; i++ {
			taskID := uuid.New()
			taskIDs[i] = taskID
			err := store.AddItem(taskID, "Benchmark Task", Low)
			if err != nil {
				return
			}
		}

		b.ResetTimer()

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				for _, taskID := range taskIDs {
					err := store.ToggleDone(taskID)
					if err != nil {
						b.Errorf("Error toggling task: %s", err)
					}
				}
			}
		})
	})

	b.Run("DeleteItem", func(b *testing.B) {
		store := NewInMemoryStore(c)
		taskIDs := make([]uuid.UUID, b.N)
		for i := 0; i < b.N; i++ {
			taskID := uuid.New()
			taskIDs[i] = taskID
			err := store.AddItem(taskID, "Benchmark Task", High)
			if err != nil {
				return
			}
		}

		b.ResetTimer()

		var mu sync.Mutex

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				mu.Lock()
				if len(taskIDs) > 0 {
					taskID := taskIDs[0]
					taskIDs = taskIDs[1:]
					mu.Unlock()
					err := store.DeleteItem(taskID)
					if err != nil {
						b.Errorf("Error deleting task: %s", err)
					}
				} else {
					mu.Unlock()
				}
			}
		})
	})
}
