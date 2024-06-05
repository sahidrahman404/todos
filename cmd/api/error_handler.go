package main

import (
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/sahidrahman404/todos/internal/utils"
	"github.com/uptrace/bunrouter"
)

func (app *application) ErrorHandler(next bunrouter.HandlerFunc) bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, r bunrouter.Request) error {
		// Call the next handler on the chain to get the error.
		err := next(w, r)

		switch err := err.(type) {
		case nil:
			// no error
		case utils.HTTPError: // already a HTTPError
			w.WriteHeader(err.StatusCode)
			_ = bunrouter.JSON(w, err)
		default:

			message := err.Error()
			method := r.Method
			url := r.URL.String()
			trace := string(debug.Stack())

			requestAttrs := slog.Group("request", "method", method, "url", url)
			httpErr := utils.NewHTTPError(err)
			w.WriteHeader(httpErr.StatusCode)
			app.logger.Error(message, requestAttrs, "trace", trace)
			_ = bunrouter.JSON(w, httpErr)
		}

		return err // return the err in case there other middlewares
	}
}
