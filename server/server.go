package server

import (
	"encoding/json"
	"net/http"
	"todoapp/store"

	"github.com/google/uuid"
)

type TaskServer struct {
	store *store.InMemoryTaskStore
}

func NewTaskServer(store *store.InMemoryTaskStore) *TaskServer {
	return &TaskServer{store: store}
}

func (s *TaskServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.getTasks(w)
	case http.MethodPost:
		s.addTask(w, r)
	case http.MethodPut:
		s.editTask(w, r)
	case http.MethodDelete:
		s.deleteTask(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *TaskServer) getTasks(w http.ResponseWriter) {
	tasks := s.store.GetAllItems()
	json.NewEncoder(w).Encode(tasks)
}

func (s *TaskServer) addTask(w http.ResponseWriter, r *http.Request) {
	var task struct {
		Title    string         `json:"title"`
		Priority store.Priority `json:"priority"`
	}
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	s.store.AddItem(uuid.New(), task.Title, task.Priority)
	w.WriteHeader(http.StatusCreated)
}

func (s *TaskServer) editTask(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	var task struct {
		Title string `json:"title"`
	}
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	if err := s.store.EditTask(id, task.Title); err != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *TaskServer) deleteTask(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	if err := s.store.DeleteItem(id); err != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}
