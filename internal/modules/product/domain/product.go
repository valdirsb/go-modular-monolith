package domain

import (
	"errors"
	"time"
	"unicode/utf8"

	"go-modular-monolith/pkg/contracts"
)

// Product representa a entidade de domínio do produto
type Product struct {
	contracts.Product
}

// ProductAggregate contém as regras de negócio do produto
type ProductAggregate struct {
	product *Product
}

// NewProduct cria um novo produto com validações de domínio
func NewProduct(id, name, description, categoryID string, price float64, stock int) (*Product, error) {
	if err := validateName(name); err != nil {
		return nil, err
	}

	if err := validateDescription(description); err != nil {
		return nil, err
	}

	if err := validatePrice(price); err != nil {
		return nil, err
	}

	if err := validateStock(stock); err != nil {
		return nil, err
	}

	if err := validateCategoryID(categoryID); err != nil {
		return nil, err
	}

	return &Product{
		Product: contracts.Product{
			ID:          id,
			Name:        name,
			Description: description,
			Price:       price,
			Stock:       stock,
			CategoryID:  categoryID,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}, nil
}

// NewProductAggregate cria um novo aggregate de produto
func NewProductAggregate(product *Product) *ProductAggregate {
	return &ProductAggregate{product: product}
}

// NewProductAggregateFromRequest cria um aggregate a partir de um request
func NewProductAggregateFromRequest(req contracts.CreateProductRequest) (*ProductAggregate, error) {
	if err := validateName(req.Name); err != nil {
		return nil, err
	}

	if err := validateDescription(req.Description); err != nil {
		return nil, err
	}

	if err := validatePrice(req.Price); err != nil {
		return nil, err
	}

	if err := validateStock(req.Stock); err != nil {
		return nil, err
	}

	if err := validateCategoryID(req.CategoryID); err != nil {
		return nil, err
	}

	product := &Product{
		Product: contracts.Product{
			Name:        req.Name,
			Description: req.Description,
			Price:       req.Price,
			Stock:       req.Stock,
			CategoryID:  req.CategoryID,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	return &ProductAggregate{product: product}, nil
}

// UpdateName atualiza o nome do produto com validação
func (pa *ProductAggregate) UpdateName(newName string) error {
	if err := validateName(newName); err != nil {
		return err
	}

	pa.product.Name = newName
	pa.product.UpdatedAt = time.Now()
	return nil
}

// UpdateDescription atualiza a descrição do produto
func (pa *ProductAggregate) UpdateDescription(newDescription string) error {
	if err := validateDescription(newDescription); err != nil {
		return err
	}

	pa.product.Description = newDescription
	pa.product.UpdatedAt = time.Now()
	return nil
}

// UpdatePrice atualiza o preço do produto com validação
func (pa *ProductAggregate) UpdatePrice(newPrice float64) error {
	if err := validatePrice(newPrice); err != nil {
		return err
	}

	pa.product.Price = newPrice
	pa.product.UpdatedAt = time.Now()
	return nil
}

// UpdateStock atualiza o estoque do produto
func (pa *ProductAggregate) UpdateStock(newStock int) error {
	if err := validateStock(newStock); err != nil {
		return err
	}

	pa.product.Stock = newStock
	pa.product.UpdatedAt = time.Now()
	return nil
}

// AddStock adiciona itens ao estoque
func (pa *ProductAggregate) AddStock(quantity int) error {
	if quantity < 0 {
		return errors.New("quantity must be positive")
	}

	pa.product.Stock += quantity
	pa.product.UpdatedAt = time.Now()
	return nil
}

// RemoveStock remove itens do estoque
func (pa *ProductAggregate) RemoveStock(quantity int) error {
	if quantity < 0 {
		return errors.New("quantity must be positive")
	}

	if pa.product.Stock < quantity {
		return errors.New("insufficient stock")
	}

	pa.product.Stock -= quantity
	pa.product.UpdatedAt = time.Now()
	return nil
}

// IsInStock verifica se o produto tem estoque suficiente
func (pa *ProductAggregate) IsInStock(quantity int) bool {
	return pa.product.Stock >= quantity
}

// IsValid verifica se o produto é válido
func (pa *ProductAggregate) IsValid() error {
	if err := validateName(pa.product.Name); err != nil {
		return err
	}

	if err := validateDescription(pa.product.Description); err != nil {
		return err
	}

	if err := validatePrice(pa.product.Price); err != nil {
		return err
	}

	if err := validateStock(pa.product.Stock); err != nil {
		return err
	}

	if err := validateCategoryID(pa.product.CategoryID); err != nil {
		return err
	}

	return nil
}

// GetProduct retorna a entidade do produto
func (pa *ProductAggregate) GetProduct() *Product {
	return pa.product
}

// Validation functions

func validateName(name string) error {
	if name == "" {
		return errors.New("product name is required")
	}

	if utf8.RuneCountInString(name) > 100 {
		return errors.New("product name must be at most 100 characters")
	}

	return nil
}

func validateDescription(description string) error {
	if utf8.RuneCountInString(description) > 500 {
		return errors.New("product description must be at most 500 characters")
	}

	return nil
}

func validatePrice(price float64) error {
	if price <= 0 {
		return errors.New("product price must be greater than zero")
	}

	return nil
}

func validateStock(stock int) error {
	if stock < 0 {
		return errors.New("product stock cannot be negative")
	}

	return nil
}

func validateCategoryID(categoryID string) error {
	if categoryID == "" {
		return errors.New("category ID is required")
	}

	return nil
}
