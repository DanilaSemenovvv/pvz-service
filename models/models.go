package models

import (
	"time"
)

type Order struct {
	OrderID  int         `json:"order_id"`
	ClientID int         `json:"client_id"`
	Status   OrderStatus `json:"status"`

	StorageDeadline time.Time `json:"storage_deadline"`
	DeliveredAt     time.Time `json:"delivered_at"`
	UpdatedAt       time.Time `json:"update_at"`
}

type OrderStatus string

const (
	StatusAccepted      OrderStatus = "принят от курьера"
	StatusDelivered     OrderStatus = "выдан клиенту"
	StatusReturnOnPvz   OrderStatus = "заказ вернули на пвз"
	StatusReturnCourier OrderStatus = "заказ вернули курьеру"
)
