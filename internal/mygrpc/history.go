package mygrpc

import (
	"context"
	"time"

	desc "github.com/DanilaSemenovvv/pvz/pkg/api/pvz/v1"
)

func (s *serverAPI) GetHistory(ctx context.Context, req *desc.GetHistoryRequest) (*desc.GetHistoryResponse, error) {
	history, err := s.services.GetOrderHistory(ctx, int(req.Limit), int(req.Offset))
	if err != nil {
		return nil, err
	}

	var pbOrders []*desc.Order

	for _, dataAboutOrder := range history {
		pbOrder := &desc.Order{
			Id:              int64(dataAboutOrder.OrderID),
			ClientId:        int64(dataAboutOrder.ClientID),
			Status:          string(dataAboutOrder.Status),
			StorageDeadline: dataAboutOrder.StorageDeadline.Format(time.RFC3339),
		}

		pbOrders = append(pbOrders, pbOrder)
	}

	return &desc.GetHistoryResponse{Orders: pbOrders}, nil
}
