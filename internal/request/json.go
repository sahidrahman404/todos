package request

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/sahidrahman404/todos/internal/utils"
)

func DecodeJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)

	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		case errors.As(err, &syntaxError):
			return utils.BadJSONRequest(fmt.Sprintf("body contains badly-formed JSON (at character %d)", syntaxError.Offset))

		case errors.Is(err, io.ErrUnexpectedEOF):
			return utils.BadJSONRequest("body contains badly-formed JSON")

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return utils.BadJSONRequest(fmt.Sprintf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field))
			}
			return utils.BadJSONRequest(fmt.Sprintf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset))

		case errors.Is(err, io.EOF):
			return utils.BadJSONRequest("body must not be empty")

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return utils.BadJSONRequest(fmt.Sprintf("body contains unknown key %s", fieldName))

		case err.Error() == "http: request body too large":
			return utils.BadJSONRequest(fmt.Sprintf("body must not be larger than %d bytes", maxBytes))

		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return utils.BadJSONRequest(err.Error())
		}
	}

	err = dec.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		return utils.BadJSONRequest("body must only contain a single JSON value")
	}

	return nil
}
