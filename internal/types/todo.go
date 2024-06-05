package types

import (
	"net/http"
	"time"

	"github.com/sahidrahman404/todos/internal/utils"
	"github.com/sahidrahman404/todos/internal/validator"
)

var taskMinLen = 5

type TodoFilters struct {
	GetCompleted *bool
}

type CreateTodoParams struct {
	Deadline *time.Time `json:"deadline"`
	Task     string     `json:"task"`
}

type CreateChildTodoParams struct {
	Deadline   *time.Time `json:"deadline"`
	Task       string     `json:"task"`
	TodoParent int        `json:"todoParent"`
}

type UpdateTodoParams struct {
	Deadline    *time.Time `json:"deadline"`
	Task        *string    `json:"task"`
	IsCompleted *bool      `json:"isCompleted"`
}

type Todo struct {
	Deadline    *time.Time `json:"deadline"`
	TodoParent  *int       `json:"todoParent"`
	Task        string     `json:"task"`
	Children    []*Todo    `bun:"-" json:"children"`
	IsCompleted bool       `json:"isCompleted"`
	ID          int        `bun:",pk,autoincrement" json:"id"`
	Version     int        `json:"version"`
}

func NewTodoFromCreateTodoParams(params CreateTodoParams) (*Todo, error) {
	todo := &Todo{
		Deadline: params.Deadline,
		Task:     params.Task,
	}
	validator := todo.validate()
	if validator.HasErrors() {
		return nil, utils.BadInputRequest(validator.FieldErrors)
	}
	return todo, nil
}

func NewChildTodoFromCreateChildTodoParams(params CreateChildTodoParams) (*Todo, error) {
	todo := &Todo{
		Deadline:   params.Deadline,
		Task:       params.Task,
		TodoParent: &params.TodoParent,
	}

	validator := todo.validate()
	validator.CheckField(*todo.TodoParent > 0, "todoParent", "todoParent cannot be empty")
	if validator.HasErrors() {
		return nil, utils.BadInputRequest(validator.FieldErrors)
	}
	return todo, nil
}

func NewTodoFromUpdateTodoParams(params UpdateTodoParams, currTodo *Todo) (*Todo, error) {
	if params.Deadline != nil {
		currTodo.Deadline = params.Deadline
	}
	if params.Task != nil {
		currTodo.Task = *params.Task
	}

	if params.IsCompleted != nil {
		currTodo.IsCompleted = *params.IsCompleted
	}

	validator := currTodo.validate()
	if validator.HasErrors() {
		return nil, utils.BadInputRequest(validator.FieldErrors)
	}
	return currTodo, nil
}

func TransformTodos(
	todos []*Todo,
) []*Todo {
	nestedTodos := []*Todo{}
	nestedTodoMap := map[int]int{}

	for i, v := range todos {
		nestedTodoMap[v.ID] = i
		if v.TodoParent != nil {
			parentPathIndex, ok := nestedTodoMap[*v.TodoParent]
			if !ok {
				v = nil
				continue
			}
			todos[parentPathIndex].Children = append(
				todos[parentPathIndex].Children,
				v,
			)
			continue
		}
		nestedTodos = append(nestedTodos, v)
	}
	return nestedTodos
}

func TodoFiltersFromRequest(r *http.Request) (*TodoFilters, error) {
	var todoFilters TodoFilters
	query := r.URL.Query()
	switch query.Get("get_completed") {
	case "":
		todoFilters.GetCompleted = nil
	case "true":
		filter := true
		todoFilters.GetCompleted = &filter
	case "false":
		filter := false
		todoFilters.GetCompleted = &filter
	default:
		return nil, utils.BadInputRequest(map[string]string{
			"get_completed": "type should be boolean",
		})
	}
	return &todoFilters, nil
}

func (t Todo) validate() validator.Validator {
	var validator validator.Validator
	validator.CheckField(len(t.Task) > taskMinLen, "task", "Task is too sort")

	return validator
}
