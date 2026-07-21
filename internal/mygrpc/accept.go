package mygrpc

import (
	"context"
	"time"

	desc "github.com/DanilaSemenovvv/pvz/pkg/api/pvz/v1"
)

func (s *serverAPI) AcceptOrder(ctx context.Context, req *desc.AcceptOrderRequest) (*desc.AcceptOrderResponse, error) {
	deadline, err := time.Parse(time.RFC3339, req.StorageDeadline)
	if err != nil {
		return nil, err
	}

	err = s.services.AcceptOrder(ctx, int(req.OrderId), int(req.ClientId), deadline)
	if err != nil {
		return nil, err
	}

	return &desc.AcceptOrderResponse{Msg: "Заказ принят"}, nil
}
