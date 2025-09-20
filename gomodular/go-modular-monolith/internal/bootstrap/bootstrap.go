package bootstrap

import (
	"log"

	"go-modular-monolith/internal/shared/database"
	"go-modular-monolith/internal/user/adapters"
	"go-modular-monolith/internal/user/handler"
	"go-modular-monolith/internal/user/repository"
	"go-modular-monolith/internal/user/service"
	"go-modular-monolith/pkg/container"
	"go-modular-monolith/pkg/contracts"
	"go-modular-monolith/pkg/events"

	"gorm.io/gorm"
)

// Bootstrap configura toda a aplicação com injeção de dependência
func Bootstrap() (*container.Container, error) {
	c := container.NewContainer()

	// Registrar infraestrutura
	registerInfrastructure(c)

	// Registrar serviços de domínio
	registerDomainServices(c)

	// Registrar handlers HTTP
	registerHandlers(c)

	return c, nil
}

func registerInfrastructure(c *container.Container) {
	// Database Connection
	c.RegisterSingleton("database", func() interface{} {
		config := database.GetDefaultConfig()
		db, err := database.Connect(config)
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}

		// Executar migrações
		if err := database.AutoMigrate(db); err != nil {
			log.Fatalf("Failed to run database migrations: %v", err)
		}

		return db
	})

	// Event Bus
	c.RegisterSingleton("eventbus", func() interface{} {
		return events.NewEventBus()
	})

	// Logger (implementação simples)
	c.RegisterSingleton("logger", func() interface{} {
		return &SimpleLogger{}
	})

	// Password Hasher
	c.RegisterSingleton("passwordHasher", func() interface{} {
		return adapters.NewArgon2PasswordHasher()
	})

	// Email Service (implementação mock)
	c.RegisterSingleton("emailService", func() interface{} {
		return &MockEmailService{}
	})

	// Token Generator (implementação mock)
	c.RegisterSingleton("tokenGenerator", func() interface{} {
		return &MockTokenGenerator{}
	})

	// User Repository (implementação MySQL)
	c.RegisterSingleton("userRepository", func() interface{} {
		db := c.MustGet("database").(*gorm.DB)
		return repository.NewMySQLUserRepository(db)
	})
}

func registerDomainServices(c *container.Container) {
	// User Service
	c.RegisterSingleton("userService", func() interface{} {
		userRepo := c.MustGet("userRepository").(contracts.UserRepository)
		passwordHasher := c.MustGet("passwordHasher").(contracts.PasswordHasher)
		emailService := c.MustGet("emailService").(contracts.EmailService)
		tokenGenerator := c.MustGet("tokenGenerator").(contracts.TokenGenerator)
		eventPublisher := c.MustGet("eventbus").(contracts.EventPublisher)
		logger := c.MustGet("logger").(contracts.Logger)

		return service.NewUserService(
			userRepo,
			passwordHasher,
			emailService,
			tokenGenerator,
			eventPublisher,
			logger,
		)
	})
}

func registerHandlers(c *container.Container) {
	// User Handler
	c.RegisterSingleton("userHandler", func() interface{} {
		userService := c.MustGet("userService").(contracts.UserService)
		return handler.NewUserHandler(userService)
	})

	// Futuros handlers dos outros módulos serão registrados aqui
	// Product Handler, Order Handler, etc.
}
