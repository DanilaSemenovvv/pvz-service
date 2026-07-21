package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/DanilaSemenovvv/pvz/internal/service"
)

type DeliverOrderRequest struct {
	ClientID  int   `json:"client_id"`
	OrdersIDs []int `json:"order_ids"`
}

func (h *Handler) DeliverOrder(w http.ResponseWriter, r *http.Request) {
	baseCtx := r.Context()
	ctx, cancel := context.WithTimeout(baseCtx, 3*time.Second)

	defer cancel()

	var req DeliverOrderRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		respondError(w, http.StatusBadRequest, "ошибка запроса")
		return
	}

	err = h.services.ProcessClient(ctx, req.ClientID, req.OrdersIDs, service.ActionDeliver)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "ошибка доставки")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"msg": "Заказ доставлен"})
}
