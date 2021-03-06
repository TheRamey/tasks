package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

// Todo is the structure of the tasks
type Todo struct {
	Task       string `json:"task"`
	IsComplete bool   `json:"isComplete"`
}

func (t Todo) toString() string {
	bytes, err := json.Marshal(t)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return string(bytes)
}

func updateTodo(n int, s string, b bool) error {
	todos := getTodos()
	var updatedTodos []Todo
	if n > len(todos) {
		return errors.New("Todo not found, nothing updated")
	}
	for i, todo := range todos {
		num := i + 1
		if num == n {
			single := &todo
			single.Task = s
			single.IsComplete = b
		}
		updatedTodos = append(updatedTodos, todo)
	}

	bytes, err := json.Marshal(updatedTodos)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	path, err := getPath()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	file := filepath.Join(path, ".todos.json")
	err = ioutil.WriteFile(file, bytes, 0755)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return nil
}

func deleteTodo(n int) error {
	todos := getTodos()
	var updatedTodos []Todo
	for i, todo := range todos {
		num := i + 1
		if num != n {
			updatedTodos = append(updatedTodos, todo)
		}
	}
	if len(todos) == len(updatedTodos) {
		return errors.New("Todo not found, nothing deleted")
	}

	bytes, err := json.Marshal(updatedTodos)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	path, err := getPath()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	file := filepath.Join(path, ".todos.json")
	err = ioutil.WriteFile(file, bytes, 0755)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return nil
}

func addTodo(t string) error {
	todos := getTodos()
	todo := Todo{Task: t, IsComplete: false}
	newTodos := append(todos, todo)
	if len(todos) == len(newTodos) {
		return errors.New("Failed to add task")
	}
	bytes, err := json.Marshal(newTodos)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	path, err := getPath()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	file := filepath.Join(path, ".todos.json")
	err = ioutil.WriteFile(file, bytes, 0755)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return nil
}

func getTodos() []Todo {
	var todos []Todo
	path, err := getPath()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	file := filepath.Join(path, ".todos.json")
	raw, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	err = json.Unmarshal(raw, &todos)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return todos
}

func getPath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", errors.New("Cannot get home directory")
	}

	path := usr.HomeDir
	return path, nil
}

func listCommands() {
	green := color.New(color.FgGreen).SprintFunc()
	magenta := color.New(color.FgMagenta).SprintFunc()
	fmt.Printf("\n\n%s\n\n%s\n%s\n%s\n%s\n%s\n%s\n\n\n",
		magenta("List of Commands:"),
		green("list - will list all tasks"),
		green("add - create a new task"),
		green("delete <num> - delete a task by ID"),
		green("edit <num> <new description> - edit the text of a task by ID"),
		green("complete <num> - mark a task complete by ID"),
		green("incomplete <num> - mark a task incomplete by ID"))
}

func listTodos() {
	blue := color.New(color.FgBlue).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	todos := getTodos()

	for i, todo := range todos {
		num := i + 1
		if todo.IsComplete == true {
			fmt.Printf("%d. %s %s\n", num, blue("[X]"), red(todo.Task))
		} else {
			fmt.Printf("%d. %s %s\n", num, blue("[ ]"), red(todo.Task))
		}
	}
}

func main() {

	path, err := getPath()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	file := filepath.Join(path, ".todos.json")
	if _, err := os.Stat(file); os.IsNotExist(err) {
		err := ioutil.WriteFile(file, nil, 0755)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	args := os.Args[1:]
	todos := getTodos()

	switch {
	case len(args) < 1:
		listCommands()
	case args[0] == "list":
		listTodos()
	case args[0] == "add":
		task := strings.Join(args[1:], " ")
		err := addTodo(task)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		listTodos()
	case args[0] == "delete":
		num, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		err = deleteTodo(num)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		listTodos()
	case args[0] == "complete":
		num, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		err = updateTodo(num, todos[num-1].Task, true)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		listTodos()
	case args[0] == "incomplete":
		num, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		err = updateTodo(num, todos[num-1].Task, false)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		listTodos()
	case args[0] == "edit":
		num, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		newTask := strings.Join(args[2:], " ")
		err = updateTodo(num, newTask, todos[num-1].IsComplete)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		listTodos()
	default:
		listCommands()
	}

}
