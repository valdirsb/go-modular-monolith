package events

import (
	"context"
	"fmt"
	"go-modular-monolith/pkg/contracts"
	"sync"
)

// EventBus implementa um sistema de eventos em memória
type EventBus struct {
	handlers map[string][]contracts.EventHandler
	mu       sync.RWMutex
}

// NewEventBus cria uma nova instância do event bus
func NewEventBus() *EventBus {
	return &EventBus{
		handlers: make(map[string][]contracts.EventHandler),
	}
}

// Publish publica um evento para todos os handlers registrados
func (e *EventBus) Publish(ctx context.Context, event contracts.Event) error {
	e.mu.RLock()
	handlers, exists := e.handlers[event.Type]
	e.mu.RUnlock()

	if !exists {
		return nil // Não há handlers registrados para este tipo de evento
	}

	// Executa todos os handlers de forma síncrona
	// Em uma implementação real, você pode querer fazer isso de forma assíncrona
	for _, handler := range handlers {
		if err := handler(ctx, event); err != nil {
			// Log do erro mas continue executando outros handlers
			fmt.Printf("Error handling event %s: %v\n", event.Type, err)
		}
	}

	return nil
}

// Subscribe registra um handler para um tipo específico de evento
func (e *EventBus) Subscribe(eventType string, handler contracts.EventHandler) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.handlers[eventType] = append(e.handlers[eventType], handler)
	return nil
}

// Constants para tipos de eventos
const (
	UserCreatedEventType         = "user.created"
	UserUpdatedEventType         = "user.updated"
	UserDeletedEventType         = "user.deleted"
	ProductCreatedEventType      = "product.created"
	ProductUpdatedEventType      = "product.updated"
	ProductStockUpdatedEventType = "product.stock.updated"
	OrderCreatedEventType        = "order.created"
	OrderStatusUpdatedEventType  = "order.status.updated"
	OrderCancelledEventType      = "order.cancelled"
)
