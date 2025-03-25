package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
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
	help := flag.Bool("help", false, "help by providing arguments list")

	flag.Usage = showHelp
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
	case flag.NFlag() == 0 || *help:
		showHelp()
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

func showHelp() {
	fmt.Println("Usage:")
	fmt.Println("  -add         Add a new task")
	fmt.Println("  -complete N  Mark task N as completed")
	fmt.Println("  -delete N    Delete task N")
	fmt.Println("  -list        Show all tasks")
	os.Exit(1)
}
