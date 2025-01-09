package store

import (
	"errors"
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

type InMemoryTaskStore struct {
	tasks []Task
	mu    sync.Mutex
}

func NewInMemoryTaskStore() *InMemoryTaskStore {
	return &InMemoryTaskStore{
		tasks: []Task{},
		mu:    sync.Mutex{},
	}
}

func (s *InMemoryTaskStore) AddItem(id uuid.UUID, t string, p Priority) {
	s.mu.Lock()
	defer s.mu.Unlock()
	task := Task{
		ID:       id,
		Title:    t,
		Priority: p,
		Done:     false,
	}
	s.tasks = append(s.tasks, task)
}

func (s *InMemoryTaskStore) DeleteItem(id uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, task := range s.tasks {
		if task.ID == id {
			s.tasks = append(s.tasks[:i], s.tasks[i+1:]...)
			return nil
		}
	}
	return errors.New("task not found")
}

func (s *InMemoryTaskStore) EditTask(id uuid.UUID, t string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, task := range s.tasks {
		if task.ID == id {
			s.tasks[i].Title = t
			return nil
		}
	}
	return errors.New("task not found")
}

func (s *InMemoryTaskStore) ToggleDone(id uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, task := range s.tasks {
		if task.ID == id {
			s.tasks[i].Done = !s.tasks[i].Done
			return nil
		}
	}
	return errors.New("task not found")
}

func (s *InMemoryTaskStore) GetAllItems() []Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.tasks
}
