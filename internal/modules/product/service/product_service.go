package service

import (
	"context"
	"fmt"
	"time"

	"go-modular-monolith/internal/modules/product/domain"
	"go-modular-monolith/pkg/contracts"

	"github.com/google/uuid"
)

// ProductService implementa a interface contracts.ProductService
type ProductService struct {
	repo           contracts.ProductRepository
	eventPublisher contracts.EventPublisher
}

// NewProductService cria uma nova instância do ProductService
func NewProductService(
	repo contracts.ProductRepository,
	eventPublisher contracts.EventPublisher,
) contracts.ProductService {
	return &ProductService{
		repo:           repo,
		eventPublisher: eventPublisher,
	}
}

// CreateProduct cria um novo produto
func (s *ProductService) CreateProduct(ctx context.Context, req contracts.CreateProductRequest) (*contracts.Product, error) {
	// Criar aggregate de domínio para validações
	aggregate, err := domain.NewProductAggregateFromRequest(req)
	if err != nil {
		return nil, fmt.Errorf("invalid product data: %w", err)
	}

	// Gerar ID único
	productID := uuid.New().String()
	product := aggregate.GetProduct()
	product.ID = productID

	// Salvar no banco de dados
	contractProduct := &product.Product
	if err := s.repo.Create(ctx, contractProduct); err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	// Publicar evento de produto criado
	event := contracts.Event{
		Type:      "ProductCreatedEventType",
		Timestamp: time.Now(),
		Payload: contracts.ProductCreatedEvent{
			ProductID:  productID,
			Name:       product.Name,
			CategoryID: product.CategoryID,
			Price:      product.Price,
		},
	}

	// Ignorar erros de evento para não falhar a operação
	_ = s.eventPublisher.Publish(ctx, event)

	return contractProduct, nil
}

// GetProductByID busca um produto por ID
func (s *ProductService) GetProductByID(ctx context.Context, id string) (*contracts.Product, error) {
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return product, nil
}

// UpdateProduct atualiza um produto existente
func (s *ProductService) UpdateProduct(ctx context.Context, id string, req contracts.UpdateProductRequest) (*contracts.Product, error) {
	// Buscar produto existente
	existingProduct, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	// Criar aggregate para validações e atualizações
	domainProduct := &domain.Product{Product: *existingProduct}
	aggregate := domain.NewProductAggregate(domainProduct)

	// Aplicar atualizações
	if req.Name != nil {
		if err := aggregate.UpdateName(*req.Name); err != nil {
			return nil, fmt.Errorf("invalid name: %w", err)
		}
	}

	if req.Description != nil {
		if err := aggregate.UpdateDescription(*req.Description); err != nil {
			return nil, fmt.Errorf("invalid description: %w", err)
		}
	}

	if req.Price != nil {
		if err := aggregate.UpdatePrice(*req.Price); err != nil {
			return nil, fmt.Errorf("invalid price: %w", err)
		}
	}

	if req.Stock != nil {
		if err := aggregate.UpdateStock(*req.Stock); err != nil {
			return nil, fmt.Errorf("invalid stock: %w", err)
		}
	}

	// Salvar alterações
	updatedProduct := &aggregate.GetProduct().Product
	if err := s.repo.Update(ctx, updatedProduct); err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	return updatedProduct, nil
}

// DeleteProduct remove um produto
func (s *ProductService) DeleteProduct(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	return nil
}

// GetProducts lista produtos com filtros
func (s *ProductService) GetProducts(ctx context.Context, filters contracts.ProductFilters) ([]*contracts.Product, error) {
	products, err := s.repo.List(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %w", err)
	}

	return products, nil
}

// UpdateStock atualiza o estoque de um produto
func (s *ProductService) UpdateStock(ctx context.Context, id string, quantity int) error {
	// Buscar produto existente
	existingProduct, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("product not found: %w", err)
	}

	// Criar aggregate e atualizar estoque
	domainProduct := &domain.Product{Product: *existingProduct}
	aggregate := domain.NewProductAggregate(domainProduct)

	if err := aggregate.UpdateStock(quantity); err != nil {
		return fmt.Errorf("failed to update stock: %w", err)
	}

	// Salvar alterações
	updatedProduct := &aggregate.GetProduct().Product
	if err := s.repo.Update(ctx, updatedProduct); err != nil {
		return fmt.Errorf("failed to update product stock: %w", err)
	}

	// Publicar evento de estoque atualizado
	event := contracts.Event{
		Type:      "ProductStockUpdatedEventType",
		Timestamp: time.Now(),
		Payload: contracts.ProductStockUpdatedEvent{
			ProductID: id,
			NewStock:  quantity,
		},
	}

	// Ignorar erros de evento para não falhar a operação
	_ = s.eventPublisher.Publish(ctx, event)

	return nil
}
