package storage

import (
	"L0/models/deliveryAddress"
	itemdetails "L0/models/itemDetails"
	"L0/models/order"
	"L0/models/payment"
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, pool *pgxpool.Pool, order *order.Order, delivery *deliveryAddress.DeliveryAddress,
		payment *payment.Payment, itemDetails []itemdetails.ItemDetails) error
	FindOne(ctx context.Context, pool *pgxpool.Pool, id string) (order.Order, error)
	FindAll(ctx context.Context, pool *pgxpool.Pool) ([]order.Order, error)
	FindN(ctx context.Context, pool *pgxpool.Pool, n int) ([]order.Order, error)
}
