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

	flag.Parse()

	todos := &todo.Todos{}

	if err := todos.Load(todofile); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	switch {
	case *add:
		task, err := getInput(os.Stdin, flag.Args()...)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		todos.Add(task)
		err = todos.Store(todofile)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	case *complete > 0:
		err := todos.Complete(*complete)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		err = todos.Store(todofile)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	case *delete > 0:
		err := todos.Delete(*delete)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		err = todos.Store(todofile)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	case *list:
		todos.Show()
	default:
		fmt.Println("invalid command")
		os.Exit(1)
	}
}

func getInput(r io.Reader, args ...string) (string, error) {

	if len(args) > 0 {
		return strings.Join(args, " "), nil
	}

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
