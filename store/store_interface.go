package store

import (
	"github.com/google/uuid"
)

type Store interface {
	GetAllItems() ([]Task, error)
	AddItem(id uuid.UUID, title string, priority Priority) error
	DeleteItem(id uuid.UUID) error
	ToggleDone(id uuid.UUID) error
	EditTask(id uuid.UUID, title string) error
}

type Priority string

const (
	Low    Priority = "Low"
	Medium Priority = "Medium"
	High   Priority = "High"
)

type Config struct {
	LoadFromFile bool
	DBName       string
}
type Task struct {
	ID       uuid.UUID `json:"ID"`
	Title    string    `json:"Title"`
	Priority Priority  `json:"Priority"`
	Done     bool      `json:"Done"`
}

type TaskOperation struct {
	Type     string
	ID       uuid.UUID
	Title    string
	Priority Priority
	Result   chan error
}
