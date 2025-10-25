package storage

import (
	itemdetails "L0/backend/models/itemDetails"
	"L0/backend/models/order"
	"context"
)

type Repository interface {
	Create(ctx context.Context, item *itemdetails.ItemDetails) error
	FindAllFromOrder(ctx context.Context, order order.Order) ([]itemdetails.ItemDetails, error)
}
