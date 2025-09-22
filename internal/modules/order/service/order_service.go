package service

import (
	"context"
	"errors"
	"time"

	"go-modular-monolith/internal/modules/order/domain"
	"go-modular-monolith/pkg/contracts"
	"go-modular-monolith/pkg/events"

	"github.com/google/uuid"
)

// OrderService implementa a lógica de negócio do módulo de pedidos
type OrderService struct {
	orderRepo      contracts.OrderRepository
	productService contracts.ProductService // Para validar produtos e verificar estoque
	userService    contracts.UserService    // Para validar usuários
	eventPublisher contracts.EventPublisher
}

// NewOrderService cria uma nova instância do serviço de pedidos
func NewOrderService(
	orderRepo contracts.OrderRepository,
	productService contracts.ProductService,
	userService contracts.UserService,
	eventPublisher contracts.EventPublisher,
) contracts.OrderService {
	return &OrderService{
		orderRepo:      orderRepo,
		productService: productService,
		userService:    userService,
		eventPublisher: eventPublisher,
	}
}

// CreateOrder cria um novo pedido
func (s *OrderService) CreateOrder(ctx context.Context, req contracts.CreateOrderRequest) (*contracts.Order, error) {
	// Validar se o usuário existe
	user, err := s.userService.GetUserByID(ctx, req.UserID)
	if err != nil || user == nil {
		return nil, errors.New("invalid user ID")
	}

	// Validar produtos e calcular preços
	orderItems := make([]contracts.OrderItem, len(req.Items))
	for i, item := range req.Items {
		// Buscar produto para validar e obter preço atual
		product, err := s.productService.GetProductByID(ctx, item.ProductID)
		if err != nil || product == nil {
			return nil, errors.New("invalid product ID: " + item.ProductID)
		}

		// Verificar se há estoque suficiente
		if product.Stock < item.Quantity {
			return nil, errors.New("insufficient stock for product: " + product.Name)
		}

		orderItems[i] = contracts.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     product.Price, // Usar o preço atual do produto
		}
	}

	// Gerar ID único para o pedido
	orderID := uuid.New().String()

	// Criar entidade de domínio
	order, err := domain.NewOrder(orderID, req.UserID, orderItems)
	if err != nil {
		return nil, err
	}

	// Criar aggregate e validar
	orderAggregate := domain.NewOrderAggregate(order)
	if err := orderAggregate.IsValid(); err != nil {
		return nil, err
	}

	// Reduzir estoque dos produtos (dentro de transação idealmente)
	for _, item := range orderItems {
		if err := s.productService.UpdateStock(ctx, item.ProductID, -item.Quantity); err != nil {
			return nil, errors.New("failed to update stock for product: " + item.ProductID)
		}
	}

	// Persistir pedido
	orderToSave := orderAggregate.GetOrder()
	if err := s.orderRepo.Create(ctx, &orderToSave.Order); err != nil {
		// Se falhar, reverter estoque (compensação simples)
		for _, item := range orderItems {
			s.productService.UpdateStock(ctx, item.ProductID, item.Quantity) // Reverter
		}
		return nil, errors.New("failed to create order")
	}

	// Publicar evento
	event := contracts.Event{
		Type:      events.OrderCreatedEventType,
		Timestamp: time.Now(),
		Payload: contracts.OrderCreatedEvent{
			OrderID: orderID,
			UserID:  req.UserID,
			Total:   orderToSave.Total,
		},
	}

	if err := s.eventPublisher.Publish(ctx, event); err != nil {
		// Log do erro mas não falhar a operação
	}

	return &orderToSave.Order, nil
}

// GetOrderByID obtém um pedido por ID
func (s *OrderService) GetOrderByID(ctx context.Context, id string) (*contracts.Order, error) {
	if id == "" {
		return nil, errors.New("order ID cannot be empty")
	}

	order, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if order == nil {
		return nil, errors.New("order not found")
	}

	return order, nil
}

// GetOrdersByUserID obtém todos os pedidos de um usuário
func (s *OrderService) GetOrdersByUserID(ctx context.Context, userID string) ([]*contracts.Order, error) {
	if userID == "" {
		return nil, errors.New("user ID cannot be empty")
	}

	// Validar se o usuário existe
	user, err := s.userService.GetUserByID(ctx, userID)
	if err != nil || user == nil {
		return nil, errors.New("invalid user ID")
	}

	orders, err := s.orderRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

// UpdateOrderStatus atualiza o status de um pedido
func (s *OrderService) UpdateOrderStatus(ctx context.Context, id string, status contracts.OrderStatus) error {
	if id == "" {
		return errors.New("order ID cannot be empty")
	}

	// Buscar pedido existente
	existingOrder, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if existingOrder == nil {
		return errors.New("order not found")
	}

	// Criar domain object
	orderDomain := &domain.Order{
		Order: *existingOrder,
	}

	// Criar aggregate e atualizar status
	orderAggregate := domain.NewOrderAggregate(orderDomain)
	if err := orderAggregate.UpdateStatus(status); err != nil {
		return err
	}

	// Persistir alterações
	updatedOrder := orderAggregate.GetOrder()
	if err := s.orderRepo.Update(ctx, &updatedOrder.Order); err != nil {
		return errors.New("failed to update order")
	}

	// Publicar evento
	event := contracts.Event{
		Type:      events.OrderStatusUpdatedEventType,
		Timestamp: time.Now(),
		Payload: map[string]interface{}{
			"order_id":   id,
			"new_status": string(status),
			"old_status": string(existingOrder.Status),
		},
	}

	if err := s.eventPublisher.Publish(ctx, event); err != nil {
		// Log do erro mas não falhar a operação
	}

	return nil
}

// CancelOrder cancela um pedido
func (s *OrderService) CancelOrder(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("order ID cannot be empty")
	}

	// Buscar pedido existente
	existingOrder, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if existingOrder == nil {
		return errors.New("order not found")
	}

	// Criar domain object
	orderDomain := &domain.Order{
		Order: *existingOrder,
	}

	// Criar aggregate e cancelar
	orderAggregate := domain.NewOrderAggregate(orderDomain)
	if err := orderAggregate.Cancel(); err != nil {
		return err
	}

	// Restaurar estoque dos produtos (se o pedido ainda estava pendente/confirmado)
	if existingOrder.Status == contracts.OrderStatusPending || existingOrder.Status == contracts.OrderStatusConfirmed {
		for _, item := range existingOrder.Items {
			if err := s.productService.UpdateStock(ctx, item.ProductID, item.Quantity); err != nil {
				// Log do erro mas continuar o cancelamento
			}
		}
	}

	// Persistir alterações
	cancelledOrder := orderAggregate.GetOrder()
	if err := s.orderRepo.Update(ctx, &cancelledOrder.Order); err != nil {
		return errors.New("failed to cancel order")
	}

	// Publicar evento
	event := contracts.Event{
		Type:      events.OrderCancelledEventType,
		Timestamp: time.Now(),
		Payload: map[string]interface{}{
			"order_id": id,
			"user_id":  existingOrder.UserID,
			"total":    existingOrder.Total,
		},
	}

	if err := s.eventPublisher.Publish(ctx, event); err != nil {
		// Log do erro mas não falhar a operação
	}

	return nil
}
