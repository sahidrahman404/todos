package utils

import (
	"database/sql"
	"errors"
	"io"
	"net/http"
)

var ErrEditConflict = errors.New("edit conflict")

type HTTPError struct {
	Code     string            `json:"code"`
	Messages map[string]string `json:"messages"`
	Message  string            `json:"message"`

	StatusCode int
}

func (e HTTPError) Error() string {
	return e.Message
}

func BadJSONRequest(msg string) HTTPError {
	return HTTPError{
		StatusCode: http.StatusBadRequest,

		Code:    "bad_json_request",
		Message: msg,
	}
}

func BadInputRequest(msg map[string]string) HTTPError {
	return HTTPError{
		StatusCode: http.StatusBadRequest,

		Code:     "bad_input",
		Messages: msg,
	}
}

func NewHTTPError(err error) HTTPError {
	switch err {
	case io.EOF:
		return HTTPError{
			StatusCode: http.StatusBadRequest,

			Code:    "eof",
			Message: "EOF reading HTTP request body",
		}
	case sql.ErrNoRows:
		return HTTPError{
			StatusCode: http.StatusNotFound,

			Code:    "not_found",
			Message: "Page Not Found",
		}
	case ErrEditConflict:
		return HTTPError{
			StatusCode: http.StatusConflict,

			Code:    "edit_conflict",
			Message: "Edit conflict",
		}
	}

	return HTTPError{
		StatusCode: http.StatusInternalServerError,

		Code:    "internal",
		Message: "Internal server error",
	}
}
