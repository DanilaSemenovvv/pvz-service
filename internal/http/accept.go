package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type AcceptOrderRequest struct {
	OrderID         int    `json:"order_id"`
	ClientID        int    `json:"client_id"`
	StorageDeadline string `json:"storage_deadline"`
}

func (h *Handler) AcceptOrder(w http.ResponseWriter, r *http.Request) {
	baseCtx := r.Context()
	ctx, cancel := context.WithTimeout(baseCtx, 3*time.Second)

	defer cancel()

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

	err = h.services.AcceptOrder(ctx, req.OrderID, req.ClientID, storageDeadline)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Не удалось принять заказ")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"msg": "Заказ принят"})
}
