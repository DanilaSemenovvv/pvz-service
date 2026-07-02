package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/DanilaSemenovvv/pvz/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStorage struct {
	db *pgxpool.Pool
}

func timeToNull(t time.Time) any {
	if t.IsZero() {
		return nil
	}

	return t
}

func NewPostgresStorage(connString string) (*PostgresStorage, error) {
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("Ошибка создания пула соединений: %w", err)
	}

	if err = pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("Ошибка, база не доступна: %w", err)
	}

	newTable := `
	CREATE TABLE IF NOT EXISTS orders (
    order_id BIGINT PRIMARY KEY,
    client_id BIGINT NOT NULL,
    status VARCHAR(50) NOT NULL,
    storage_deadline TIMESTAMP NOT NULL,
    delivered_at TIMESTAMP,
    updated_at TIMESTAMP NOT NULL
	);
	`

	_, err = pool.Exec(ctx, newTable)
	if err != nil {
		return nil, fmt.Errorf("Ошибка создания таблицы: %w", err)
	}

	return &PostgresStorage{db: pool}, nil
}

func (s *PostgresStorage) Save(order models.Order) error {
	ctx := context.Background()

	saveQuery := `
	INSERT INTO orders (order_id, client_id, status, storage_deadline, delivered_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := s.db.Exec(ctx, saveQuery,
		order.OrderID,
		order.ClientID,
		string(order.Status),
		order.StorageDeadline,
		timeToNull(order.DeliveredAt),
		order.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("Ошибка сохранения данных заказа: %w", err)
	}

	return nil
}

func (s *PostgresStorage) GetByID(id int) (models.Order, error) {
	ctx := context.Background()

	var order models.Order
	var statusStr string
	var deliverAt sql.NullTime

	query := "SELECT order_id, client_id, status, storage_deadline, delivered_at, updated_at FROM orders WHERE order_id = $1"

	err := s.db.QueryRow(ctx, query, id).Scan(
		&order.OrderID,
		&order.ClientID,
		&statusStr,
		&order.StorageDeadline,
		&deliverAt,
		&order.UpdatedAt,
	)
	if err != nil {
		return order, fmt.Errorf("Ошибка получения данных из БД: %w", err)
	}

	order.Status = models.OrderStatus(statusStr)
	if deliverAt.Valid {
		order.DeliveredAt = deliverAt.Time
	}

	return order, nil
}

func (s *PostgresStorage) GetAll() ([]models.Order, error) {
	ctx := context.Background()

	query := "SELECT order_id, client_id, status, storage_deadline, delivered_at, updated_at FROM orders"

	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("Ошибка получения данных из БД: %w", err)
	}

	var orders []models.Order

	defer rows.Close()
	for rows.Next() {
		var order models.Order
		var statusStr string
		var deliverAt sql.NullTime

		err := rows.Scan(
			&order.OrderID,
			&order.ClientID,
			&statusStr,
			&order.StorageDeadline,
			&deliverAt,
			&order.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("Ошибка чтения данных: %w", err)
		}

		order.Status = models.OrderStatus(statusStr)
		if deliverAt.Valid {
			order.DeliveredAt = deliverAt.Time
		}

		orders = append(orders, order)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при чтении строк: %w", err)
	}

	return orders, nil
}

func (s *PostgresStorage) GetByIDs(ids []int) ([]models.Order, error) {
	ctx := context.Background()

	query := "SELECT order_id, client_id, status, storage_deadline, delivered_at, updated_at FROM orders WHERE order_id = ANY($1)"

	rows, err := s.db.Query(ctx, query, ids)
	if err != nil {
		return nil, fmt.Errorf("Ошибка получения данных из БД: %w", err)
	}

	var orders []models.Order

	defer rows.Close()
	for rows.Next() {
		var order models.Order
		var statusStr string
		var deliverAt sql.NullTime

		err := rows.Scan(
			&order.OrderID,
			&order.ClientID,
			&statusStr,
			&order.StorageDeadline,
			&deliverAt,
			&order.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("Ошибка чтения данных: %w", err)
		}

		order.Status = models.OrderStatus(statusStr)
		if deliverAt.Valid {
			order.DeliveredAt = deliverAt.Time
		}

		orders = append(orders, order)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при чтении строк: %w", err)
	}

	return orders, nil
}

func (s *PostgresStorage) Update(order models.Order) error {
	ctx := context.Background()

	updateQuery := "UPDATE orders SET status = $1, delivered_at = $2, updated_at = $3 WHERE order_id = $4"

	_, err := s.db.Exec(ctx, updateQuery,
		string(order.Status),
		timeToNull(order.DeliveredAt),
		order.UpdatedAt,
		order.OrderID,
	)
	if err != nil {
		return fmt.Errorf("Ошибка апдейта: %w", err)
	}
	return nil
}

func (s *PostgresStorage) DeleteByID(id int) error {
	ctx := context.Background()

	deleteQuery := "DELETE FROM orders WHER order_id = $1"

	_, err := s.db.Exec(ctx, deleteQuery, id)
	if err != nil {
		return fmt.Errorf("Ошибка удаления заказа: %w", err)
	}

	return nil
}
