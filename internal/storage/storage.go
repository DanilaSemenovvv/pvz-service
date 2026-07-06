package storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/DanilaSemenovvv/pvz/internal/models"
)

type OrderStorage interface {
	Save(ctx context.Context, order models.Order) error
	GetByID(ctx context.Context, id int) (models.Order, error)
	GetByIDs(ctx context.Context, id []int) ([]models.Order, error)
	GetAll(ctx context.Context) ([]models.Order, error)
	Update(ctx context.Context, order models.Order) error
	DeleteByID(ctx context.Context, id int) error
}

type JSONStorage struct {
	filePath string
	cache    map[int]models.Order
}

func NewJSONStorage(path string) (*JSONStorage, error) {
	st := &JSONStorage{
		filePath: path,
		cache:    make(map[int]models.Order),
	}

	data, err := os.ReadFile(st.filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return st, nil
		}

		return nil, fmt.Errorf("Ошибка чтения БД: %w", err)
	}

	var orders = make([]models.Order, 0, len(st.cache))
	err = json.Unmarshal(data, &orders)
	if err != nil {
		return nil, fmt.Errorf("Ошибка JSON: %w", err)
	}

	for _, order := range orders {
		st.cache[order.OrderID] = order
	}

	return st, nil

}

func (s *JSONStorage) GetAll() ([]models.Order, error) {
	var orders []models.Order
	for _, order := range s.cache {
		orders = append(orders, order)
	}

	return orders, nil
}

func (s *JSONStorage) Save(order models.Order) error {
	s.cache[order.OrderID] = order

	orders, err := s.GetAll()
	if err != nil {
		return fmt.Errorf("Ошибка получения списка заказов: %w", err)
	}

	data, err := json.MarshalIndent(orders, "", " ")
	if err != nil {
		return fmt.Errorf("Ошибка кодирования данных: %w", err)
	}

	err = os.WriteFile(s.filePath, data, 0644)
	if err != nil {
		return fmt.Errorf("Ошибка записи в файл: %w", err)
	}

	return nil
}

func (s *JSONStorage) GetByID(id int) (models.Order, error) {
	order, ok := s.cache[id]
	if !ok {
		return models.Order{}, errors.New("Такого заказа не существует")
	}

	return order, nil
}

func (s *JSONStorage) GetByIDs(ids []int) ([]models.Order, error) {
	orders := make([]models.Order, 0, len(ids))

	for _, id := range ids {
		if order, ok := s.cache[id]; ok {
			orders = append(orders, order)
		} else {
			return nil, fmt.Errorf("заказ с ID %d не найден", id)
		}
	}

	return orders, nil
}

func (s *JSONStorage) Update(order models.Order) error {
	return s.Save(order)
}

func (s *JSONStorage) DeleteByID(id int) error {
	delete(s.cache, id)

	orders, err := s.GetAll()
	if err != nil {
		return fmt.Errorf("Ошибка получения списка заказов после удаления: %w", err)
	}

	data, err := json.MarshalIndent(orders, "", " ")
	if err != nil {
		return fmt.Errorf("Ошибка кодирования данных: %w", err)
	}

	err = os.WriteFile(s.filePath, data, 0644)
	if err != nil {
		return fmt.Errorf("Ошибка записи в файл: %w", err)
	}

	return nil
}
