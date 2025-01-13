package main

import (
	"todoapp/cli"
	"todoapp/store"
)

func main() {
	s := store.NewInMemoryStore()

	cli.Start(s)

}
