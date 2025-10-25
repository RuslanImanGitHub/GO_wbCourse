package service

import (
	"L0/backend/models/deliveryAddress"
	"L0/backend/models/deliveryAddress/storage"
	"context"
	"fmt"
)

type Service struct {
	repository storage.Repository
}

func NewService(repository storage.Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) Create(ctx context.Context, delivery *deliveryAddress.DeliveryAddress) (int, error) {
	delivery_id, err := s.repository.Create(ctx, delivery)
	if err != nil {
		return 0, fmt.Errorf("failed to create delivery: %v", err)
	}
	return delivery_id, nil
}

func (s *Service) FindFromId(ctx context.Context, id int) (deliveryAddress.DeliveryAddress, error) {
	delivery, err := s.repository.FindFromId(ctx, id)
	if err != nil {
		return delivery, fmt.Errorf("failed to get delivery for %d: %v", id, err)
	}
	return delivery, nil
}

func (s *Service) FindAll(ctx context.Context) ([]deliveryAddress.DeliveryAddress, error) {
	deliveries, err := s.repository.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get deliveries: %v", err)
	}
	return deliveries, nil
}
