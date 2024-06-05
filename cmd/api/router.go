package main

import (
	"net/http"

	"github.com/sahidrahman404/todos/internal/database"
	"github.com/uptrace/bunrouter"
)

func (app *application) routes() http.Handler {
	router := bunrouter.New()

	// todo handlers
	todoHandler := NewTodoHandler(database.NewSqliteTodoStore(app.db.Bun))
	apiv1 := router.NewGroup("/api/v1", bunrouter.Use(app.ErrorHandler))
	apiv1.GET("/todos", todoHandler.HandleGetTodos)
	apiv1.POST("/todos", todoHandler.HandlePostTodo)
	apiv1.POST("/todos/children", todoHandler.HandlePostChildTodo)
	apiv1.DELETE("/todos/:id", todoHandler.HandleDeleteTodo)
	apiv1.PATCH("/todos/:id", todoHandler.HandleUpdateTodo)

	handler := http.Handler(router)
	handler = PanicHandler{Next: handler}

	return handler
}
