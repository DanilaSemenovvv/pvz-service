package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/DanilaSemenovvv/pvz/internal/service"
)

type ReturnOrderRequest struct {
	ClientID  int   `json:"client_id"`
	OrdersIDs []int `json:"order_ids"`
}

func (h *Handler) ReturnOrder(w http.ResponseWriter, r *http.Request) {
	baseCtx := r.Context()
	ctx, cancel := context.WithTimeout(baseCtx, 3*time.Second)

	defer cancel()

	var req ReturnOrderRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		respondError(w, http.StatusBadRequest, "ошибка запроса")
		return
	}

	err = h.services.ProcessClient(ctx, req.ClientID, req.OrdersIDs, service.ActionReturn)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Ошибка возврата")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"msg": "Заказ возвращен"})
}
