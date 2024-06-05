package main

import (
	"context"
	"net/http"
	"time"

	"github.com/sahidrahman404/todos/internal/database"
	"github.com/sahidrahman404/todos/internal/request"
	"github.com/sahidrahman404/todos/internal/types"
	"github.com/uptrace/bunrouter"
)

const defaultTimeout = 3 * time.Second

type TodoHandler struct {
	todoStorer database.TodoStorer
}

func NewTodoHandler(todoStorer database.TodoStorer) *TodoHandler {
	return &TodoHandler{
		todoStorer: todoStorer,
	}
}

func (h *TodoHandler) HandleGetTodos(w http.ResponseWriter, r bunrouter.Request) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	todoFilters, err := types.TodoFiltersFromRequest(r.Request)
	if err != nil {
		return err
	}
	todos, err := h.todoStorer.GetTodos(ctx, todoFilters)
	if err != nil {
		return err
	}
	return bunrouter.JSON(w, envelope{
		"data": todos,
	})
}

func (h *TodoHandler) HandlePostTodo(w http.ResponseWriter, r bunrouter.Request) error {
	var params types.CreateTodoParams
	if err := request.DecodeJSON(w, r.Request, &params); err != nil {
		return err
	}
	todo, err := types.NewTodoFromCreateTodoParams(params)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	todo, err = h.todoStorer.CreateTodo(ctx, todo)
	if err != nil {
		return err
	}
	return bunrouter.JSON(w, envelope{
		"data": todo,
	})
}

func (h *TodoHandler) HandlePostChildTodo(w http.ResponseWriter, r bunrouter.Request) error {
	var params types.CreateChildTodoParams
	if err := request.DecodeJSON(w, r.Request, &params); err != nil {
		return err
	}
	todo, err := types.NewChildTodoFromCreateChildTodoParams(params)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	todo, err = h.todoStorer.CreateTodo(ctx, todo)
	if err != nil {
		return err
	}
	return bunrouter.JSON(w, envelope{
		"data": todo,
	})
}

func (h *TodoHandler) HandleDeleteTodo(w http.ResponseWriter, r bunrouter.Request) error {
	id := r.Param("id")

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	err := h.todoStorer.DeleteTodo(ctx, id)
	if err != nil {
		return err
	}
	return bunrouter.JSON(w, map[string]interface{}{
		"success": true,
	})
}

func (h *TodoHandler) HandleUpdateTodo(w http.ResponseWriter, r bunrouter.Request) error {
	id := r.Param("id")
	var params types.UpdateTodoParams

	if err := request.DecodeJSON(w, r.Request, &params); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	todo, err := h.todoStorer.GetTodo(ctx, id)
	if err != nil {
		return err
	}

	todo, err = types.NewTodoFromUpdateTodoParams(params, todo)
	if err != nil {
		return err
	}

	todo, err = h.todoStorer.UpdateTodo(ctx, todo)
	if err != nil {
		return err
	}
	return bunrouter.JSON(w, envelope{
		"data": todo,
	})
}
