package contracts

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

// Interfaces para comunicação entre módulos (Ports)

// UserService define operações de negócio relacionadas a usuários
type UserService interface {
	CreateUser(ctx context.Context, req CreateUserRequest) (*User, error)
	GetUserByID(ctx context.Context, id string) (*User, error)
	UpdateUser(ctx context.Context, id string, req UpdateUserRequest) (*User, error)
	DeleteUser(ctx context.Context, id string) error
	ValidateUser(ctx context.Context, email, password string) (*User, error)
}

// ProductService define operações de negócio relacionadas a produtos
type ProductService interface {
	CreateProduct(ctx context.Context, req CreateProductRequest) (*Product, error)
	GetProductByID(ctx context.Context, id string) (*Product, error)
	UpdateProduct(ctx context.Context, id string, req UpdateProductRequest) (*Product, error)
	DeleteProduct(ctx context.Context, id string) error
	GetProducts(ctx context.Context, filters ProductFilters) ([]*Product, error)
	UpdateStock(ctx context.Context, id string, quantity int) error
}

// OrderService define operações de negócio relacionadas a pedidos
type OrderService interface {
	CreateOrder(ctx context.Context, req CreateOrderRequest) (*Order, error)
	GetOrderByID(ctx context.Context, id string) (*Order, error)
	GetOrdersByUserID(ctx context.Context, userID string) ([]*Order, error)
	UpdateOrderStatus(ctx context.Context, id string, status OrderStatus) error
	CancelOrder(ctx context.Context, id string) error
}

// Repository Interfaces (Adapters)

// UserRepository define a interface para persistência de usuários
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id string) error
}

// ProductRepository define a interface para persistência de produtos
type ProductRepository interface {
	Create(ctx context.Context, product *Product) error
	GetByID(ctx context.Context, id string) (*Product, error)
	Update(ctx context.Context, product *Product) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filters ProductFilters) ([]*Product, error)
}

// OrderRepository define a interface para persistência de pedidos
type OrderRepository interface {
	Create(ctx context.Context, order *Order) error
	GetByID(ctx context.Context, id string) (*Order, error)
	GetByUserID(ctx context.Context, userID string) ([]*Order, error)
	Update(ctx context.Context, order *Order) error
	Delete(ctx context.Context, id string) error
}

// Event Publisher para comunicação assíncrona entre módulos
type EventPublisher interface {
	Publish(ctx context.Context, event Event) error
	Subscribe(eventType string, handler EventHandler) error
}

type EventHandler func(ctx context.Context, event Event) error

// Domain Models

// User representa o modelo de domínio do usuário
type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // Não expor na serialização
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Product representa o modelo de domínio do produto
type Product struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	CategoryID  string    `json:"category_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Order representa o modelo de domínio do pedido
type Order struct {
	ID        string      `json:"id"`
	UserID    string      `json:"user_id"`
	Items     []OrderItem `json:"items"`
	Status    OrderStatus `json:"status"`
	Total     float64     `json:"total"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

type OrderItem struct {
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusConfirmed OrderStatus = "confirmed"
	OrderStatusShipped   OrderStatus = "shipped"
	OrderStatusDelivered OrderStatus = "delivered"
	OrderStatusCancelled OrderStatus = "cancelled"
)

// Request/Response DTOs

type CreateUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type UpdateUserRequest struct {
	Username *string `json:"username,omitempty" validate:"omitempty,min=3,max=50"`
	Email    *string `json:"email,omitempty" validate:"omitempty,email"`
}

type CreateProductRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=100"`
	Description string  `json:"description" validate:"max=500"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	Stock       int     `json:"stock" validate:"required,gte=0"`
	CategoryID  string  `json:"category_id" validate:"required"`
}

type UpdateProductRequest struct {
	Name        *string  `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	Description *string  `json:"description,omitempty" validate:"omitempty,max=500"`
	Price       *float64 `json:"price,omitempty" validate:"omitempty,gt=0"`
	Stock       *int     `json:"stock,omitempty" validate:"omitempty,gte=0"`
	CategoryID  *string  `json:"category_id,omitempty"`
}

type CreateOrderRequest struct {
	UserID string            `json:"user_id" validate:"required"`
	Items  []CreateOrderItem `json:"items" validate:"required,dive"`
}

type CreateOrderItem struct {
	ProductID string `json:"product_id" validate:"required"`
	Quantity  int    `json:"quantity" validate:"required,gt=0"`
}

type ProductFilters struct {
	CategoryID *string
	MinPrice   *float64
	MaxPrice   *float64
	Name       *string
	Limit      int
	Offset     int
}

// Events para comunicação entre módulos

type Event struct {
	Type      string      `json:"type"`
	Payload   interface{} `json:"payload"`
	Timestamp time.Time   `json:"timestamp"`
}

type UserCreatedEvent struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
}

type OrderCreatedEvent struct {
	OrderID string  `json:"order_id"`
	UserID  string  `json:"user_id"`
	Total   float64 `json:"total"`
}

type ProductCreatedEvent struct {
	ProductID  string  `json:"product_id"`
	Name       string  `json:"name"`
	CategoryID string  `json:"category_id"`
	Price      float64 `json:"price"`
}

type ProductStockUpdatedEvent struct {
	ProductID string `json:"product_id"`
	NewStock  int    `json:"new_stock"`
}

type PasswordHasher interface {
	Hash(password string) (string, error)
	Verify(password, hash string) bool
}

type EmailService interface {
	SendWelcomeEmail(ctx context.Context, userID, email string) error
	SendPasswordResetEmail(ctx context.Context, userID, email, token string) error
}

type TokenGenerator interface {
	GenerateAccessToken(userID string) (string, error)
	GenerateRefreshToken(userID string) (string, error)
	ValidateToken(token string) (string, error) // retorna userID
}

// Handler Interfaces
type UserHandler interface {
	CreateUser(ctx *gin.Context)
	GetUser(ctx *gin.Context)
	UpdateUser(ctx *gin.Context)
	DeleteUser(ctx *gin.Context)
	ValidateUser(ctx *gin.Context)
}

type ProductHandler interface {
	CreateProduct(ctx *gin.Context)
	GetProduct(ctx *gin.Context)
	UpdateProduct(ctx *gin.Context)
	DeleteProduct(ctx *gin.Context)
	GetProducts(ctx *gin.Context)
	UpdateStock(ctx *gin.Context)
}

type OrderHandler interface {
	CreateOrder(ctx *gin.Context)
	GetOrder(ctx *gin.Context)
	GetOrdersByUser(ctx *gin.Context)
	UpdateOrderStatus(ctx *gin.Context)
	CancelOrder(ctx *gin.Context)
}
