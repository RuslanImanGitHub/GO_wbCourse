package storage

import (
	itemdetails "L0/models/itemDetails"
	"L0/models/order"
	"context"
)

type Repository interface {
	Create(ctx context.Context, item *itemdetails.ItemDetails) error
	FindAllFromOrder(ctx context.Context, order order.Order) ([]itemdetails.ItemDetails, error)
}
