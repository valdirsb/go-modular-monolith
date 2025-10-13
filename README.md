# Monólito Modular em Go

Este projeto implementa um **monólito modular** em Go seguindo as melhores práticas de arquitetura de software. Cada módulo é independente, com baixo acoplamento e alta coesão, permitindo evolução e manutenção facilitadas.

## 🎯 Características Principais

- **Arquitetura Hexagonal**: Separação clara entre domínio, aplicação e infraestrutura
- **Clean Architecture**: Dependências sempre apontando para o domínio
- **Domain-Driven Design**: Modelagem rica do domínio com agregados e entidades
- **Injeção de Dependência**: Container DI para desacoplamento
- **Comunicação por Eventos**: Módulos se comunicam via eventos assíncronos
- **Interfaces Bem Definidas**: Contratos claros entre módulos
- **Dual API Support**: REST (HTTP) e gRPC implementados
- **Framework Modular**: Sistema de providers para protocolos e integrações

## 🏗️ Módulos Implementados

### 📁 Organização Modular
Todos os módulos de domínio estão organizados dentro de `internal/modules/` para melhor organização e escalabilidade:

- **User** (`internal/modules/user/`): Gerenciamento de usuários com autenticação e persistência MySQL
- **Product** (`internal/modules/product/`): Catálogo de produtos com controle de estoque
- **Order** (`internal/modules/order/`): Sistema de pedidos com gestão de estoque e eventos 
- **Order** (`internal/modules/order/`): Processamento de pedidos com estados

### 🔧 Estrutura de Cada Módulo
```
modules/{module}/
├── domain/        # Entidades e regras de negócio
├── ports/         # Interfaces (Primary e Secondary Ports)
├── service/       # Casos de uso/aplicação
├── adapters/      # Implementações de interfaces externas
├── repository/    # Implementação de persistência
└── handler/       # Controllers/HTTP Handlers
```

## 📖 Documentação

## 📚 Documentation

- [Development Guide](docs/DEVELOPMENT.md) - Quick start and development workflow
- [Architecture Guide](ARCHITECTURE.md) - Detailed architectural patterns and conventions
- [API Documentation](docs/API.md) - Complete API endpoints and examples
- [Database Schema](docs/DATABASE.md) - Database structure and relationships
- [Database Migrations](docs/MIGRATIONS.md) - Schema evolution and seeded data

## Project Structure

```
go-modular-monolith
├── cmd
│   └── server
│       └── main.go
├── internal
│   ├── bootstrap                    # Configuração de DI e inicialização
│   │   ├── bootstrap.go
│   │   └── mocks.go
│   ├── shared                      # Recursos compartilhados
│   │   ├── config
│   │   │   └── config.go
│   │   ├── database
│   │   │   └── database.go
│   │   ├── middleware
│   │   │   └── middleware.go
│   │   └── logger
│   │       └── logger.go
│   └── modules                     # Módulos de domínio organizados
│       ├── user
│       │   ├── domain
│       │   │   ├── user.go
│       │   │   └── repository.go
│       │   ├── ports               # Interfaces (Primary e Secondary Ports)
│       │   │   └── ports.go
│       │   ├── adapters           # Implementações de interfaces externas
│       │   │   └── password_hasher.go
│       │   ├── repository
│       │   │   ├── user_repository.go
│       │   │   └── mysql_user_repository.go
│       │   ├── service
│       │   │   └── user_service.go
│       │   └── handler
│       │       └── user_handler.go
│       ├── order
│       │   ├── domain
│       │   │   ├── order.go
│       │   │   └── repository.go
│       │   ├── repository
│       │   │   └── order_repository.go
│       │   ├── service
│       │   │   └── order_service.go
│       │   └── handler
│       │       └── order_handler.go
│       └── product
│           ├── domain
│           │   ├── product.go
│           │   └── repository.go
│           ├── repository
│           │   └── product_repository.go
│           ├── service
│           │   └── product_service.go
│           └── handler
│               └── product_handler.go
├── pkg                            # Código reutilizável
│   ├── container                  # DI Container
│   │   └── container.go
│   ├── contracts                  # Interfaces e contratos globais
│   │   ├── interfaces.go
│   │   └── infrastructure.go
│   └── events                     # Sistema de eventos
│       └── eventbus.go
├── docs                          # Documentação
│   └── DATABASE.md
├── scripts                       # Scripts utilitários
│   ├── init_database.sql
│   └── test_api.sh
├── go.mod
├── go.sum
└── README.md
```

## Getting Started

### Prerequisites

- Go 1.16 or later
- A working database (e.g., PostgreSQL, MySQL)

### Installation

1. Clone the repository:
   ```
   git clone <repository-url>
   cd go-modular-monolith
   ```

2. Install dependencies:
   ```
   go mod tidy
   ```

3. Configure the application by setting environment variables or modifying the configuration files in `internal/shared/config`.

### 🐳 Configuração Rápida com Docker

```bash
# Setup completo (recomendado)
make setup

# Ou passo a passo:
cp .env.example .env
make docker-up    # Inicia MySQL via Docker
make build        # Compila aplicação
```

### ⚙️ Configuração Manual (alternativa)

1. **MySQL com Docker:**
   ```bash
   docker-compose up -d mysql
   ```

2. **MySQL local:**
   ```bash
   sudo apt install mysql-server
   mysql -u root -p
   CREATE DATABASE app_db;
   ```

3. **Variáveis de ambiente (.env):**
   ```
   DB_HOST=localhost
   DB_PORT=3306
   DB_USERNAME=root
   DB_PASSWORD=123456
   DB_DATABASE=app_db
   ```

### 🚀 Executando a Aplicação

```bash
# Modo mais simples
make run

# Ou manualmente:
go mod tidy
go run cmd/server/main.go

# A aplicação irá:
# 1. Conectar ao MySQL (Docker ou local)
# 2. Executar migrações automáticas  
# 3. Iniciar servidor na porta 8080
```

### 📋 Comandos Úteis

```bash
make help         # Ver todos os comandos
make setup        # Configuração inicial completa
make docker-up    # Iniciar MySQL via Docker
make build        # Compilar aplicação
make run          # Executar aplicação (REST apenas)
make run-framework # Executar com framework (REST + gRPC)
make test         # Executar testes
make api-test     # Testar API REST completa
make grpc-test    # Testar serviços gRPC
make proto-gen    # Gerar código a partir dos arquivos .proto
make db-shell     # Conectar ao MySQL
```

## 🔗 APIs Disponíveis

### 🌐 REST API (HTTP)
- **Porta**: 8080 (configurável via `PORT` no .env)
- **Framework**: Gin
- **Documentação**: Endpoints em `/api/v1/{module}`

### ⚡ gRPC API  
- **Porta**: 50051 (configurável via `GRPC_PORT` no .env)
- **Protocolos**: Protocol Buffers
- **Serviços**: UserService, ProductService, OrderService
- **Health Check**: `grpc.health.v1.Health/Check`
- **Reflection**: Habilitada para desenvolvimento
- **Documentação**: [internal/grpc/README.md](internal/grpc/README.md)

### 🧪 Testando a API

#### Health Check
```bash
curl http://localhost:8080/health
```

#### Testando Produtos (com dados seedados)
```bash
# Listar todos os produtos
curl http://localhost:8080/api/v1/products/

# Filtrar por categoria
curl "http://localhost:8080/api/v1/products/?category_id=electronics"

# Filtrar por faixa de preço
curl "http://localhost:8080/api/v1/products/?min_price=2000&max_price=5000"

# Buscar produto específico
curl http://localhost:8080/api/v1/products/prod-001
```

#### Testando Usuários
```bash
# Criar usuário
curl -X POST http://localhost:8080/api/v1/users/ \
  -H "Content-Type: application/json" \
  -d '{"username":"teste","email":"teste@example.com","password":"123456"}'
```

#### Script de Testes Automatizados
```bash
# Executar suite completa de testes
./scripts/test_api.sh
```

### 🌱 Dados Iniciais (Seeds)

A aplicação **automaticamente popula** o banco com 12 produtos de exemplo nas seguintes categorias:
- **Electronics**: iPhone 15 Pro Max, Samsung Galaxy S24 Ultra
- **Computers**: MacBook Air M2, Dell XPS 13  
- **Accessories**: AirPods Pro, Sony WH-1000XM5
- **Tablets**: iPad Air, Microsoft Surface Pro 9
- **Gaming**: Nintendo Switch OLED, PlayStation 5
- **TV**: LG OLED C3 55"
- **Wearables**: Apple Watch Series 9

Os seeds são executados automaticamente na primeira inicialização e não duplicam dados existentes.

### 📡 API Endpoints

#### 👤 Users (`/api/v1/users/`)
- `POST /` - Criar usuário
- `GET /:id` - Buscar usuário por ID
- `PUT /:id` - Atualizar usuário
- `DELETE /:id` - Remover usuário
- `POST /validate` - Validar credenciais

#### 📦 Products (`/api/v1/products/`)
- `POST /` - Criar produto
- `GET /` - Listar produtos (com filtros: `category_id`, `min_price`, `max_price`, `name`, `limit`, `offset`)
- `GET /:id` - Buscar produto por ID
- `PUT /:id` - Atualizar produto
- `DELETE /:id` - Remover produto
- `PUT /:id/stock` - Atualizar estoque

#### 🛒 Orders (`/api/v1/orders/`) 
- `POST /` - Criar pedido
- `GET /:id` - Buscar pedido por ID
- `PUT /:id/status` - Atualizar status do pedido
- `POST /:id/cancel` - Cancelar pedido
- `GET /user/:user_id` - Listar pedidos do usuário

#### 🔧 System
- `GET /health` - Health check da aplicação

### Contributing

Contributions are welcome! Please open an issue or submit a pull request for any enhancements or bug fixes.

### License

This project is licensed under the MIT License. See the LICENSE file for more details.