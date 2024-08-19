package todo

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alexeyco/simpletable"
	"os"
	"time"
)

type item struct {
	Task        string
	Done        bool
	CreatedAt   time.Time
	CompletedAt time.Time
}

type Todos []item

func (t *Todos) Add(task string) {
	todo := item{
		Task:        task,
		Done:        false,
		CreatedAt:   time.Now(),
		CompletedAt: time.Time{},
	}
	*t = append(*t, todo)
}

func (t *Todos) Complete(index int) error {
	todos := *t
	if index <= 0 || index > len(todos) {
		fmt.Println("invalid index")
	}
	todos[index-1].CompletedAt = time.Now()
	todos[index-1].Done = true
	*t = todos
	return nil
}

func (t *Todos) Delete(index int) error {
	todos := *t
	if index <= 0 || index > len(todos) {
		fmt.Println("invalid index")
	}
	*t = append(todos[:index-1], todos[index:]...)
	return nil
}

func (t *Todos) Get(index int) (string, error) {
	todos := *t
	if index <= 0 || index > len(todos) {
		return "", fmt.Errorf("invalid index")
	}
	return todos[index-1].Task, nil
}

func (t *Todos) Load(filename string) error {
	file, err := os.ReadFile(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}

	if len(file) == 0 {
		return err
	}
	err = json.Unmarshal(file, t)
	if err != nil {
		return err
	}

	return nil
}

func (t *Todos) Store(filename string) error {
	file, err := json.MarshalIndent(t, "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, file, 0644)
}

func (t *Todos) Print() {

	table := simpletable.New()

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "#"},
			{Align: simpletable.AlignCenter, Text: "Task"},
			{Align: simpletable.AlignCenter, Text: "Done"},
			{Align: simpletable.AlignLeft, Text: "Created At"},
			{Align: simpletable.AlignLeft, Text: "Completed At"},
		},
	}
	var cells [][]*simpletable.Cell

	for i, todo := range *t {
		i++
		cells = append(cells, []*simpletable.Cell{
			{Text: fmt.Sprintf("%d", i)},
			{Text: todo.Task},
			{Text: fmt.Sprintf("%t", todo.Done)},
			{Text: todo.CreatedAt.Format(time.RFC822)},
			{Text: todo.CompletedAt.Format(time.RFC822)},
		})
	}
	table.Body = &simpletable.Body{Cells: cells}

	table.SetStyle(simpletable.StyleUnicode)

	fmt.Println(table.String())
}
