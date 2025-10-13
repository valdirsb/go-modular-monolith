# Monólito Modular em Go - Guia de Arquitetura

Este projeto implementa um **monólito modular** em Go seguindo os melhores padrões de arquitetura de software, incluindo **Clean Architecture**, **Arquitetura Hexagonal** e **Domain-Driven Design (DDD)**.

## 🏗️ Visão Geral da Arquitetura

### Princípios Fundamentais

1. **Baixo Acoplamento**: Módulos independentes que se comunicam apenas através de interfaces bem definidas
2. **Alta Coesão**: Cada módulo tem uma responsabilidade específica e bem definida
3. **Inversão de Dependência**: Dependemos de abstrações, não de implementações
4. **Separação de Responsabilidades**: Domínio, aplicação, infraestrutura claramente separados

### Estrutura do Projeto

```
go-modular-monolith/
├── cmd/                          # Pontos de entrada da aplicação
│   └── server/
│       └── main.go              # Main da aplicação
├── internal/                     # Código interno da aplicação
│   ├── bootstrap/               # Configuração de DI e inicialização
│   │   ├── bootstrap.go
│   │   └── mocks.go
│   ├── shared/                  # Recursos compartilhados
│   │   ├── config/
│   │   ├── database/
│   │   ├── logger/
│   │   └── middleware/
│   └── modules/                 # Módulos de domínio organizados
│       └── {module}/            # Cada módulo (user, product, order)
│           ├── domain/          # Entidades e regras de negócio
│           │   ├── {entity}.go
│           │   └── repository.go # Interface do repositório
│           ├── ports/           # Interfaces (Primary e Secondary Ports)
│           │   └── ports.go
│           ├── service/         # Casos de uso/aplicação
│           │   └── {module}_service.go
│           ├── adapters/        # Implementações de interfaces externas
│           │   └── {adapter}.go
│           ├── repository/      # Implementação de persistência
│           │   └── {module}_repository.go
│           └── handler/         # Controllers/HTTP Handlers
│               └── {module}_handler.go
├── pkg/                         # Código reutilizável
│   ├── contracts/               # Interfaces e contratos globais
│   │   ├── interfaces.go        # Interfaces de domínio
│   │   └── infrastructure.go    # Interfaces de infraestrutura
│   ├── container/               # DI Container
│   │   └── container.go
│   └── events/                  # Sistema de eventos
│       └── eventbus.go
├── go.mod
├── go.sum
└── README.md
```

## 🔄 Padrões Implementados

### 1. Arquitetura Hexagonal (Ports & Adapters)

Cada módulo segue o padrão de arquitetura hexagonal:

- **Domínio** (centro): Entidades, aggregates e regras de negócio
- **Ports** (interfaces): Contratos para comunicação
- **Adapters** (implementações): Implementações concretas dos ports

### 2. Clean Architecture

Camadas bem definidas com dependências sempre apontando para dentro:

```
Frameworks & Drivers (HTTP, DB, External APIs)
    ↓
Interface Adapters (Handlers, Repositories)
    ↓
Application Business Rules (Services)
    ↓
Enterprise Business Rules (Domain)
```

### 3. Domain-Driven Design (DDD)

- **Entities**: Objetos com identidade única
- **Value Objects**: Objetos imutáveis sem identidade
- **Aggregates**: Clusters de entidades tratadas como unidade
- **Domain Services**: Lógica que não pertence a uma entidade específica
- **Repositories**: Interfaces para persistência

## � Módulos Implementados

### User Module (`internal/modules/user/`)
**Responsabilidades:**
- Gerenciamento de usuários (CRUD)
- Autenticação com hash de senha
- Validação de credenciais

**Principais Componentes:**
- `domain/user.go`: Agregado de usuário com validações
- `service/user_service.go`: Casos de uso (criar, buscar, autenticar)
- `repository/mysql_user_repository.go`: Persistência MySQL
- `handler/user_handler.go`: Endpoints HTTP
- `adapters/password_hasher.go`: Hash de senhas com Argon2

**Endpoints:**
- `POST /api/v1/users/` - Criar usuário
- `GET /api/v1/users/:id` - Buscar usuário
- `PUT /api/v1/users/:id` - Atualizar usuário
- `DELETE /api/v1/users/:id` - Remover usuário
- `POST /api/v1/users/validate` - Validar credenciais

### Product Module (`internal/modules/product/`)
**Responsabilidades:**
- Catálogo de produtos com categorias
- Controle de estoque
- Filtros avançados de busca

**Principais Componentes:**
- `domain/product.go`: Entidade produto com validações de negócio
- `service/product_service.go`: Lógica de negócio e estoque
- `repository/product_repository.go`: Persistência com filtros
- `handler/product_handler.go`: API REST com query parameters

**Endpoints:**
- `POST /api/v1/products/` - Criar produto
- `GET /api/v1/products/` - Listar com filtros (categoria, preço, nome)
- `GET /api/v1/products/:id` - Buscar produto específico
- `PUT /api/v1/products/:id` - Atualizar produto
- `DELETE /api/v1/products/:id` - Remover produto
- `PUT /api/v1/products/:id/stock` - Atualizar estoque

### Order Module (`internal/modules/order/`)
**Responsabilidades:**
- Sistema de pedidos com múltiplos itens
- Gestão de status do pedido
- Integração com controle de estoque
- Eventos de negócio (pedido criado, status atualizado)

**Principais Componentes:**
- `domain/order.go`: Agregado de pedido com validações e transições de status
- `service/order_service.go`: Orquestração de criação de pedidos, validações de estoque, cache de produtos
- `repository/mysql_order_repository.go`: Persistência transacional
- `handler/order_handler.go`: Endpoints HTTP RESTful

**Endpoints:**
- `POST /api/v1/orders/` - Criar pedido (com validação de estoque)
- `GET /api/v1/orders/:id` - Buscar pedido específico
- `GET /api/v1/orders/user/:user_id` - Listar pedidos do usuário
- `PUT /api/v1/orders/:id/status` - Atualizar status
- `POST /api/v1/orders/:id/cancel` - Cancelar pedido (com reversão de estoque)

**Recursos Avançados:**
- Cache de produtos para otimização de performance
- Agregação de quantidades por produto
- Transações para garantir consistência
- Reversão automática de estoque em cancelamentos
- Eventos assíncronos para módulos interessados

## �🔧 Como Adicionar um Novo Módulo

### 1. Criar a Estrutura de Diretórios

```bash
mkdir -p internal/{nome_modulo}/{domain,ports,service,adapters,repository,handler}
```

### 2. Definir o Domínio

Crie as entidades e regras de negócio em `domain/`:

```go
// internal/modules/{modulo}/domain/{entidade}.go
package domain

type MinhaEntidade struct {
    contracts.MinhaEntidade
}

type MinhaEntidadeAggregate struct {
    entidade *MinhaEntidade
}

func NewMinhaEntidade(id, nome string) (*MinhaEntidade, error) {
    // Validações de domínio
    if nome == "" {
        return nil, errors.New("nome não pode ser vazio")
    }
    
    return &MinhaEntidade{
        MinhaEntidade: contracts.MinhaEntidade{
            ID: id,
            Nome: nome,
        },
    }, nil
}

// Métodos de negócio...
```

### 3. Definir as Interfaces (Ports)

```go
// internal/modules/{modulo}/ports/ports.go
package ports

import (
    "context"
    "go-modular-monolith/pkg/contracts"
)

// Primary Ports (casos de uso)
type MinhaEntidadeService interface {
    contracts.MinhaEntidadeService
}

// Secondary Ports (adaptadores)
type MinhaEntidadeRepository interface {
    contracts.MinhaEntidadeRepository
}
```

### 4. Implementar o Serviço de Aplicação

```go
// internal/modules/{modulo}/service/{modulo}_service.go
package service

type minhaEntidadeService struct {
    repo           ports.MinhaEntidadeRepository
    eventPublisher contracts.EventPublisher
    logger         contracts.Logger
}

func NewMinhaEntidadeService(
    repo ports.MinhaEntidadeRepository,
    eventPublisher contracts.EventPublisher,
    logger contracts.Logger,
) ports.MinhaEntidadeService {
    return &minhaEntidadeService{
        repo:           repo,
        eventPublisher: eventPublisher,
        logger:         logger,
    }
}

func (s *minhaEntidadeService) Create(ctx context.Context, req contracts.CreateMinhaEntidadeRequest) (*contracts.MinhaEntidade, error) {
    // Implementar lógica de negócio
}
```

### 5. Implementar os Adapters

```go
// internal/modules/{modulo}/repository/{modulo}_repository.go
package repository

type gormMinhaEntidadeRepository struct {
    db *gorm.DB
}

func NewGormMinhaEntidadeRepository(db *gorm.DB) ports.MinhaEntidadeRepository {
    return &gormMinhaEntidadeRepository{db: db}
}

func (r *gormMinhaEntidadeRepository) Create(ctx context.Context, entidade *domain.MinhaEntidade) error {
    // Implementar persistência
}
```

### 6. Registrar no Container DI

```go
// internal/bootstrap/bootstrap.go

func registerDomainServices(c *container.Container) {
    // ... outros serviços
    
    // Meu Módulo Service
    c.RegisterSingleton("minhaEntidadeService", func() interface{} {
        repo := c.MustGet("minhaEntidadeRepository").(contracts.MinhaEntidadeRepository)
        eventPublisher := c.MustGet("eventbus").(contracts.EventPublisher)
        logger := c.MustGet("logger").(contracts.Logger)

        return service.NewMinhaEntidadeService(repo, eventPublisher, logger)
    })
}
```

## 📡 Comunicação entre Módulos

### 1. Através de Interfaces

```go
// Módulo A usa serviço do Módulo B
type ModuloAService struct {
    moduloBService contracts.ModuloBService
}

func (s *ModuloAService) AlgumaOperacao(ctx context.Context) error {
    // Usa ModuloB através de interface
    resultado, err := s.moduloBService.Operacao(ctx, params)
    if err != nil {
        return err
    }
    // ... lógica
}
```

### 2. Através de Eventos

```go
// Módulo A publica evento
event := contracts.Event{
    Type: events.AlgoAconteceuEventType,
    Payload: AlgoAconteceuEvent{
        ID: "123",
        Dados: "importantes",
    },
    Timestamp: time.Now(),
}

s.eventPublisher.Publish(ctx, event)

// Módulo B assina evento
eventBus.Subscribe(events.AlgoAconteceuEventType, func(ctx context.Context, event contracts.Event) error {
    // Processar evento
    return nil
})
```

## 🛠️ Boas Práticas

### 1. Validação de Domínio

- Sempre validar dados no nível de domínio
- Usar aggregates para manter invariantes
- Falhar rapidamente com erros descritivos

### 2. Tratamento de Erros

- Usar erros específicos do domínio
- Não vazar detalhes de implementação
- Logar adequadamente para debugging

### 3. Testes

- Testes unitários para domínio (sem dependências externas)
- Testes de integração para repositórios
- Mocks para dependências externas

### 4. Performance

- Usar context para cancelamento
- Implementar cache quando apropriado
- Considerar lazy loading para aggregates grandes

## � Módulos Implementados

### 👤 Módulo User (`internal/modules/user/`)
- **Domain**: Entidade User com validações de domínio
- **Service**: Criação, autenticação e gerenciamento de usuários
- **Repository**: Persistência MySQL com GORM
- **Handler**: Endpoints REST para CRUD completo
- **Features**: Hash de senhas com Argon2, validação de e-mail

### 📦 Módulo Product (`internal/modules/product/`)
- **Domain**: Aggregate Product com validações de negócio (preço, estoque, nome)
- **Service**: CRUD completo + controle de estoque + eventos
- **Repository**: Persistência MySQL com filtros avançados (categoria, preço, nome)
- **Handler**: Endpoints REST com query parameters para filtros
- **Features**: 
  - Filtros: categoria, faixa de preço, busca por nome, paginação
  - Controle de estoque com validações
  - Publicação de eventos (ProductCreated, StockUpdated)
  - Seeds automáticos com 12 produtos de exemplo

### 🛒 Módulo Order (`internal/modules/order/`)
*Em desenvolvimento* - Estrutura básica criada seguindo os mesmos padrões

### 🔄 Sistema de Events
- Event Bus para comunicação assíncrona entre módulos
- Eventos implementados: UserCreated, ProductCreated, ProductStockUpdated
- Padrão Observer para subscribers de eventos

### 🗄️ Seeds Automáticos
- Sistema de seeds integrado ao processo de migração
- Prevenção automática de duplicação de dados
- Seeds de produtos com dados realistas em 7 categorias
- Execução somente na primeira inicialização

## �🚀 Executando o Projeto

```bash
# Setup completo (recomendado)
make setup

# Instalar dependências
go mod tidy

# Executar testes
go test ./...

# Executar aplicação
make run
# ou
go run cmd/server/main.go
```

## 📚 Recursos Adicionais

- [Clean Architecture - Robert C. Martin](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Hexagonal Architecture - Alistair Cockburn](https://alistair.cockburn.us/hexagonal-architecture/)
- [Domain-Driven Design - Eric Evans](https://www.domainlanguage.com/ddd/)
- [Go Project Layout](https://github.com/golang-standards/project-layout)

## 🤝 Contribuição

1. Siga os padrões de arquitetura estabelecidos
2. Mantenha baixo acoplamento entre módulos
3. Documente interfaces públicas
4. Escreva testes apropriados
5. Use nomes descritivos e consistentes