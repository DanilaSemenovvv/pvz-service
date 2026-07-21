package mygrpc

import (
	"context"

	"github.com/DanilaSemenovvv/pvz/internal/service"
	desc "github.com/DanilaSemenovvv/pvz/pkg/api/pvz/v1"
)

func (s *serverAPI) DeliverOrder(ctx context.Context, req *desc.DeliverOrderRequest) (*desc.DeliverOrderResponse, error) {
	ordersIDs := make([]int, len(req.OrderIds))
	for i, id := range req.OrderIds {
		ordersIDs[i] = int(id)
	}

	err := s.services.ProcessClient(ctx, int(req.ClientId), ordersIDs, service.ActionDeliver)
	if err != nil {
		return nil, err
	}

	return &desc.DeliverOrderResponse{Msg: "Заказ доставлен"}, nil
}
