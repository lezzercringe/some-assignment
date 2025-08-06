package api

import (
	"encoding/json"
	"net/http"
)

var responseInternalError = newErrorResponse(500, "Internal server error")

type field struct {
	Name    string
	Message string
}

type ErrorBody struct {
	Message string  `json:"message,omitempty"`
	Fields  []field `json:"fields,omitempty"`
}

type HTTPError struct {
	Code int
	Body ErrorBody
}

func (r *HTTPError) Write(w http.ResponseWriter) error {
	w.WriteHeader(r.Code)
	return respondJSON(r.Body, w)
}

func (r *HTTPError) WithField(name, message string) *HTTPError {
	r.Body.Fields = append(r.Body.Fields, field{
		Name:    name,
		Message: message,
	})

	return r
}

func newErrorResponse(code int, message string) *HTTPError {
	return &HTTPError{
		Code: code,
		Body: ErrorBody{
			Message: message,
			Fields:  nil,
		},
	}
}

func respondJSON(v any, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	encoded, err := json.Marshal(v)
	if err != nil {
		return err
	}

	w.Write(encoded)
	return nil
}
