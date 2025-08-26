package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// Error response from the server
func Error(w http.ResponseWriter, statusCode int, text string) {
	slog.Error(text)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(map[string]string{"error": text}); err != nil {
		slog.Error(err.Error())
	}
}

// Success response
func SendResponse(w http.ResponseWriter, r *http.Request, object interface{}, textInfo string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if object == nil {
		m := make(map[string]string)
		m["status"] = textInfo
		object = m
	}

	if err := json.NewEncoder(w).Encode(object); err != nil {
		Error(w, 500, "Error encoding a file")
		return
	}

	slog.Info(textInfo)
}
