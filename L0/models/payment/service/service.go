package service

import (
	"L0/models/payment"
	"L0/models/payment/storage"
	"context"
	"fmt"
)

type Service struct {
	repository storage.Repository
}

func NewService(repository storage.Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) Create(ctx context.Context, payment *payment.Payment) error {
	err := s.repository.Create(ctx, payment)
	if err != nil {
		return fmt.Errorf("failed to create Payment: %v", err)
	}
	return nil
}

func (s *Service) FindOneFromTransaction(ctx context.Context, transaction string) (payment.Payment, error) {
	order, err := s.repository.FindOneFromTransaction(ctx, transaction)
	if err != nil {
		return order, fmt.Errorf("failed to get order with uid %s: %v", transaction, err)
	}
	return order, nil
}

func (s *Service) FindAll(ctx context.Context) ([]payment.Payment, error) {
	orders, err := s.repository.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get payments: %v", err)
	}
	return orders, nil
}
