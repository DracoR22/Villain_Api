package utils

import (
	"encoding/json"
	"net/http"
)

// --------------------------------------------------//TURN DATA INTO JSON//------------------------------------//
func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode((v))
}
