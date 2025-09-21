package repository

import (
	"context"
	"fmt"

	"go-modular-monolith/pkg/contracts"

	"gorm.io/gorm"
)

// MySQLProductRepository implementa ProductRepository usando MySQL
type MySQLProductRepository struct {
	db *gorm.DB
}

// NewMySQLProductRepository cria uma nova instância do repository
func NewMySQLProductRepository(db *gorm.DB) contracts.ProductRepository {
	return &MySQLProductRepository{db: db}
}

// Create cria um novo produto
func (r *MySQLProductRepository) Create(ctx context.Context, product *contracts.Product) error {
	if err := r.db.WithContext(ctx).Create(product).Error; err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}
	return nil
}

// GetByID busca um produto por ID
func (r *MySQLProductRepository) GetByID(ctx context.Context, id string) (*contracts.Product, error) {
	var product contracts.Product
	if err := r.db.WithContext(ctx).First(&product, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("product not found")
		}
		return nil, fmt.Errorf("failed to get product: %w", err)
	}
	return &product, nil
}

// Update atualiza um produto existente
func (r *MySQLProductRepository) Update(ctx context.Context, product *contracts.Product) error {
	if err := r.db.WithContext(ctx).Save(product).Error; err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}
	return nil
}

// Delete remove um produto
func (r *MySQLProductRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&contracts.Product{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete product: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("product not found")
	}
	return nil
}

// List lista produtos com filtros
func (r *MySQLProductRepository) List(ctx context.Context, filters contracts.ProductFilters) ([]*contracts.Product, error) {
	query := r.db.WithContext(ctx).Model(&contracts.Product{})

	// Aplicar filtros
	if filters.CategoryID != nil {
		query = query.Where("category_id = ?", *filters.CategoryID)
	}

	if filters.MinPrice != nil {
		query = query.Where("price >= ?", *filters.MinPrice)
	}

	if filters.MaxPrice != nil {
		query = query.Where("price <= ?", *filters.MaxPrice)
	}

	if filters.Name != nil {
		query = query.Where("name ILIKE ?", "%"+*filters.Name+"%")
	}

	// Aplicar paginação
	if filters.Limit > 0 {
		query = query.Limit(filters.Limit)
	}

	if filters.Offset > 0 {
		query = query.Offset(filters.Offset)
	}

	var products []*contracts.Product
	if err := query.Find(&products).Error; err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}

	return products, nil
}
