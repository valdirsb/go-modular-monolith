package service

import (
	"errors"
	"go-modular-monolith/internal/product/domain"
)

type ProductService struct {
	repo domain.ProductRepository
}

func NewProductService(repo domain.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) CreateProduct(product *domain.Product) error {
	if product == nil {
		return errors.New("product cannot be nil")
	}
	return s.repo.Save(product)
}

func (s *ProductService) GetProductByID(id string) (*domain.Product, error) {
	return s.repo.FindByID(id)
}

func (s *ProductService) UpdateProduct(product *domain.Product) error {
	if product == nil {
		return errors.New("product cannot be nil")
	}
	return s.repo.Save(product)
}

func (s *ProductService) DeleteProduct(id string) error {
	return s.repo.Delete(id)
}
