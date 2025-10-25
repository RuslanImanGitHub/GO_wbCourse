package service

import (
	itemdetails "L0/backend/models/itemDetails"
	"L0/backend/models/itemDetails/storage"
	"L0/backend/models/order"
	"context"
	"fmt"
)

type Service struct {
	repository storage.Repository
}

func NewService(repository storage.Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) Create(ctx context.Context, item *itemdetails.ItemDetails) error {
	err := s.repository.Create(ctx, item)
	if err != nil {
		return fmt.Errorf("failed to create ItemDetail: %v", err)
	}
	return nil
}

func (s *Service) FindAllFromOrder(ctx context.Context, order order.Order) ([]itemdetails.ItemDetails, error) {
	itemDetails, err := s.repository.FindAllFromOrder(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("failed to get itemDetails from order with track %s: %v", order.Track_number, err)
	}
	return itemDetails, nil
}
