package storage

import (
	"L0/backend/models/payment"
	"context"
)

type Repository interface {
	Create(ctx context.Context, payment *payment.Payment) error
	FindOneFromTransaction(ctx context.Context, transaction string) (payment.Payment, error)
	FindAll(ctx context.Context) ([]payment.Payment, error)
}
