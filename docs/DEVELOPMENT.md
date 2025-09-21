# Development Guide

Este guia ajuda novos desenvolvedores a começar rapidamente no projeto.

## 🚀 Quick Start (5 minutos)

### 1. Pré-requisitos
```bash
# Verificar se tem Go 1.21+
go version

# Verificar se tem Docker
docker --version

# Verificar se tem Make
make --version
```

### 2. Setup Inicial
```bash
# 1. Clonar repositório
git clone <repository-url>
cd monolito_modular

# 2. Setup automático (cria .env e inicia MySQL)
make setup

# 3. Executar aplicação
make run
```

### 3. Validar Setup
```bash
# Testar health check
curl http://localhost:8080/health

# Executar testes automatizados
make api-test

# Acessar phpMyAdmin (opcional)
open http://localhost:8081
```

## 🏗️ Estrutura de Desenvolvimento

### Comandos Essenciais
```bash
make setup           # Setup inicial completo
make docker-up       # Apenas subir MySQL
make run             # Executar aplicação
make api-test        # Testes de API
make docker-down     # Parar containers
```

### Estrutura de Pastas
```
internal/modules/{module}/
├── domain/          # Lógica de negócio pura
├── service/         # Orquestração de domínio
├── repository/      # Acesso a dados
├── handler/         # Controladores HTTP
└── adapters/        # Integrações externas
```

### Fluxo de Dados
```
HTTP Request → Handler → Service → Repository → Database
                  ↓
              Event Bus → Other Modules
```

## 🔄 Workflow de Desenvolvimento

### Implementar Novo Módulo

1. **Definir Contratos** (`pkg/contracts/interfaces.go`)
```go
type {Module}Handler interface {
    Create{Module}(c *gin.Context)
    Get{Module}(c *gin.Context)
    // ...
}

type {Module} struct {
    ID          string `json:"id" gorm:"primaryKey"`
    // campos...
}
```

2. **Criar Estrutura de Pasta**
```bash
mkdir -p internal/modules/{module}/{domain,service,repository,handler}
```

3. **Implementar Camadas** (seguir padrão dos módulos existentes)
   - Domain: Aggregate + validações
   - Repository: Interface + implementação MySQL
   - Service: Lógica de negócio + eventos
   - Handler: Endpoints REST

4. **Registrar no DI** (`internal/bootstrap/bootstrap.go`)
```go
// Repository
c.RegisterSingleton("{module}Repository", func() interface{} {
    db := c.MustGet("db").(*gorm.DB)
    return repository.NewMySQL{Module}Repository(db)
})

// Service  
c.RegisterSingleton("{module}Service", func() interface{} {
    repo := c.MustGet("{module}Repository").(contracts.{Module}Repository)
    eventBus := c.MustGet("eventBus").(contracts.EventPublisher)
    return service.New{Module}Service(repo, eventBus)
})

// Handler
c.RegisterSingleton("{module}Handler", func() interface{} {
    service := c.MustGet("{module}Service").(contracts.{Module}Service)
    return handler.New{Module}Handler(service)
})
```

5. **Adicionar Rotas** (`cmd/server/main.go`)
```go
func register{Module}Routes(router *gin.Engine, container *container.Container) {
    handler := container.MustGet("{module}Handler").(contracts.{Module}Handler)
    group := router.Group("/api/v1/{module}s")
    {
        group.POST("/", handler.Create{Module})
        group.GET("/:id", handler.Get{Module})
        // ...
    }
}

// No main()
register{Module}Routes(router, container)
```

6. **Atualizar Migração** (`internal/shared/database/database.go`)
```go
err = db.AutoMigrate(
    &contracts.User{},
    &contracts.Product{},
    &contracts.{Module}{}, // Adicionar aqui
)
```

7. **Adicionar Testes** (`scripts/test_api.sh`)

### Debugging Comum

#### Database Issues
```bash
# Verificar se MySQL está rodando
docker ps | grep mysql

# Conectar ao MySQL
docker exec -it mysql-container mysql -u root -p123456 app_db

# Ver logs do container
docker logs mysql-container
```

#### API Issues
```bash
# Ver logs da aplicação
tail -f app.log

# Testar endpoint específico
curl -v http://localhost:8080/api/v1/products/

# Verificar porta ocupada
lsof -i :8080
```

#### Build Issues
```bash
# Limpar cache do Go
go clean -cache

# Atualizar dependências
go mod tidy

# Rebuild completo
make clean && make setup
```

## 🎯 Padrões de Código

### Repository Pattern
```go
func (r *MySQL{Module}Repository) Create(ctx context.Context, entity *contracts.{Module}) error {
    return r.db.WithContext(ctx).Create(entity).Error
}

func (r *MySQL{Module}Repository) FindByID(ctx context.Context, id string) (*contracts.{Module}, error) {
    var entity contracts.{Module}
    if err := r.db.WithContext(ctx).First(&entity, "id = ?", id).Error; err != nil {
        return nil, err
    }
    return &entity, nil
}
```

### Service Pattern
```go
func (s *{Module}Service) Create{Module}(ctx context.Context, req contracts.Create{Module}Request) (*contracts.{Module}, error) {
    // 1. Validação via domain
    aggregate := domain.New{Module}Aggregate(req)
    if err := aggregate.IsValid(); err != nil {
        return nil, err
    }

    // 2. Persistência
    entity := aggregate.ToEntity()
    if err := s.repository.Create(ctx, entity); err != nil {
        return nil, err
    }

    // 3. Evento
    s.eventPublisher.Publish(ctx, contracts.Event{
        Type: "{Module}CreatedEventType",
        Payload: contracts.{Module}CreatedEvent{EntityID: entity.ID},
    })

    return entity, nil
}
```

### Handler Pattern
```go
func (h *{Module}Handler) Create{Module}(c *gin.Context) {
    var req contracts.Create{Module}Request
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    entity, err := h.service.Create{Module}(c.Request.Context(), req)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    c.JSON(201, entity)
}
```

## 🧪 Testing Strategy

### Levels de Teste
1. **Unit Tests**: Testes de domínio e serviços
2. **Integration Tests**: Testes de repository com banco
3. **API Tests**: Testes end-to-end via HTTP

### Test Data Management
```bash
# Seeds automáticos na inicialização
# Ver internal/shared/database/database.go - createDefaultProducts()

# Reset completo do banco
docker-compose down -v
make setup
```

### Performance Testing
```bash
# Teste de carga básico
for i in {1..100}; do
  curl -s http://localhost:8080/api/v1/products/ > /dev/null &
done
wait
```

## 📊 Monitoring & Logs

### Application Logs
```bash
# Logs em tempo real
tail -f app.log

# Filtrar por nível
grep "ERROR" app.log
```

### Database Monitoring
```bash
# Via phpMyAdmin
open http://localhost:8081

# Via linha de comando
docker exec mysql-container mysql -u root -p123456 -e "SHOW PROCESSLIST;" app_db
```

## 🚀 Production Deployment

### Environment Variables
```bash
# Copiar e ajustar para produção
cp .env.example .env

# Ajustar:
# - DATABASE_URL para produção
# - JWT_SECRET com chave forte
# - PORT conforme necessário
```

### Build para Produção
```bash
# Build binário
go build -o bin/server cmd/server/main.go

# Executar
./bin/server
```

### Docker Deploy *(Futuro)*
```dockerfile
FROM golang:1.21-alpine AS builder
# ... build steps

FROM alpine:latest
# ... runtime setup
```