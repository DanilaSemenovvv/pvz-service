package handler

import (
	"encoding/json"
	"net/http"

	"github.com/DanilaSemenovvv/pvz/internal/service"
)

type Handler struct {
	services *service.OrderService
}

func NewHandler(srv *service.OrderService) *Handler {
	return &Handler{
		services: srv,
	}
}

func respondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func respondError(w http.ResponseWriter, status int, msg string) {
	errorData := map[string]string{"error": msg}
	respondJSON(w, status, errorData)
}

func HandlePing(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong! Система ПВЗ работает!"))
}
