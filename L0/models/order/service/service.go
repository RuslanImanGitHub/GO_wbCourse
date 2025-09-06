package service

import (
	"L0/models/order"
	"L0/models/order/storage"
	"context"
	"fmt"
)

type Service struct {
	repository storage.Repository
}

func NewService(repository storage.Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) Create(ctx context.Context, order *order.Order) error {
	err := s.repository.Create(ctx, order)
	if err != nil {
		return fmt.Errorf("failed to create Order: %v", err)
	}
	return nil
}

func (s *Service) FindOne(ctx context.Context, id string) (order.Order, error) {
	order, err := s.repository.FindOne(ctx, id)
	if err != nil {
		return order, fmt.Errorf("failed to get order with uid %s: %v", id, err)
	}
	return order, nil
}

func (s *Service) FindAll(ctx context.Context) ([]order.Order, error) {
	orders, err := s.repository.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %v", err)
	}
	return orders, nil
}

func (s *Service) FindN(ctx context.Context, n int) ([]order.Order, error) {
	orders, err := s.repository.FindN(ctx, n)
	if err != nil {
		return nil, fmt.Errorf("failed to get N orders: %v", err)
	}
	return orders, nil
}
