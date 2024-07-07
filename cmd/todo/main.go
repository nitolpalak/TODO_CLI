package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"todo"
)

const (
	todofile = ".todo.json"
)

func main() {
	add := flag.Bool("add", false, "add a todo")
	complete := flag.Int("complete", 0, "mark a todo as completed")
	delete := flag.Int("delete", 0, "delete a todo")
	list := flag.Bool("list", false, "show the list of all todos")

	flag.Parse()

	handleError := func(err error, msg string) {
		if err != nil {
			fmt.Println(msg, err)
			os.Exit(1)
		}
	}

	todos := &todo.Todos{}

	if err := todos.Load(todofile); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	switch {
	case *add:
		msg := "The name of the task: "
		task, err := getInput(msg, os.Stdin, flag.Args()...)
		handleError(err, "Error getting input:")

		todos.Add(task)
		err = todos.Store(todofile)
		handleError(err, "Error storing data:")
	case *complete == 0 && !*list:
		msg := "ID of the task you want to list as completed: "
		text, err := getInput(msg, os.Stdin, flag.Args()...)
		handleError(err, "Error getting input:")
		id, err := strconv.Atoi(text)
		handleError(err, "Error converting data:")
		err = todos.Complete(id)
		handleError(err, "Error function call:")
		err = todos.Store(todofile)
		handleError(err, "Error storing data:")
	case *complete > 0:
		err := todos.Complete(*complete)
		handleError(err, "Error function call:")

		err = todos.Store(todofile)
		handleError(err, "Error storing data:")
	case *delete > 0:
		err := todos.Delete(*delete)
		handleError(err, "Error function call:")
		err = todos.Store(todofile)
		handleError(err, "Error storing data:")
	case *list:
		todos.Show()
	default:
		fmt.Println("invalid command")
		os.Exit(1)
	}
}

func getInput(msg string, r io.Reader, args ...string) (string, error) {

	if len(args) > 0 {
		return strings.Join(args, " "), nil
	}

	fmt.Print(msg)

	scanner := bufio.NewScanner(r)
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		return "", err
	}

	text := scanner.Text()

	if len(text) == 0 {
		return "", errors.New("empty todo is not allowed")
	}

	return text, nil

}
