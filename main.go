package main

import (
	"log"
	"net/http"
	"todoapp/server"
	"todoapp/store"
)

func main() {
	store := store.NewInMemoryTaskStore()
	taskServer := server.NewTaskServer(store)
	log.Fatal(http.ListenAndServe(":5000", taskServer))
}
