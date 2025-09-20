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

- **User**: Gerenciamento de usuários com autenticação
- **Product**: Catálogo de produtos com controle de estoque
- **Order**: Processamento de pedidos com estados

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

### Running the Application

To start the application, run the following command:
```
go run cmd/server/main.go
```

### API Endpoints

The application exposes several API endpoints for managing users, orders, and products. Refer to the respective handler files for detailed information on available routes and their functionalities.

### Contributing

Contributions are welcome! Please open an issue or submit a pull request for any enhancements or bug fixes.

### License

This project is licensed under the MIT License. See the LICENSE file for more details.