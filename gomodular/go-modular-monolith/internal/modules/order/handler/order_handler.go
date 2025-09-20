package handler

import (
	"go-modular-monolith/internal/modules/order/service"
	"net/http"

	"github.com/gorilla/mux"
)

type OrderHandler struct {
	orderService service.OrderService
}

func NewOrderHandler(orderService service.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	// Implementation for creating an order
}

func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	// Implementation for retrieving an order
}

func (h *OrderHandler) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	// Implementation for updating an order
}

func (h *OrderHandler) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	// Implementation for deleting an order
}

func (h *OrderHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/orders", h.CreateOrder).Methods(http.MethodPost)
	router.HandleFunc("/orders/{id}", h.GetOrder).Methods(http.MethodGet)
	router.HandleFunc("/orders/{id}", h.UpdateOrder).Methods(http.MethodPut)
	router.HandleFunc("/orders/{id}", h.DeleteOrder).Methods(http.MethodDelete)
}
