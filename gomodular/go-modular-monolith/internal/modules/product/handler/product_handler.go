package handler

import (
	"go-modular-monolith/internal/modules/product/service"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productService service.ProductService
}

func NewProductHandler(productService service.ProductService) *ProductHandler {
	return &ProductHandler{productService: productService}
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	// Implementation for creating a product
}

func (h *ProductHandler) GetProduct(c *gin.Context) {
	// Implementation for getting a product by ID
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	// Implementation for updating a product
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	// Implementation for deleting a product
}
