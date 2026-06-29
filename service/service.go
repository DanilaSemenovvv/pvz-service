package service

import (
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/DanilaSemenovvv/pvz/models"
	"github.com/DanilaSemenovvv/pvz/storage"
)

type ActionType string

const (
	ActionDeliver ActionType = "deliver"
	ActionReturn  ActionType = "return"
)

type OrderService struct {
	storage storage.OrderStorage
}

func NewOrderService(st storage.OrderStorage) *OrderService {
	return &OrderService{
		storage: st,
	}
}

func (s *OrderService) AcceptOrder(orderID int, clientID int, deadline time.Time) error {

	if deadline.Before(time.Now()) {
		return errors.New("срок хранения в прошлом")
	}

	_, err := s.storage.GetByID(orderID)
	if err == nil {
		return errors.New("Заказ уже существует")
	}

	order := models.Order{
		ClientID:        clientID,
		OrderID:         orderID,
		Status:          models.StatusAccepted,
		StorageDeadline: deadline,
		UpdatedAt:       time.Now(),
	}

	err = s.storage.Save(order)
	if err != nil {
		return fmt.Errorf("Не удалось сохранить заказ: %w", err)
	}

	return nil

}

func (s *OrderService) ReturnToCourier(orderID int) error {
	order, err := s.storage.GetByID(orderID)
	if err != nil {
		return fmt.Errorf("Ошибка поиска заказа: %w", err)
	}

	if order.Status == models.StatusDelivered {
		return errors.New("Заказ выдан клиенту")
	}

	if order.StorageDeadline.After(time.Now()) {
		return errors.New("Срок хранения еще не истек")
	}

	err = s.storage.DeleteByID(orderID)
	if err != nil {
		return fmt.Errorf("Ошибка удаления заказа: %w", err)
	}

	return nil
}

func (s *OrderService) ProcessClient(clientID int, orderIDs []int, action ActionType) error {
	orders, err := s.storage.GetByIDs(orderIDs)
	if err != nil {
		return fmt.Errorf("Ошибка получения заказов: %w", err)
	}

	for _, order := range orders {
		if order.ClientID != clientID {
			return fmt.Errorf("заказ %d принадлежит другому клиенту", order.OrderID)
		}
	}

	switch action {
	case ActionDeliver:
		return s.deliverOrders(orders)
	case ActionReturn:
		return s.returnOrders(orders)
	default:
		return errors.New("Неизвестное действие")
	}
}

func (s *OrderService) deliverOrders(orders []models.Order) error {
	for _, order := range orders {
		if order.Status != models.StatusAccepted {
			return errors.New("Заказ не принят")
		}

		if order.StorageDeadline.Before(time.Now()) {
			return errors.New("Срок хранения истек")
		}
	}

	for _, order := range orders {
		order.Status = models.StatusDelivered
		order.DeliveredAt = time.Now()
		order.UpdatedAt = time.Now()

		err := s.storage.Update(order)
		if err != nil {
			return fmt.Errorf("Ошибка обновления состояния заказа: %w", err)
		}
	}

	return nil
}

func (s *OrderService) returnOrders(orders []models.Order) error {
	for _, order := range orders {
		if order.Status != models.StatusDelivered {
			return errors.New("Заказ еще не был выдан получателю или уже был возвращен")
		}
		if time.Since(order.DeliveredAt) > 48*time.Hour {
			return errors.New("Прошло более 48 часов после выдачи заказа")
		}
	}

	for _, order := range orders {
		order.Status = models.StatusReturnOnPvz
		order.UpdatedAt = time.Now()

		err := s.storage.Update(order)
		if err != nil {
			return fmt.Errorf("Ошибка обновления состояния заказа: %w", err)
		}
	}

	return nil
}

func (s *OrderService) GetOrderHistory(limit int, offset int) ([]models.Order, error) {
	orders, err := s.storage.GetAll()
	if err != nil {
		return nil, fmt.Errorf("Ошибка получения списка заказов: %w", err)
	}

	sort.Slice(orders, func(i, j int) bool {
		return orders[i].UpdatedAt.After(orders[j].UpdatedAt)
	})

	if offset >= len(orders) {
		return []models.Order{}, nil
	}

	end := offset + limit

	if end > len(orders) {
		end = len(orders)
	}

	return orders[offset:end], nil
}
