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

// ProductGRPCService implements the gRPC ProductService
type ProductGRPCService struct {
	pb.UnimplementedProductServiceServer
	productService contracts.ProductService
}

// NewProductGRPCService creates a new gRPC product service
func NewProductGRPCService(productService contracts.ProductService) *ProductGRPCService {
	return &ProductGRPCService{
		productService: productService,
	}
}

// RegisterWithServer registers the service with the gRPC server
func (s *ProductGRPCService) RegisterWithServer(server *grpc.Server) {
	pb.RegisterProductServiceServer(server, s)
}

// CreateProduct creates a new product
func (s *ProductGRPCService) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	// Validate input
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}
	if req.Price <= 0 {
		return nil, status.Error(codes.InvalidArgument, "price must be greater than 0")
	}

	// Create product request
	createReq := contracts.CreateProductRequest{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       int(req.Quantity),
		CategoryID:  "default", // Valor padrão para categoria
	}

	// Call service
	product, err := s.productService.CreateProduct(ctx, createReq)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to create product: %v", err))
	}

	// Convert to proto message
	protoProduct := &pb.Product{
		Id:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Quantity:    int32(product.Stock),
		CreatedAt:   product.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   product.UpdatedAt.Format(time.RFC3339),
	}

	return &pb.CreateProductResponse{
		Product: protoProduct,
		Message: "Product created successfully",
	}, nil
}

// GetProduct gets a product by ID
func (s *ProductGRPCService) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	product, err := s.productService.GetProductByID(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("product not found: %v", err))
	}

	protoProduct := &pb.Product{
		Id:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Quantity:    int32(product.Stock),
		CreatedAt:   product.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   product.UpdatedAt.Format(time.RFC3339),
	}

	return &pb.GetProductResponse{
		Product: protoProduct,
		Message: "Product retrieved successfully",
	}, nil
}

// UpdateProduct updates a product
func (s *ProductGRPCService) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.UpdateProductResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	name := req.Name
	description := req.Description
	price := req.Price
	stock := int(req.Quantity)

	updateReq := contracts.UpdateProductRequest{
		Name:        &name,
		Description: &description,
		Price:       &price,
		Stock:       &stock,
	}

	product, err := s.productService.UpdateProduct(ctx, req.Id, updateReq)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to update product: %v", err))
	}

	protoProduct := &pb.Product{
		Id:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Quantity:    int32(product.Stock),
		CreatedAt:   product.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   product.UpdatedAt.Format(time.RFC3339),
	}

	return &pb.UpdateProductResponse{
		Product: protoProduct,
		Message: "Product updated successfully",
	}, nil
}

// DeleteProduct deletes a product
func (s *ProductGRPCService) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	err := s.productService.DeleteProduct(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to delete product: %v", err))
	}

	return &pb.DeleteProductResponse{
		Message: "Product deleted successfully",
	}, nil
}

// ListProducts lists products with pagination
func (s *ProductGRPCService) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	page := int(req.Page)
	if page <= 0 {
		page = 1
	}

	pageSize := int(req.PageSize)
	if pageSize <= 0 {
		pageSize = 10
	}

	filters := contracts.ProductFilters{
		Limit:  pageSize,
		Offset: (page - 1) * pageSize,
	}

	products, err := s.productService.GetProducts(ctx, filters)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to list products: %v", err))
	}

	var protoProducts []*pb.Product
	for _, product := range products {
		protoProducts = append(protoProducts, &pb.Product{
			Id:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Quantity:    int32(product.Stock),
			CreatedAt:   product.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   product.UpdatedAt.Format(time.RFC3339),
		})
	}

	return &pb.ListProductsResponse{
		Products: protoProducts,
		Total:    int32(len(protoProducts)),
		Page:     int32(page),
		PageSize: int32(pageSize),
		Message:  "Products retrieved successfully",
	}, nil
}

// SearchProducts searches products by query
func (s *ProductGRPCService) SearchProducts(ctx context.Context, req *pb.SearchProductsRequest) (*pb.ListProductsResponse, error) {
	page := int(req.Page)
	if page <= 0 {
		page = 1
	}

	pageSize := int(req.PageSize)
	if pageSize <= 0 {
		pageSize = 10
	}

	filters := contracts.ProductFilters{
		Name:   &req.Query,
		Limit:  pageSize,
		Offset: (page - 1) * pageSize,
	}

	products, err := s.productService.GetProducts(ctx, filters)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to search products: %v", err))
	}

	var protoProducts []*pb.Product
	for _, product := range products {
		protoProducts = append(protoProducts, &pb.Product{
			Id:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Quantity:    int32(product.Stock),
			CreatedAt:   product.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   product.UpdatedAt.Format(time.RFC3339),
		})
	}

	return &pb.ListProductsResponse{
		Products: protoProducts,
		Total:    int32(len(protoProducts)),
		Page:     int32(page),
		PageSize: int32(pageSize),
		Message:  "Products found successfully",
	}, nil
}
