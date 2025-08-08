package api

import (
	"encoding/json"
	"net/http"
)

var responseInternalError = newErrorResponse(500, "Internal server error")

type HTTPError struct {
	Code int
	Body Error
}

type Error struct {
	Message string `json:"message,omitempty"`
}

func (r *HTTPError) Write(w http.ResponseWriter) error {
	w.WriteHeader(r.Code)
	return respondJSON(r.Body, w)
}

func newErrorResponse(code int, message string) *HTTPError {
	return &HTTPError{
		Code: code,
		Body: Error{
			Message: message,
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
