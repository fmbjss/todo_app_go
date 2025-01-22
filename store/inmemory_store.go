package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/google/uuid"
)

type InMemoryStore struct {
	tasks       []Task
	mu          sync.Mutex
	taskChannel chan TaskOperation
	stopChannel chan struct{}
	filePath    string
}

func NewInMemoryStore(config Config) (*InMemoryStore, error) {
	store := &InMemoryStore{
		tasks:       []Task{},
		taskChannel: make(chan TaskOperation),
		stopChannel: make(chan struct{}),
		filePath:    "tasks.json",
	}
	if config.LoadFromFile {
		store.loadTasksFromFile()
	}
	go store.processTasks()
	return store, nil
}

func (s *InMemoryStore) GetAllItems() ([]Task, error) {
	return s.tasks, nil
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
						s.tasks[i].Done = !s.tasks[i].Done
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

type TaskFile struct {
	Tasks []Task `json:"tasks"`
}

func (s *InMemoryStore) loadTasksFromFile() {
	s.mu.Lock()
	defer s.mu.Unlock()

	file, err := os.Open(s.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			s.tasks = []Task{}
			return
		}
		panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	var taskFile TaskFile
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&taskFile); err != nil {
		panic(err)
	}

	s.tasks = taskFile.Tasks
}

func (s *InMemoryStore) SaveTasksToFile() {
	s.mu.Lock()
	defer s.mu.Unlock()

	file, err := os.Create(s.filePath)
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	taskFile := TaskFile{
		Tasks: s.tasks,
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(taskFile); err != nil {
		panic(err)
	}
}
