package util

import (
	"encoding/json"
	"net/http"
)

// WriteError writes an error response to the client with the given HTTP status code and message.
func WriteError(w http.ResponseWriter, statusCode int, message string) {
	response := map[string]string{"error": message}
	WriteJSON(w, statusCode, response)
}

// WriteJSON serializes data to JSON and writes it to the response writer with the given HTTP status code.
func WriteJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to serialize data to JSON")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(jsonData)
}
