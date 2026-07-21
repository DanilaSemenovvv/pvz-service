package storage

import (
	"context"

	"github.com/DanilaSemenovvv/pvz/internal/models"
)

type OrderStorage interface {
	Save(ctx context.Context, order models.Order) error
	GetByID(ctx context.Context, id int) (models.Order, error)
	GetByIDs(ctx context.Context, id []int) ([]models.Order, error)
	GetAll(ctx context.Context) ([]models.Order, error)
	Update(ctx context.Context, order models.Order) error
	DeleteByID(ctx context.Context, id int) error
	Close()
}
