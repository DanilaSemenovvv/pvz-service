package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/DanilaSemenovvv/pvz/internal/service"
)

type ReturnOrderRequest struct {
	ClientID  int   `json:"client_id"`
	OrdersIDs []int `json:"order_ids"`
}
type DeliverOrderRequest struct {
	ClientID  int   `json:"client_id"`
	OrdersIDs []int `json:"order_ids"`
}
type AcceptOrderRequest struct {
	OrderID         int    `json:"order_id"`
	ClientID        int    `json:"client_id"`
	StorageDeadline string `json:"storage_deadline"`
}

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

func (h *Handler) ReturnOrder(w http.ResponseWriter, r *http.Request) {
	var req ReturnOrderRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		respondError(w, http.StatusBadRequest, "ошибка запроса")
		return
	}

	err = h.services.ProcessClient(req.ClientID, req.OrdersIDs, service.ActionReturn)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Ошибка возврата")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"msg": "Заказ возвращен"})
}

func (h *Handler) DeliverOrder(w http.ResponseWriter, r *http.Request) {
	var req DeliverOrderRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		respondError(w, http.StatusBadRequest, "ошибка запроса")
		return
	}

	err = h.services.ProcessClient(req.ClientID, req.OrdersIDs, service.ActionDeliver)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "ошибка доставки")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"msg": "Заказ доставлен"})
}

func (h *Handler) GetHistory(w http.ResponseWriter, r *http.Request) {
	limit := 10
	offset := 0

	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")
	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "параметр limit должен быть числом")
			return
		}
		limit = parsedLimit
	}

	if offsetStr != "" {
		parsedOffset, err := strconv.Atoi(offsetStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "параметр offset должен быть числом")
			return
		}
		offset = parsedOffset
	}

	orders, err := h.services.GetOrderHistory(limit, offset)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Ошибка получения заказа")
		return
	}

	respondJSON(w, http.StatusOK, orders)
}

func (h *Handler) AcceptOrder(w http.ResponseWriter, r *http.Request) {
	var req AcceptOrderRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		respondError(w, http.StatusBadRequest, "")
		return
	}

	storageDeadline, err := time.Parse(time.RFC3339, req.StorageDeadline)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Ошибка формата даты!")
		return
	}

	err = h.services.AcceptOrder(req.OrderID, req.ClientID, storageDeadline)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Не удалось принять заказ")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"msg": "Заказ принят"})
}

func HandlePing(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong! Система ПВЗ работает!"))
}
