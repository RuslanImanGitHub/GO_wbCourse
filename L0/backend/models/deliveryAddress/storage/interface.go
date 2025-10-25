package storage

import (
	"L0/backend/models/deliveryAddress"
	"context"
)

type Repository interface {
	Create(ctx context.Context, delivery *deliveryAddress.DeliveryAddress) (int, error)
	FindFromId(ctx context.Context, id int) (deliveryAddress.DeliveryAddress, error)
	FindAll(ctx context.Context) ([]deliveryAddress.DeliveryAddress, error)
}
