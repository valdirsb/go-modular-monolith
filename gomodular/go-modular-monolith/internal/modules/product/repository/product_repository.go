package repository

import (
	"errors"
	"go-modular-monolith/internal/modules/product/domain"
)

type InMemoryProductRepository struct {
	products map[string]*domain.Product
}

func NewInMemoryProductRepository() *InMemoryProductRepository {
	return &InMemoryProductRepository{
		products: make(map[string]*domain.Product),
	}
}

func (r *InMemoryProductRepository) FindByID(id string) (*domain.Product, error) {
	product, exists := r.products[id]
	if !exists {
		return nil, errors.New("product not found")
	}
	return product, nil
}

func (r *InMemoryProductRepository) Save(product *domain.Product) error {
	r.products[product.ID] = product
	return nil
}

func (r *InMemoryProductRepository) Delete(id string) error {
	if _, exists := r.products[id]; !exists {
		return errors.New("product not found")
	}
	delete(r.products, id)
	return nil
}
