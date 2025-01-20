package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"todoapp/server"
	"todoapp/store"
)

func main() {
	killChan := make(chan os.Signal, 1)
	signal.Notify(killChan, os.Interrupt, syscall.SIGTERM)

	c := store.Config{LoadFromFile: true, DBName: "todo_app"}

	//s := store.NewInMemoryStore(c)
	s, _ := store.NewPostgresStore(c)

	go server.Start(s)

	<-killChan

	//s.SaveTasksToFile()
	fmt.Println("Server shut down ...")
}
