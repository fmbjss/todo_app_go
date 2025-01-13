package store

import (
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"
)

type Priority string

const (
	Low    Priority = "Low"
	Medium Priority = "Medium"
	High   Priority = "High"
)

type Task struct {
	ID       uuid.UUID
	Title    string
	Priority Priority
	Done     bool
}

type TaskOperation struct {
	Type     string
	ID       uuid.UUID
	Title    string
	Priority Priority
	Result   chan error
}

type InMemoryStore struct {
	tasks       []Task
	mu          sync.Mutex
	taskChannel chan TaskOperation
	stopChannel chan struct{}
}

func NewInMemoryStore() *InMemoryStore {
	store := &InMemoryStore{
		tasks:       []Task{},
		taskChannel: make(chan TaskOperation),
		stopChannel: make(chan struct{}),
	}
	go store.processTasks()
	return store
}

func (s *InMemoryStore) GetAllItems() []Task {
	return s.tasks
}

func (s *InMemoryStore) processTasks() {
	for {
		select {
		case op, ok := <-s.taskChannel:
			if !ok {
				fmt.Println("Task channel closed")
				return
			}
			var err error
			switch op.Type {

			case "Add":
				task := Task{
					ID:       op.ID,
					Title:    op.Title,
					Priority: op.Priority,
					Done:     false,
				}
				s.tasks = append(s.tasks, task)
			case "Delete":
				found := false
				for i, task := range s.tasks {
					if task.ID == op.ID {
						s.tasks = append(s.tasks[:i], s.tasks[i+1:]...)
						found = true
						break
					}
				}
				if !found {
					err = errors.New("task not found")
				}
			case "Edit":
				found := false
				for i, task := range s.tasks {
					if task.ID == op.ID {
						s.tasks[i].Title = op.Title
						found = true
						break
					}
				}
				if !found {
					err = errors.New("task not found")
				}
			case "ToggleDone":
				found := false
				for i, task := range s.tasks {
					if task.ID == op.ID {
						s.tasks[i].Done = true
						found = true
						break
					}
				}
				if !found {
					err = errors.New("task not found")
				}
			}
			if op.Result != nil {
				op.Result <- err
				close(op.Result)
			}

		case <-s.stopChannel:
			fmt.Println("case <- s.stopChannel")
			close(s.taskChannel)
			return
		}
	}
}

func (s *InMemoryStore) AddItem(id uuid.UUID, t string, p Priority) error {
	result := make(chan error)
	s.taskChannel <- TaskOperation{
		Type:     "Add",
		ID:       id,
		Title:    t,
		Priority: p,
		Result:   result,
	}
	return <-result
}

func (s *InMemoryStore) DeleteItem(id uuid.UUID) error {
	result := make(chan error)
	s.taskChannel <- TaskOperation{
		Type:   "Delete",
		ID:     id,
		Result: result,
	}
	return <-result
}

func (s *InMemoryStore) EditTask(id uuid.UUID, t string) error {
	result := make(chan error)
	s.taskChannel <- TaskOperation{
		Type:   "Edit",
		ID:     id,
		Title:  t,
		Result: result,
	}
	return <-result
}

func (s *InMemoryStore) ToggleDone(id uuid.UUID) error {
	result := make(chan error)
	s.taskChannel <- TaskOperation{
		Type:   "ToggleDone",
		ID:     id,
		Result: result,
	}
	return <-result
}
