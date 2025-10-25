package storage

import (
	"L0/backend/models/order"
	"context"
)

type Repository interface {
	Create(ctx context.Context, order *order.Order) error
	FindOne(ctx context.Context, id string) (order.Order, error)
	FindAll(ctx context.Context) ([]order.Order, error)
	FindN(ctx context.Context, n int) ([]order.Order, error)
}
