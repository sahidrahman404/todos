package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/sahidrahman404/todos/internal/types"
	"github.com/sahidrahman404/todos/internal/utils"
	"github.com/uptrace/bun"
)

type TodoStorer interface {
	GetTodo(context.Context, string) (*types.Todo, error)
	GetTodos(context.Context, *types.TodoFilters) ([]*types.Todo, error)
	CreateTodo(context.Context, *types.Todo) (*types.Todo, error)
	DeleteTodo(context.Context, string) error
	UpdateTodo(context.Context, *types.Todo) (*types.Todo, error)
}

type SqliteTodoStore struct {
	client *bun.DB
}

func NewSqliteTodoStore(s *bun.DB) *SqliteTodoStore {
	return &SqliteTodoStore{
		client: s,
	}
}

func (s *SqliteTodoStore) GetTodo(ctx context.Context, id string) (*types.Todo, error) {
	var todo types.Todo
	err := s.client.NewSelect().
		Model(&todo).
		Where("? = ?", bun.Ident("id"), id).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &todo, nil
}

func (s *SqliteTodoStore) GetTodos(ctx context.Context, todoFilters *types.TodoFilters) ([]*types.Todo, error) {
	var todos []*types.Todo
	q, err := todoFilterQuery(s.client.NewSelect().Model(&todos).Order("id ASC"), todoFilters)
	if err != nil {
		return nil, err
	}
	err = q.Scan(ctx)
	if err != nil {
		return nil, err
	}
	todos = types.TransformTodos(todos)
	return todos, nil
}

func (s *SqliteTodoStore) CreateTodo(ctx context.Context, user *types.Todo) (*types.Todo, error) {
	_, err := s.client.NewInsert().
		Model(user).
		Exec(ctx)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *SqliteTodoStore) DeleteTodo(ctx context.Context, id string) error {
	_, err := s.client.NewDelete().
		Table("todos").
		Where("? = ?", bun.Ident("id"), bun.Safe(id)).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (s *SqliteTodoStore) UpdateTodo(ctx context.Context, todo *types.Todo) (*types.Todo, error) {
	originVersion := todo.Version
	todo.Version = todo.Version + 1
	_, err := s.client.NewUpdate().
		Model(todo).
		Where("? = ?", bun.Ident("id"), bun.Safe(fmt.Sprint(todo.ID))).
		Where("? = ?", bun.Ident("version"), bun.Safe(fmt.Sprint(originVersion))).
		Exec(ctx)
	if err != nil {
		switch {
		case errors.Is(err, utils.ErrEditConflict):
			return nil, utils.ErrEditConflict
		default:
			return nil, err
		}
	}
	return todo, nil
}

func todoFilterQuery(q *bun.SelectQuery, todoFilters *types.TodoFilters) (*bun.SelectQuery, error) {
	if todoFilters.GetCompleted == nil {
		return q, nil
	}
	q = q.Where("? = ?", bun.Ident("is_completed"), bun.Safe(fmt.Sprint(*todoFilters.GetCompleted)))
	return q, nil
}
