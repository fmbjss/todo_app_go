package server

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"todoapp/store"

	"github.com/google/uuid"
)

type TaskServer struct {
	store store.Store
}

func NewTaskServer(store store.Store) *TaskServer {
	return &TaskServer{store: store}
}

func LoadTemplate() (*template.Template, error) {
	tmplPath := filepath.Join("server", "todo_app.html")

	return template.ParseFiles(tmplPath)
}

func (s *TaskServer) renderTasksPage(w http.ResponseWriter) {

	tasks, err := s.store.GetAllItems()
	if err != nil {
		log.Println("Error loading tasks")
		return
	}

	tmpl, err := LoadTemplate()
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, tasks)
	if err != nil {
		http.Error(w, "Error rendering tasks", http.StatusInternalServerError)
	}
}

func ParseID(r *http.Request) (uuid.UUID, error) {
	idStr := r.FormValue("ID")
	taskID, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.UUID{}, err
	}
	return taskID, nil
}

func (s *TaskServer) home(w http.ResponseWriter, _ *http.Request) {
	s.renderTasksPage(w)
}

func (s *TaskServer) addTask(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	task := store.Task{
		ID:       uuid.New(),
		Title:    r.FormValue("title"),
		Priority: store.Priority(r.FormValue("priority")),
	}

	if err := s.store.AddItem(task.ID, task.Title, task.Priority); err != nil {
		http.Error(w, "Error adding task", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *TaskServer) deleteTask(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	taskID, err := ParseID(r)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	if err := s.store.DeleteItem(taskID); err != nil {
		http.Error(w, "Error deleting task", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *TaskServer) toggleDone(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	taskID, err := ParseID(r)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	if err := s.store.ToggleDone(taskID); err != nil {
		http.Error(w, "Error toggling task", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *TaskServer) edit(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	taskID, err := ParseID(r)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	taskTitle := r.FormValue("title")
	if err := s.store.EditTask(taskID, taskTitle); err != nil {
		http.Error(w, "Error editing task", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func Start(store store.Store) {
	log.Println("Web API server is running on http://localhost:8080")

	taskServer := NewTaskServer(store)
	http.HandleFunc("/", taskServer.home)
	http.HandleFunc("/add", taskServer.addTask)
	http.HandleFunc("/delete", taskServer.deleteTask)
	http.HandleFunc("/toggle", taskServer.toggleDone)
	http.HandleFunc("/edit", taskServer.edit)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println(err)
	}
}
