package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

type ErrorResponse struct {
	Errors []string `json:"errors"`
}

func WriteError(w http.ResponseWriter, statusCode int, msg string) {
	response := ErrorResponse{Errors: []string{msg}}
	if err := WriteJSON(w, statusCode, response); err != nil {
		log.Printf("JSON encode error: %v", err)
	}
}

func WriteValidationErrors(w http.ResponseWriter, errors []string) {
	response := ErrorResponse{Errors: errors}
	if err := WriteJSON(w, http.StatusBadRequest, response); err != nil {
		log.Printf("JSON encode error: %v", err)
	}
}

func WriteJSON(w http.ResponseWriter, statusCode int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(data)
}
