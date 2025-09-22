package domain

import (
	"errors"
	"time"

	"go-modular-monolith/pkg/contracts"
)

// Order representa a entidade de domínio do pedido
type Order struct {
	contracts.Order
}

// OrderAggregate contém as regras de negócio do pedido
type OrderAggregate struct {
	order *Order
}

// NewOrder cria um novo pedido com validações de domínio
func NewOrder(id, userID string, items []contracts.OrderItem) (*Order, error) {
	if err := validateUserID(userID); err != nil {
		return nil, err
	}

	if err := validateItems(items); err != nil {
		return nil, err
	}

	total := calculateTotal(items)

	return &Order{
		Order: contracts.Order{
			ID:        id,
			UserID:    userID,
			Items:     items,
			Status:    contracts.OrderStatusPending,
			Total:     total,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}, nil
}

// NewOrderAggregate cria um novo aggregate de pedido
func NewOrderAggregate(order *Order) *OrderAggregate {
	return &OrderAggregate{order: order}
}

// UpdateStatus atualiza o status do pedido com validação
func (oa *OrderAggregate) UpdateStatus(newStatus contracts.OrderStatus) error {
	if err := validateStatusTransition(oa.order.Status, newStatus); err != nil {
		return err
	}

	oa.order.Status = newStatus
	oa.order.UpdatedAt = time.Now()
	return nil
}

// Cancel cancela o pedido se possível
func (oa *OrderAggregate) Cancel() error {
	if oa.order.Status == contracts.OrderStatusDelivered {
		return errors.New("cannot cancel delivered order")
	}

	if oa.order.Status == contracts.OrderStatusCancelled {
		return errors.New("order is already cancelled")
	}

	oa.order.Status = contracts.OrderStatusCancelled
	oa.order.UpdatedAt = time.Now()
	return nil
}

// AddItem adiciona um item ao pedido
func (oa *OrderAggregate) AddItem(item contracts.OrderItem) error {
	if err := validateOrderItem(item); err != nil {
		return err
	}

	if oa.order.Status != contracts.OrderStatusPending {
		return errors.New("can only add items to pending orders")
	}

	oa.order.Items = append(oa.order.Items, item)
	oa.order.Total = calculateTotal(oa.order.Items)
	oa.order.UpdatedAt = time.Now()
	return nil
}

// GetOrder retorna o pedido do aggregate
func (oa *OrderAggregate) GetOrder() *Order {
	return oa.order
}

// IsValid verifica se o pedido é válido
func (oa *OrderAggregate) IsValid() error {
	if oa.order.ID == "" {
		return errors.New("order ID cannot be empty")
	}

	if err := validateUserID(oa.order.UserID); err != nil {
		return err
	}

	if err := validateItems(oa.order.Items); err != nil {
		return err
	}

	if oa.order.Total <= 0 {
		return errors.New("order total must be greater than zero")
	}

	return nil
}

// Domain validation functions

func validateUserID(userID string) error {
	if userID == "" {
		return errors.New("user ID cannot be empty")
	}
	return nil
}

func validateItems(items []contracts.OrderItem) error {
	if len(items) == 0 {
		return errors.New("order must have at least one item")
	}

	for i, item := range items {
		if err := validateOrderItem(item); err != nil {
			return errors.New("invalid item at position " + string(rune(i)) + ": " + err.Error())
		}
	}

	return nil
}

func validateOrderItem(item contracts.OrderItem) error {
	if item.ProductID == "" {
		return errors.New("product ID cannot be empty")
	}

	if item.Quantity <= 0 {
		return errors.New("quantity must be greater than zero")
	}

	if item.Price <= 0 {
		return errors.New("price must be greater than zero")
	}

	return nil
}

func validateStatusTransition(currentStatus, newStatus contracts.OrderStatus) error {
	validTransitions := map[contracts.OrderStatus][]contracts.OrderStatus{
		contracts.OrderStatusPending:   {contracts.OrderStatusConfirmed, contracts.OrderStatusCancelled},
		contracts.OrderStatusConfirmed: {contracts.OrderStatusShipped, contracts.OrderStatusCancelled},
		contracts.OrderStatusShipped:   {contracts.OrderStatusDelivered},
		contracts.OrderStatusDelivered: {}, // Status final
		contracts.OrderStatusCancelled: {}, // Status final
	}

	allowedStatuses, exists := validTransitions[currentStatus]
	if !exists {
		return errors.New("invalid current status")
	}

	for _, status := range allowedStatuses {
		if status == newStatus {
			return nil
		}
	}

	return errors.New("invalid status transition from " + string(currentStatus) + " to " + string(newStatus))
}

func calculateTotal(items []contracts.OrderItem) float64 {
	total := 0.0
	for _, item := range items {
		total += item.Price * float64(item.Quantity)
	}
	return total
}
