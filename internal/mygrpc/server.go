package mygrpc

import (
	"github.com/DanilaSemenovvv/pvz/internal/service"
	desc "github.com/DanilaSemenovvv/pvz/pkg/api/pvz/v1"
)

type serverAPI struct {
	desc.UnimplementedPVZServiceServer
	services *service.OrderService
}

func NewServer(srv *service.OrderService) *serverAPI {
	return &serverAPI{
		services: srv,
	}
}
