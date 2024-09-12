package httpServer

import (
	"encoding/json"
	"net/http"
)

func MethodNotAllowed(w http.ResponseWriter) {
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func BadRequest(w http.ResponseWriter, message string) {
	http.Error(w, message, http.StatusBadRequest)
}

func Ok(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusOK)
	jsonResp, _ := json.Marshal(message)
	w.Write(jsonResp)
}

func StatusInternalServerError(w http.ResponseWriter, message string) {
	http.Error(w, message, http.StatusInternalServerError)
}
