package api

import (
	"encoding/json"
	"net/http"
)

func write(w http.ResponseWriter, msg string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(Error{
		Message: msg,
	})
}

func WriteBadRequest(w http.ResponseWriter, msg string) {
	write(w, msg, http.StatusBadRequest)
}

func WriteUnauthorized(w http.ResponseWriter, msg string) {
	write(w, msg, http.StatusUnauthorized)
}

func WriteForbidden(w http.ResponseWriter, msg string) {
	write(w, msg, http.StatusForbidden)
}
