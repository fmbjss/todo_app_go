package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"todoapp/store"

	"github.com/google/uuid"
)

func Start() {
	memStore := &store.Store{}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Task Manager CLI")
	fmt.Println("Commands: add title priority, delete task_id, edit task_id new_title, toggle task_id, list, quit")

	for {
		fmt.Print("> ")
		scanner.Scan()
		command := scanner.Text()
		args := strings.Fields(command)

		if len(args) == 0 {
			continue
		}

		switch args[0] {
		case "add":
			if len(args) < 3 {
				fmt.Println("Usage: add title priority")
				continue
			}
			title := args[1]
			priority := args[2]
			p, valid := mapStringToPriorityType(priority)
			if !valid {
				fmt.Println("Invalid priority. Valid values are: low, medium, high")
				continue
			}
			id := uuid.New()
			store.AddItem(memStore, id, title, p)
			fmt.Printf("Task added with ID: %s\n", id)

		case "delete":
			if len(args) < 2 {
				fmt.Println("Usage: delete task_id")
				continue
			}
			id, err := uuid.Parse(args[1])
			if err != nil {
				fmt.Println("Invalid UUID format")
				continue
			}
			err = store.DeleteItem(memStore, id)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("Task deleted")
			}

		case "edit":
			if len(args) < 3 {
				fmt.Println("Usage: edit task_id new_title")
				continue
			}
			id, err := uuid.Parse(args[1])
			if err != nil {
				fmt.Println("Invalid UUID format")
				continue
			}
			newTitle := args[2]
			err = store.EditTask(memStore, id, newTitle)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("Task edited")
			}

		case "toggle":
			if len(args) < 2 {
				fmt.Println("Usage: toggle task_id")
				continue
			}
			id, err := uuid.Parse(args[1])
			if err != nil {
				fmt.Println("Invalid UUID format")
				continue
			}
			err = store.ToggleDone(memStore, id)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("Task completion toggled")
			}

		case "list":
			tasks := store.GetAllItems(memStore)
			if len(tasks) == 0 {
				fmt.Println("No tasks available.")
			} else {
				fmt.Println("Tasks:")
				for _, task := range tasks {
					status := "Incomplete"
					if task.Done {
						status = "Complete"
					}
					fmt.Printf("ID: %s, Title: %s, Priority: %s, Status: %s\n", task.ID, task.Title, task.Priority, status)
				}
			}

		case "quit":
			fmt.Println("Exiting Task Manager CLI.")
			return

		default:
			fmt.Println("Unknown command!")
		}
	}
}

func mapStringToPriorityType(priority string) (store.Priority, bool) {
	switch strings.ToLower(priority) {
	case "low":
		return store.Low, true
	case "medium":
		return store.Medium, true
	case "high":
		return store.High, true
	default:
		return "", false
	}
}
