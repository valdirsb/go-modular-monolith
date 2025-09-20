# AI Coding Agent Instructions - Go Modular Monolith

## 🏗️ Architecture Overview

This is a **modular monolith** implementing **Clean Architecture**, **Hexagonal Architecture**, and **DDD** patterns. Each module is self-contained with clear boundaries.

### Key Architectural Principles

- **Dependency Injection**: All services registered in `internal/bootstrap/bootstrap.go` using custom DI container
- **Interface-First**: All contracts defined in `pkg/contracts/interfaces.go` - implement these, never depend on concrete types
- **Event-Driven**: Modules communicate via `pkg/events/eventbus.go` - use `EventPublisher` interface
- **Repository Pattern**: Data access through repository interfaces with adapter implementations

### Module Structure Pattern
```
internal/{module}/
├── domain/          # Business entities and aggregates (pure domain logic)
├── service/         # Application services (orchestrate domain + infrastructure)
├── repository/      # Data persistence adapters
├── handler/         # HTTP handlers (thin layer, delegate to services)
└── adapters/        # External service adapters (email, etc.)
```

## 🚀 Development Workflow

### Essential Commands
```bash
make setup           # Initial project setup (creates .env, starts MySQL)
make docker-up       # Start MySQL/phpMyAdmin containers
make run             # Run application (port 8080)
make api-test        # Run full API integration tests
./scripts/test_api.sh # Manual API testing script
```

### Database & Environment
- **MySQL**: Always use Docker via `make docker-up` (port 3306, root/123456, db: app_db)
- **Config**: Environment variables loaded via `internal/shared/config/config.go`
- **Migrations**: Auto-run via GORM in `internal/shared/database/database.go`
- **phpMyAdmin**: Available at http://localhost:8081 for DB management

### Testing Patterns
- **Integration Tests**: Use `scripts/test_api.sh` - tests full CRUD flow with dynamic user creation
- **API Endpoints**: All routes under `/api/v1/{module}` (e.g., `/api/v1/users/`)
- **Health Check**: Always available at `/health`

## 🔧 Code Patterns & Conventions

### Dependency Injection Pattern
```go
// 1. Register in bootstrap.go
c.RegisterSingleton("serviceName", func() interface{} {
    dep := c.MustGet("dependency").(contracts.Interface)
    return service.NewService(dep)
})

// 2. Use in main.go via DI
handler := container.MustGet("userHandler").(contracts.UserHandler)
```

### Handler Registration Pattern
```go
// In main.go - get handler from DI, register routes
func register{Module}Routes(router *gin.Engine, container *container.Container) {
    handler := container.MustGet("{module}Handler").(contracts.{Module}Handler)
    group := router.Group("/api/v1/{module}s")
    {
        group.POST("/", handler.Create{Module})
        group.GET("/:id", handler.Get{Module})
        // ... standard CRUD pattern
    }
}
```

### Repository Implementation Pattern
```go
// Always implement the contract interface
type MySQL{Module}Repository struct {
    db *gorm.DB
}

func (r *MySQL{Module}Repository) Create(ctx context.Context, entity *contracts.{Module}) error {
    return r.db.WithContext(ctx).Create(entity).Error
}
```

### Service Layer Pattern
```go
func (s *{Module}Service) Create{Module}(ctx context.Context, req contracts.Create{Module}Request) (*contracts.{Module}, error) {
    // 1. Domain validation via aggregates
    aggregate := domain.New{Module}Aggregate(req)
    if err := aggregate.IsValid(); err != nil {
        return nil, err
    }
    
    // 2. Repository interaction
    // 3. Event publishing
    // 4. Return domain entity
}
```

### Event Publishing Pattern
```go
// In services, after successful operations
event := contracts.Event{
    Type:      "{Module}CreatedEventType",
    Timestamp: time.Now(),
    Payload:   contracts.{Module}CreatedEvent{...},
}
s.eventPublisher.Publish(ctx, event)
```

## 🔐 Security & Validation

- **Password Hashing**: Use `adapters.NewArgon2PasswordHasher()` - registered as "passwordHasher" in DI
- **Input Validation**: Gin binding + domain aggregate validation pattern
- **Error Handling**: Return domain errors, map to HTTP status in handlers

## 🎯 Implementation Guidelines

### Adding New Modules
1. **Create module structure** following `internal/user/` pattern
2. **Define contracts** in `pkg/contracts/interfaces.go` first
3. **Register in bootstrap.go** following layered registration (infrastructure → services → handlers)
4. **Add route registration** in `main.go` with DI handler retrieval
5. **Update scripts/test_api.sh** with new endpoint tests

### Database Changes
- **Models**: Define in `contracts/interfaces.go` as structs, implement in domain layer
- **Migrations**: GORM auto-migrate in `database.AutoMigrate()` - add new models there
- **Repositories**: Always implement repository contract, register MySQL adapter in bootstrap

### Critical Files to Understand
- `internal/bootstrap/bootstrap.go` - Complete DI wiring and service registration
- `pkg/contracts/interfaces.go` - All service contracts and data models
- `cmd/server/main.go` - Route registration pattern via DI
- `internal/shared/database/database.go` - Database connection and migration setup
- `pkg/events/eventbus.go` - Inter-module communication mechanism

### Never Modify
- `pkg/container/container.go` - DI container implementation (stable)
- Core interface contracts without updating all implementations
- Database configuration without updating docker-compose.yml

### Common Gotchas
- All routes MUST use `/api/v1/{modules}/` pattern (note trailing slash)
- Repository interfaces use `*contracts.{Model}` pointers, not domain entities
- Event handlers are synchronous - consider async for heavy operations
- Bootstrap registration order matters: infrastructure → domain services → handlers