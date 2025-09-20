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
```
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

## 🔧 Como Adicionar um Novo Módulo

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

## 🚀 Executando o Projeto

```bash
# Instalar dependências
go mod tidy

# Executar testes
go test ./...

# Executar aplicação
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