# Monólito Modular em Go

Este projeto implementa um **monólito modular** em Go seguindo as melhores práticas de arquitetura de software. Cada módulo é independente, com baixo acoplamento e alta coesão, permitindo evolução e manutenção facilitadas.

## 🎯 Características Principais

- **Arquitetura Hexagonal**: Separação clara entre domínio, aplicação e infraestrutura
- **Clean Architecture**: Dependências sempre apontando para o domínio
- **Domain-Driven Design**: Modelagem rica do domínio com agregados e entidades
- **Injeção de Dependência**: Container DI para desacoplamento
- **Comunicação por Eventos**: Módulos se comunicam via eventos assíncronos
- **Interfaces Bem Definidas**: Contratos claros entre módulos

## 🏗️ Módulos Implementados

- **User**: Gerenciamento de usuários com autenticação e persistência MySQL
- **Product**: Catálogo de produtos com controle de estoque 
- **Order**: Processamento de pedidos com estados

## 📖 Documentação

- [ARCHITECTURE.md](./ARCHITECTURE.md) - Guia completo da arquitetura e padrões
- [DATABASE.md](./docs/DATABASE.md) - Configuração e uso do MySQL

## Project Structure

```
go-modular-monolith
├── cmd
│   └── server
│       └── main.go
├── internal
│   ├── shared
│   │   ├── config
│   │   │   └── config.go
│   │   ├── database
│   │   │   └── database.go
│   │   ├── middleware
│   │   │   └── middleware.go
│   │   └── logger
│   │       └── logger.go
│   ├── user
│   │   ├── domain
│   │   │   ├── user.go
│   │   │   └── repository.go
│   │   ├── repository
│   │   │   └── user_repository.go
│   │   ├── service
│   │   │   └── user_service.go
│   │   └── handler
│   │       └── user_handler.go
│   ├── order
│   │   ├── domain
│   │   │   ├── order.go
│   │   │   └── repository.go
│   │   ├── repository
│   │   │   └── order_repository.go
│   │   ├── service
│   │   │   └── order_service.go
│   │   └── handler
│   │       └── order_handler.go
│   └── product
│       ├── domain
│       │   ├── product.go
│       │   └── repository.go
│       ├── repository
│       │   └── product_repository.go
│       ├── service
│       │   └── product_service.go
│       └── handler
│           └── product_handler.go
├── pkg
│   └── contracts
│       └── interfaces.go
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
make run          # Executar aplicação
make test         # Executar testes
make api-test     # Testar API completa
make db-shell     # Conectar ao MySQL
```

### Testando a API

```bash
# Testar health check
curl http://localhost:8080/health

# Criar usuário
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"username":"teste","email":"teste@example.com","password":"123456"}'

# Executar script completo de testes
./scripts/test_api.sh
```

### API Endpoints

The application exposes several API endpoints for managing users, orders, and products. Refer to the respective handler files for detailed information on available routes and their functionalities.

### Contributing

Contributions are welcome! Please open an issue or submit a pull request for any enhancements or bug fixes.

### License

This project is licensed under the MIT License. See the LICENSE file for more details.