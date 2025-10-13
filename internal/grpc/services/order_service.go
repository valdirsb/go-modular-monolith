package services

import (
	"context"
	"fmt"
	"time"

	"go-modular-monolith/pkg/contracts"
	pb "go-modular-monolith/pkg/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// OrderGRPCService implements the gRPC OrderService
type OrderGRPCService struct {
	pb.UnimplementedOrderServiceServer
	orderService contracts.OrderService
}

// NewOrderGRPCService creates a new gRPC order service
func NewOrderGRPCService(orderService contracts.OrderService) *OrderGRPCService {
	return &OrderGRPCService{
		orderService: orderService,
	}
}

// RegisterWithServer registers the service with the gRPC server
func (s *OrderGRPCService) RegisterWithServer(server *grpc.Server) {
	pb.RegisterOrderServiceServer(server, s)
}

// CreateOrder creates a new order
func (s *OrderGRPCService) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	// Validate input
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	if req.ProductId == "" {
		return nil, status.Error(codes.InvalidArgument, "product_id is required")
	}
	if req.Quantity <= 0 {
		return nil, status.Error(codes.InvalidArgument, "quantity must be greater than 0")
	}

	// Create order request
	createReq := contracts.CreateOrderRequest{
		UserID: req.UserId,
		Items: []contracts.CreateOrderItem{
			{
				ProductID: req.ProductId,
				Quantity:  int(req.Quantity),
			},
		},
	}

	// Call service
	order, err := s.orderService.CreateOrder(ctx, createReq)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to create order: %v", err))
	}

	// Convert to proto message
	protoOrder := &pb.Order{
		Id:         order.ID,
		UserId:     order.UserID,
		ProductId:  req.ProductId, // Simplificado para um produto
		Quantity:   req.Quantity,
		TotalPrice: order.Total,
		Status:     string(order.Status),
		CreatedAt:  order.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  order.UpdatedAt.Format(time.RFC3339),
	}

	return &pb.CreateOrderResponse{
		Order:   protoOrder,
		Message: "Order created successfully",
	}, nil
}

// GetOrder gets an order by ID
func (s *OrderGRPCService) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	order, err := s.orderService.GetOrderByID(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("order not found: %v", err))
	}

	// Para simplificar, vamos pegar o primeiro item
	var productId string = ""
	var quantity int32 = 0
	if len(order.Items) > 0 {
		productId = order.Items[0].ProductID
		quantity = int32(order.Items[0].Quantity)
	}

	protoOrder := &pb.Order{
		Id:         order.ID,
		UserId:     order.UserID,
		ProductId:  productId,
		Quantity:   quantity,
		TotalPrice: order.Total,
		Status:     string(order.Status),
		CreatedAt:  order.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  order.UpdatedAt.Format(time.RFC3339),
	}

	return &pb.GetOrderResponse{
		Order:   protoOrder,
		Message: "Order retrieved successfully",
	}, nil
}

// UpdateOrderStatus updates an order status
func (s *OrderGRPCService) UpdateOrderStatus(ctx context.Context, req *pb.UpdateOrderStatusRequest) (*pb.UpdateOrderStatusResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	if req.Status == "" {
		return nil, status.Error(codes.InvalidArgument, "status is required")
	}

	err := s.orderService.UpdateOrderStatus(ctx, req.Id, contracts.OrderStatus(req.Status))
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to update order status: %v", err))
	}

	// Get updated order
	order, err := s.orderService.GetOrderByID(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get updated order: %v", err))
	}

	var productId string = ""
	var quantity int32 = 0
	if len(order.Items) > 0 {
		productId = order.Items[0].ProductID
		quantity = int32(order.Items[0].Quantity)
	}

	protoOrder := &pb.Order{
		Id:         order.ID,
		UserId:     order.UserID,
		ProductId:  productId,
		Quantity:   quantity,
		TotalPrice: order.Total,
		Status:     string(order.Status),
		CreatedAt:  order.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  order.UpdatedAt.Format(time.RFC3339),
	}

	return &pb.UpdateOrderStatusResponse{
		Order:   protoOrder,
		Message: "Order status updated successfully",
	}, nil
}

// DeleteOrder deletes an order
func (s *OrderGRPCService) DeleteOrder(ctx context.Context, req *pb.DeleteOrderRequest) (*pb.DeleteOrderResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	err := s.orderService.CancelOrder(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to cancel order: %v", err))
	}

	return &pb.DeleteOrderResponse{
		Message: "Order cancelled successfully",
	}, nil
}

// ListOrders lists orders with pagination
func (s *OrderGRPCService) ListOrders(ctx context.Context, req *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	// Para simplificar, vamos retornar uma lista vazia por enquanto
	return &pb.ListOrdersResponse{
		Orders:   []*pb.Order{},
		Total:    0,
		Page:     req.Page,
		PageSize: req.PageSize,
		Message:  "Orders retrieved successfully",
	}, nil
}

// GetOrdersByUser gets orders by user ID
func (s *OrderGRPCService) GetOrdersByUser(ctx context.Context, req *pb.GetOrdersByUserRequest) (*pb.ListOrdersResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	orders, err := s.orderService.GetOrdersByUserID(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get user orders: %v", err))
	}

	var protoOrders []*pb.Order
	for _, order := range orders {
		var productId string = ""
		var quantity int32 = 0
		if len(order.Items) > 0 {
			productId = order.Items[0].ProductID
			quantity = int32(order.Items[0].Quantity)
		}

		protoOrders = append(protoOrders, &pb.Order{
			Id:         order.ID,
			UserId:     order.UserID,
			ProductId:  productId,
			Quantity:   quantity,
			TotalPrice: order.Total,
			Status:     string(order.Status),
			CreatedAt:  order.CreatedAt.Format(time.RFC3339),
			UpdatedAt:  order.UpdatedAt.Format(time.RFC3339),
		})
	}

	return &pb.ListOrdersResponse{
		Orders:   protoOrders,
		Total:    int32(len(protoOrders)),
		Page:     req.Page,
		PageSize: req.PageSize,
		Message:  "User orders retrieved successfully",
	}, nil
}
