package store

import "github.com/google/uuid"

type Store interface {
	GetAllItems() ([]Task, error)
	AddItem(id uuid.UUID, title string, priority Priority) error
	DeleteItem(id uuid.UUID) error
	ToggleDone(id uuid.UUID) error
	EditTask(id uuid.UUID, title string) error
}
