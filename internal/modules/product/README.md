# Product Module 📦

O módulo Product implementa um sistema completo de catálogo de produtos com controle de estoque, seguindo Clean Architecture e DDD.

## 🎯 Responsabilidades

- **Catálogo de Produtos**: CRUD completo com informações detalhadas
- **Controle de Estoque**: Gestão de quantidades disponíveis
- **Sistema de Categorias**: Organização por categorias
- **Filtros Avançados**: Busca por categoria, preço, nome
- **Paginação**: Listagem eficiente de grandes volumes
- **Eventos de Domínio**: Notificações de criação e atualização de estoque

## 🏗️ Estrutura do Módulo

```
product/
├── domain/
│   ├── product.go        # Entidade Product e ProductAggregate
│   └── repository.go     # Interface ProductRepository (Port)
├── service/
│   └── product_service.go # Casos de uso e lógica de aplicação
├── repository/
│   └── product_repository.go # Implementação MySQL
└── handler/
    └── product_handler.go # HTTP handlers com filtros
```

## 📡 API Endpoints

### Base URL: `/api/v1/products/`

#### Criar Produto
```http
POST /api/v1/products/
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
  "created_at": "2025-09-23T20:33:20Z",
  "updated_at": "2025-09-23T20:33:20Z"
}
```

#### Listar Produtos (com Filtros)
```http
GET /api/v1/products/
```

**Query Parameters:**
- `category_id` (string): Filtrar por categoria
- `min_price` (float): Preço mínimo
- `max_price` (float): Preço máximo
- `name` (string): Busca por nome (LIKE)
- `limit` (int): Quantidade máxima de resultados
- `offset` (int): Pular N resultados (paginação)

**Exemplos:**
```bash
# Por categoria
GET /api/v1/products/?category_id=electronics

# Por faixa de preço
GET /api/v1/products/?min_price=2000&max_price=5000

# Busca por nome
GET /api/v1/products/?name=iphone

# Paginação
GET /api/v1/products/?limit=10&offset=20
```

#### Buscar Produto
```http
GET /api/v1/products/{id}
```

#### Atualizar Produto
```http
PUT /api/v1/products/{id}
Content-Type: application/json

{
  "name": "iPhone 15 Pro Max Updated",
  "price": 7999.99,
  "stock": 20
}
```

#### Deletar Produto
```http
DELETE /api/v1/products/{id}
```

#### Atualizar Estoque
```http
PUT /api/v1/products/{id}/stock
Content-Type: application/json

{
  "stock": 25
}
```

## 📊 Dados Seedados

A aplicação inicia com 12 produtos pré-carregados organizados em 7 categorias:

### Categorias Disponíveis
- **electronics**: Smartphones de alta qualidade
- **computers**: Laptops e computadores
- **accessories**: Fones e acessórios de áudio
- **tablets**: Tablets e dispositivos móveis
- **gaming**: Consoles e jogos
- **tv**: TVs e displays
- **wearables**: Smartwatches e wearables

### Produtos Exemplo
```json
[
  {
    "id": "prod-001",
    "name": "iPhone 15 Pro Max",
    "description": "Apple iPhone 15 Pro Max 256GB - Titânio Natural",
    "price": 8999.99,
    "stock": 15,
    "category_id": "electronics"
  },
  {
    "id": "prod-002",
    "name": "MacBook Air M2", 
    "description": "MacBook Air 13\" com chip M2, 8GB RAM, 256GB SSD",
    "price": 12999.99,
    "stock": 8,
    "category_id": "computers"
  }
]
```

## 🔄 Eventos Publicados

### ProductCreatedEvent
Disparado quando um novo produto é criado.

```json
{
  "type": "ProductCreatedEventType",
  "payload": {
    "product_id": "uuid",
    "name": "iPhone 15 Pro Max",
    "category_id": "electronics",
    "price": 8999.99
  },
  "timestamp": "2025-09-23T20:33:20Z"
}
```

### ProductStockUpdatedEvent
Disparado quando o estoque de um produto é atualizado.

```json
{
  "type": "ProductStockUpdatedEventType",
  "payload": {
    "product_id": "uuid",
    "new_stock": 25
  },
  "timestamp": "2025-09-23T20:33:20Z"
}
```

## 🏛️ Arquitetura

### Camada de Domínio
- **Product**: Entidade com validações de negócio
- **ProductAggregate**: Agregado com regras complexas
- **ProductRepository**: Interface para persistência

### Camada de Aplicação
- **ProductService**: Orquestra casos de uso
  - CRUD de produtos
  - Gestão de estoque
  - Aplicação de filtros
  - Publicação de eventos

### Camada de Infraestrutura
- **ProductRepository**: Persistência MySQL com filtros otimizados
- **ProductHandler**: Endpoints HTTP com query parameters

## 🗄️ Modelo de Dados

### Tabela: `products`
```sql
CREATE TABLE products (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10,2) NOT NULL,
    stock INT DEFAULT 0,
    category_id VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_products_category (category_id),
    INDEX idx_products_price (price),
    INDEX idx_products_name (name)
);
```

### Índices para Performance
- **category_id**: Filtros por categoria
- **price**: Consultas por faixa de preço
- **name**: Busca textual por nome

## 🧪 Testes

O módulo Product é testado extensivamente no `test_api.sh`:

1. **Listar produtos seedados**: Verificação dos 12 produtos iniciais
2. **Buscar produto específico**: iPhone (prod-001)
3. **Filtros por categoria**: electronics, computers, etc.
4. **Filtros por preço**: min_price e max_price
5. **Criação de produto**: Produto de teste completo
6. **Atualização de estoque**: Novo valor de estoque
7. **Atualização de produto**: Nome e preço
8. **Exclusão**: Remoção do produto criado

### Filtros Testados
```bash
# Categoria específica
curl "http://localhost:8080/api/v1/products/?category_id=electronics"

# Faixa de preço R$ 2.000 - R$ 5.000
curl "http://localhost:8080/api/v1/products/?min_price=2000&max_price=5000"

# Busca por nome
curl "http://localhost:8080/api/v1/products/?name=iphone"
```

## ⚡ Performance

### Otimizações Implementadas
- **Índices Estratégicos**: category_id, price, name
- **Query Otimizada**: Filtros combinados eficientes
- **Paginação**: Evita carregamento de grandes volumes
- **Seeds Batch**: Inserção em lote dos dados iniciais

### Métricas Esperadas
- **Listagem com filtros**: < 50ms
- **Busca por ID**: < 10ms
- **Criação de produto**: < 30ms
- **Atualização de estoque**: < 20ms

## 🔗 Integração com Outros Módulos

### Order Module
- **Validação de Estoque**: Verificação automática na criação de pedidos
- **Reserva de Produtos**: Desconto de estoque em pedidos confirmados
- **Reversão de Estoque**: Aumento automático em cancelamentos

### User Module
- **Validação de Usuário**: Orders verificam se o usuário existe

## 🚨 Regras de Negócio

### Validações
- **Nome**: Obrigatório, 1-100 caracteres
- **Preço**: Maior que zero
- **Estoque**: Não pode ser negativo
- **Categoria**: Obrigatória

### Estoque
- **Atualização Direta**: Via endpoint específico `/stock`
- **Atualização Automática**: Via módulo Order
- **Validação**: Não permite valores negativos

## 🚀 Melhorias Futuras

### Funcionalidades Planejadas
- [ ] **Sistema de Reviews**: Avaliações e comentários
- [ ] **Imagens de Produto**: Upload e gestão de imagens
- [ ] **Variações**: Cores, tamanhos, versões
- [ ] **Promoções**: Descontos e preços especiais
- [ ] **Categorias Hierárquicas**: Sub-categorias
- [ ] **Busca Full-Text**: Elasticsearch para busca avançada
- [ ] **Cache de Produtos**: Redis para produtos populares
- [ ] **Relatórios**: Analytics de vendas e estoque

### Otimizações Técnicas
- [ ] **Cache de Categorias**: Evitar consultas repetidas
- [ ] **Compressão de Imagens**: Otimização automática
- [ ] **CDN**: Distribuição de conteúdo estático
- [ ] **Search Suggestions**: Autocompletar busca