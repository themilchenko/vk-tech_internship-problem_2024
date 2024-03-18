package pkg

import (
	"encoding/json"
	"net/http"
)

func HandleError(w http.ResponseWriter, message string, statusCode int) {
	errorResponse := struct {
		Error string `json:"error"`
	}{
		Error: message,
	}

	jsonResponse, err := json.Marshal(errorResponse)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(statusCode)
	w.Write(jsonResponse)
}
