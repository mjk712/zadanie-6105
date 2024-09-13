package handlers

import (
	"api/internal/http-server/services"
	"encoding/json"
	"net/http"
)

func Ping(w http.ResponseWriter, r *http.Request) {
	pong := services.Pong()
	json.NewEncoder(w).Encode(pong)
	w.WriteHeader(http.StatusOK)
}
