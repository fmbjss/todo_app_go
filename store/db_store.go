package store

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5431
	user     = "postgres"
	password = ""
	dbname   = "todo_app"
)

type PostgresStore struct {
	db          *sql.DB
	taskChannel chan TaskOperation
	stopChannel chan struct{}
}

func NewPostgresStore(config Config) (*PostgresStore, error) {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		log.Println("Error connecting to database")
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		log.Println("Error pinging database")
	}

	store := &PostgresStore{
		db:          db,
		taskChannel: make(chan TaskOperation),
		stopChannel: make(chan struct{}),
	}

	if config.LoadFromFile {

		if err := store.initSchema(); err != nil {
			log.Println("Error initialising db")
			return nil, err
		}
	}
	go func() {
		err := store.processTasks()
		if err != nil {
			log.Fatal(err)
		}
	}()
	return store, nil
}

func (s *PostgresStore) GetAllItems() ([]Task, error) {
	if s.db == nil {
		log.Println("s.db == nil")
		return nil, fmt.Errorf("database connection is not initialized")
	}
	rows, err := s.db.Query("SELECT * FROM tasks")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(rows)

	var tasks []Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Priority, &task.Done); err != nil {
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tasks, err
}

func (s *PostgresStore) processTasks() error {
	for {
		select {
		case op, ok := <-s.taskChannel:
			if !ok {
				fmt.Println("Task channel closed")
				return nil
			}

			switch op.Type {

			case "Add":
				_, err := s.db.Exec("INSERT INTO tasks (id, title, priority, done) VALUES ($1, $2, $3, $4)", op.ID, op.Title, op.Priority, false)
				if err != nil {
					log.Printf("Error adding task: %v", err)
				}
				op.Result <- err
				close(op.Result)

			case "Delete":
				_, err := s.db.Exec("DELETE FROM tasks WHERE id = $1", op.ID)
				if err != nil {
					log.Printf("Error deleting task: %v", err)
				}
				op.Result <- err
				close(op.Result)

			case "Edit":
				_, err := s.db.Exec("UPDATE tasks SET title = $1 WHERE id = $2", op.Title, op.ID)
				if err != nil {
					log.Printf("Error editing task: %v", err)
				}
				op.Result <- err
				close(op.Result)

			case "ToggleDone":
				_, err := s.db.Exec("UPDATE tasks SET done = NOT done WHERE id = $1", op.ID)
				if err != nil {
					log.Printf("Error toggling task done status: %v", err)
				}
				op.Result <- err
				close(op.Result)

			}

		case <-s.stopChannel:
			close(s.taskChannel)
			err := s.db.Close()
			if err != nil {

			}
			return nil
		}
	}
}

func (s *PostgresStore) AddItem(id uuid.UUID, t string, p Priority) error {
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

func (s *PostgresStore) DeleteItem(id uuid.UUID) error {
	result := make(chan error)
	s.taskChannel <- TaskOperation{
		Type:   "Delete",
		ID:     id,
		Result: result,
	}
	return <-result
}

func (s *PostgresStore) EditTask(id uuid.UUID, t string) error {
	result := make(chan error)
	s.taskChannel <- TaskOperation{
		Type:   "Edit",
		ID:     id,
		Title:  t,
		Result: result,
	}
	return <-result
}

func (s *PostgresStore) ToggleDone(id uuid.UUID) error {
	result := make(chan error)
	s.taskChannel <- TaskOperation{
		Type:   "ToggleDone",
		ID:     id,
		Result: result,
	}
	return <-result
}

func (s *PostgresStore) initSchema() error {

	query := `
	CREATE TABLE IF NOT EXISTS tasks (
		id UUID PRIMARY KEY,
		title TEXT NOT NULL,
		priority TEXT NOT NULL CHECK (priority IN ('Low', 'Medium', 'High')),
		done BOOLEAN NOT NULL
	)`
	_, err := s.db.Exec(query)
	return err
}
