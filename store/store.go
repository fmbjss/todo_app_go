package store

import (
	"errors"

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

type Store struct {
	tasks []Task
}

func AddItem(s *Store, id uuid.UUID, t string, p Priority) {
	task := Task{
		ID:       id,
		Title:    t,
		Priority: p,
		Done:     false,
	}
	s.tasks = append(s.tasks, task)
}

func DeleteItem(s *Store, id uuid.UUID) error {
	for i, task := range s.tasks {
		if task.ID == id {
			s.tasks = append(s.tasks[:i], s.tasks[i+1:]...)
			return nil
		}
	}
	return errors.New("task not found")
}

func EditTask(s *Store, id uuid.UUID, t string) error {
	for i, task := range s.tasks {
		if task.ID == id {
			s.tasks[i].Title = t
			return nil
		}
	}
	return errors.New("task not found")
}

func ToggleDone(s *Store, id uuid.UUID) error {
	for i, task := range s.tasks {
		if task.ID == id {
			s.tasks[i].Done = !s.tasks[i].Done
			return nil
		}
	}
	return errors.New("task not found")
}

func GetAllItems(memStore *Store) []Task {
	return memStore.tasks
}
