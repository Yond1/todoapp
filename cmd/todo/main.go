package main

import (
	todo "TodoCLI"
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

var (
	addTodo      = flag.String("add", "", "add todo")
	completeTodo = flag.Int("complete", 0, "complete todo")
	deleteTodo   = flag.Int("delete", 0, "delete todo")
	getTodo      = flag.Int("get", 0, "get todo")
	print        = flag.Bool("print", false, "print todo")
)

const (
	todoFile   = ".todo.json"
	emptyValue = ""
)

func main() {

	flag.Parse()

	todos := &todo.Todos{}

	if err := todos.Load(todoFile); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	switch {
	case *addTodo != emptyValue:

		task, err := getInput(os.Stdin, *addTodo)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		todos.Add(task)

		err = todos.Store(todoFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

	case *completeTodo > 0:
		todos.Complete(*completeTodo)
		err := todos.Store(todoFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
	case *deleteTodo > 0:
		todos.Delete(*deleteTodo)
		err := todos.Store(todoFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
	case *getTodo > 0:
		todos.Get(*getTodo)
		err := todos.Store(todoFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
	case *print:
		todos.Print()
	default:
		fmt.Fprintln(os.Stdout, "invalid command")
		os.Exit(0)
	}

}

func getInput(r io.Reader, args ...string) (string, error) {
	if len(args) > 0 {
		return strings.Join(args, " "), nil
	}

	scanner := bufio.NewScanner(r)
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		return emptyValue, err
	}

	if len(scanner.Text()) == 0 {
		return emptyValue, errors.New("empty input")
	}

	return scanner.Text(), nil

}
