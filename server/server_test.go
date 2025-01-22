package server

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"todoapp/store"
)

func BenchmarkServer(b *testing.B) {
	c := store.Config{LoadFromFile: false, DBName: "benchmark_tests"}
	s, _ := store.NewPostgresStore(c)
	defer func() {
		_, err := s.Db.Exec("TRUNCATE TABLE tasks RESTART IDENTITY CASCADE")
		if err != nil {
			return
		}
	}()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		NewTaskServer(s).addTask(w, r)
	}))

	defer ts.Close()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := http.PostForm(ts.URL+"/add", url.Values{
				"title":    {"benchmark"},
				"priority": {"Low"},
			})
			if err != nil {
				b.Errorf("Failed request: %v", err)
			}
		}
	})
}
