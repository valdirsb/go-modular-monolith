package service

import (
	"errors"
	"go-modular-monolith/internal/modules/order/domain"
)

type OrderService struct {
	repository domain.OrderRepository
}

func NewOrderService(repo domain.OrderRepository) *OrderService {
	return &OrderService{repository: repo}
}

func (s *OrderService) CreateOrder(order *domain.Order) error {
	if order == nil {
		return errors.New("order cannot be nil")
	}
	return s.repository.Save(order)
}

func (s *OrderService) GetOrderByID(id string) (*domain.Order, error) {
	return s.repository.FindByID(id)
}

func (s *OrderService) UpdateOrder(order *domain.Order) error {
	if order == nil {
		return errors.New("order cannot be nil")
	}
	return s.repository.Update(order)
}

func (s *OrderService) DeleteOrder(id string) error {
	return s.repository.Delete(id)
}
