package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"
)

func (h *Handler) GetHistory(w http.ResponseWriter, r *http.Request) {
	baseCtx := r.Context()
	ctx, cancel := context.WithTimeout(baseCtx, 3*time.Second)

	defer cancel()

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

	orders, err := h.services.GetOrderHistory(ctx, limit, offset)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Ошибка получения заказа")
		return
	}

	respondJSON(w, http.StatusOK, orders)
}
