# API Documentation

Este documento descreve todos os endpoints disponíveis na API do Monólito Modular.

## Base URL
```
http://localhost:8080/api/v1
```

## 🔧 System Endpoints

### Health Check
```http
GET /health
```
Verifica se a aplicação está rodando.

**Response:**
```json
{
  "status": "ok",
  "timestamp": "2025-09-20T20:33:20.123456789Z"
}
```

## 👤 Users Module

### Create User
```http
POST /users/
Content-Type: application/json

{
  "username": "joao123",
  "email": "joao@example.com", 
  "password": "senha123456"
}
```

**Response (201):**
```json
{
  "id": "uuid-generated",
  "username": "joao123",
  "email": "joao@example.com",
  "created_at": "2025-09-20T20:33:20Z",
  "updated_at": "2025-09-20T20:33:20Z"
}
```

### Get User
```http
GET /users/:id
```

### Update User
```http
PUT /users/:id
Content-Type: application/json

{
  "username": "joao_updated",
  "email": "joao_new@example.com"
}
```

### Delete User
```http
DELETE /users/:id
```

### Validate User
```http
POST /users/validate
Content-Type: application/json

{
  "email": "joao@example.com",
  "password": "senha123456"
}
```

## 📦 Products Module

### Create Product
```http
POST /products/
Content-Type: application/json

{
  "name": "iPhone 15 Pro Max",
  "description": "Apple iPhone 15 Pro Max 256GB",
  "price": 8999.99,
  "stock": 15,
  "category_id": "electronics"
}
```

**Response (201):**
```json
{
  "id": "uuid-generated",
  "name": "iPhone 15 Pro Max",
  "description": "Apple iPhone 15 Pro Max 256GB",
  "price": 8999.99,
  "stock": 15,
  "category_id": "electronics",
  "created_at": "2025-09-20T20:33:20Z",
  "updated_at": "2025-09-20T20:33:20Z"
}
```

### Get Product
```http
GET /products/:id
```

### List Products
```http
GET /products/
```

**Query Parameters:**
- `category_id` (string): Filtrar por categoria
- `min_price` (float): Preço mínimo
- `max_price` (float): Preço máximo  
- `name` (string): Busca por nome (LIKE)
- `limit` (int): Quantidade máxima de resultados
- `offset` (int): Pular N resultados (paginação)

**Examples:**
```http
GET /products/?category_id=electronics
GET /products/?min_price=2000&max_price=5000
GET /products/?name=iphone
GET /products/?limit=10&offset=0
```

### Update Product
```http
PUT /products/:id
Content-Type: application/json

{
  "name": "iPhone 15 Pro Max Updated",
  "price": 7999.99,
  "stock": 20
}
```

### Delete Product
```http
DELETE /products/:id
```

### Update Stock
```http
PUT /products/:id/stock
Content-Type: application/json

{
  "stock": 25
}
```

**Response (200):**
```json
{
  "message": "Stock updated successfully"
}
```

## 🛒 Orders Module

### Create Order
```http
POST /orders/
Content-Type: application/json

{
  "user_id": "uuid-of-user",
  "items": [
    {
      "product_id": "prod-001",
      "quantity": 2
    },
    {
      "product_id": "prod-002", 
      "quantity": 1
    }
  ]
}
```

**Response (201):**
```json
{
  "id": "uuid-generated",
  "user_id": "uuid-of-user",
  "items": [
    {
      "product_id": "prod-001",
      "quantity": 2,
      "price": 8999.99
    },
    {
      "product_id": "prod-002",
      "quantity": 1,
      "price": 2499.99
    }
  ],
  "status": "pending",
  "total": 20499.97,
  "created_at": "2025-09-23T20:33:20Z",
  "updated_at": "2025-09-23T20:33:20Z"
}
```

### Get Order
```http
GET /orders/:id
```

**Response (200):**
```json
{
  "id": "order-uuid",
  "user_id": "user-uuid",
  "items": [...],
  "status": "pending",
  "total": 20499.97,
  "created_at": "2025-09-23T20:33:20Z",
  "updated_at": "2025-09-23T20:33:20Z"
}
```

### Get Orders by User
```http
GET /orders/user/:user_id
```

**Response (200):**
```json
[
  {
    "id": "order-uuid-1",
    "user_id": "user-uuid",
    "items": [...],
    "status": "pending",
    "total": 20499.97,
    "created_at": "2025-09-23T20:33:20Z",
    "updated_at": "2025-09-23T20:33:20Z"
  }
]
```

### Update Order Status
```http
PUT /orders/:id/status
Content-Type: application/json

{
  "status": "confirmed"
}
```

**Valid Status Values:**
- `pending` - Pedido pendente
- `confirmed` - Pedido confirmado
- `shipped` - Pedido enviado
- `delivered` - Pedido entregue
- `cancelled` - Pedido cancelado

### Cancel Order
```http
POST /orders/:id/cancel
```

**Response (200):**
```json
{
  "id": "order-uuid",
  "status": "cancelled",
  "message": "Order cancelled successfully"
}
```

## 📊 Seeded Data

A aplicação inicia com 12 produtos pré-carregados:

### Categories Available:
- `electronics` - iPhone 15 Pro Max, Samsung Galaxy S24 Ultra
- `computers` - MacBook Air M2, Dell XPS 13
- `accessories` - AirPods Pro, Sony WH-1000XM5
- `tablets` - iPad Air, Microsoft Surface Pro 9
- `gaming` - Nintendo Switch OLED, PlayStation 5
- `tv` - LG OLED C3 55"
- `wearables` - Apple Watch Series 9

### Example Seeded Products:
```json
[
  {
    "id": "prod-001",
    "name": "iPhone 15 Pro Max",
    "description": "Apple iPhone 15 Pro Max 256GB - Titânio Natural com câmera profissional de 48MP",
    "price": 8999.99,
    "stock": 15,
    "category_id": "electronics"
  },
  {
    "id": "prod-002", 
    "name": "MacBook Air M2",
    "description": "MacBook Air 13\" com chip M2, 8GB RAM, 256GB SSD - Cor Meia-noite",
    "price": 12999.99,
    "stock": 8,
    "category_id": "computers"
  }
]
```

## 🔄 Events Published

### Product Events
- `ProductCreatedEventType` - Quando um produto é criado
- `ProductStockUpdatedEventType` - Quando estoque é atualizado

### User Events  
- `UserCreatedEventType` - Quando um usuário é criado

### Order Events
- `OrderCreatedEventType` - Quando um pedido é criado
- `OrderStatusUpdatedEventType` - Quando status do pedido é atualizado
- `OrderCancelledEventType` - Quando um pedido é cancelado

## ❌ Error Responses

### Validation Error (400)
```json
{
  "error": "validation failed: name is required"
}
```

### Not Found (404)
```json
{
  "error": "product not found"
}
```

### Internal Server Error (500)
```json
{
  "error": "failed to create product: database connection failed"
}
```

## 🧪 Testing

### Automated Test Script
```bash
./scripts/test_api.sh
```

### Manual Testing Examples
```bash
# Health check
curl http://localhost:8080/health

# List products
curl http://localhost:8080/api/v1/products/

# Create user
curl -X POST http://localhost:8080/api/v1/users/ \
  -H "Content-Type: application/json" \
  -d '{"username":"test","email":"test@example.com","password":"123456"}'

# Filter products by category
curl "http://localhost:8080/api/v1/products/?category_id=electronics"

# Create order
curl -X POST http://localhost:8080/api/v1/orders/ \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user-uuid",
    "items": [
      {"product_id": "prod-001", "quantity": 1}
    ]
  }'

# Update order status
curl -X PUT http://localhost:8080/api/v1/orders/order-uuid/status \
  -H "Content-Type: application/json" \
  -d '{"status": "confirmed"}'
```