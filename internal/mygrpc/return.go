package mygrpc

import (
	"context"

	"github.com/DanilaSemenovvv/pvz/internal/service"
	desc "github.com/DanilaSemenovvv/pvz/pkg/api/pvz/v1"
)

func (s *serverAPI) ReturnOrder(ctx context.Context, req *desc.ReturnOrderRequest) (*desc.ReturnOrderResponse, error) {
	ordersIDs := make([]int, len(req.OrderIds))
	for i, id := range req.OrderIds {
		ordersIDs[i] = int(id)
	}

	err := s.services.ProcessClient(ctx, int(req.ClientId), ordersIDs, service.ActionReturn)
	if err != nil {
		return nil, err
	}

	return &desc.ReturnOrderResponse{Msg: "Заказ возвращен"}, nil
}
