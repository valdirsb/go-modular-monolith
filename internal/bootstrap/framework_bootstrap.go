package bootstrap

import (
	"context"
	"fmt"
	"log"

	"go-modular-monolith/internal/modules/user/adapters"
	userHandler "go-modular-monolith/internal/modules/user/handler"
	userRepository "go-modular-monolith/internal/modules/user/repository"
	userService "go-modular-monolith/internal/modules/user/service"

	productHandler "go-modular-monolith/internal/modules/product/handler"
	productRepository "go-modular-monolith/internal/modules/product/repository"
	productService "go-modular-monolith/internal/modules/product/service"

	orderHandler "go-modular-monolith/internal/modules/order/handler"
	orderRepository "go-modular-monolith/internal/modules/order/repository"
	orderService "go-modular-monolith/internal/modules/order/service"

	"go-modular-monolith/internal/grpc/services"
	"go-modular-monolith/internal/shared/config"
	"go-modular-monolith/internal/shared/database"
	"go-modular-monolith/pkg/container"
	"go-modular-monolith/pkg/contracts"
	"go-modular-monolith/pkg/events"
	"go-modular-monolith/pkg/framework"
	"go-modular-monolith/pkg/framework/providers"
	grpcProvider "go-modular-monolith/pkg/framework/providers/grpc"

	"gorm.io/gorm"
)

// initializeGRPCProvider initializes the gRPC provider
func initializeGRPCProvider(fw *framework.Framework) error {
	// Load app config to get gRPC port
	appConfig, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load app config: %v", err)
	}

	grpcProvider := grpcProvider.NewGRPCProvider(appConfig.GRPCPort)
	fw.RegisterProvider("grpc", grpcProvider)

	// Note: não inicializamos ainda, apenas registramos
	// A inicialização será feita após o registro dos serviços

	log.Printf("✅ gRPC provider registered on port %s", appConfig.GRPCPort)
	return nil
}

// FrameworkBootstrap configures the application with framework integration
func FrameworkBootstrap(configPath string) (*container.Container, *framework.Framework, error) {
	// Create DI container
	c := container.NewContainer()

	// Initialize framework
	fw, err := framework.NewFramework(configPath, c)
	if err != nil {
		return nil, nil, err
	}

	// Register framework in container for access by other services
	c.RegisterSingleton("framework", func() interface{} {
		return fw
	})

	// Initialize framework
	ctx := context.Background()
	if err := fw.Initialize(ctx); err != nil {
		return nil, nil, err
	}

	// Register traditional services based on enabled modules
	registerCoreInfrastructure(c, fw)

	if framework.IsEnabled("modules", "user") {
		registerUserModule(c)
	}

	if framework.IsEnabled("modules", "product") {
		registerProductModule(c)
	}

	if framework.IsEnabled("modules", "order") {
		registerOrderModule(c)
	}

	// Register handlers for enabled modules
	registerModuleHandlers(c, fw)

	// Initialize and register gRPC services if enabled
	if framework.IsEnabled("protocols", "grpc") {
		if err := initializeGRPCProvider(fw); err != nil {
			return nil, nil, err
		}
		registerGRPCServices(c, fw)
	}

	return c, fw, nil
}

// registerCoreInfrastructure registers core infrastructure services
func registerCoreInfrastructure(c *container.Container, fw *framework.Framework) {
	// Database Connection (if enabled)
	if framework.IsEnabled("database", "mysql") {
		c.RegisterSingleton("database", func() interface{} {
			config := database.GetDefaultConfig()
			db, err := database.Connect(config)
			if err != nil {
				log.Fatalf("Failed to connect to database: %v", err)
			}
			return db
		})
	}

	// Event Bus (if enabled)
	if framework.IsEnabled("core", "event_system") {
		c.RegisterSingleton("eventbus", func() interface{} {
			return events.NewEventBus()
		})
	}

	// Password Hasher
	c.RegisterSingleton("passwordHasher", func() interface{} {
		return adapters.NewArgon2PasswordHasher()
	})

	// Mock Services (would be replaced with real implementations)
	c.RegisterSingleton("emailService", func() interface{} {
		return &MockEmailService{}
	})

	c.RegisterSingleton("tokenGenerator", func() interface{} {
		return &MockTokenGenerator{}
	})

	c.RegisterSingleton("logger", func() interface{} {
		return &SimpleLogger{}
	})
}

// registerUserModule registers user module services
func registerUserModule(c *container.Container) {
	// User Repository
	c.RegisterSingleton("userRepository", func() interface{} {
		db := c.MustGet("database").(*gorm.DB)
		return userRepository.NewMySQLUserRepository(db)
	})

	// User Service
	c.RegisterSingleton("userService", func() interface{} {
		userRepo := c.MustGet("userRepository").(contracts.UserRepository)
		passwordHasher := c.MustGet("passwordHasher").(contracts.PasswordHasher)
		emailService := c.MustGet("emailService").(contracts.EmailService)
		tokenGenerator := c.MustGet("tokenGenerator").(contracts.TokenGenerator)
		eventPublisher := c.MustGet("eventbus").(contracts.EventPublisher)
		logger := c.MustGet("logger").(contracts.Logger)

		return userService.NewUserService(
			userRepo,
			passwordHasher,
			emailService,
			tokenGenerator,
			eventPublisher,
			logger,
		)
	})
}

// registerProductModule registers product module services
func registerProductModule(c *container.Container) {
	// Product Repository
	c.RegisterSingleton("productRepository", func() interface{} {
		db := c.MustGet("database").(*gorm.DB)
		return productRepository.NewMySQLProductRepository(db)
	})

	// Product Service
	c.RegisterSingleton("productService", func() interface{} {
		productRepo := c.MustGet("productRepository").(contracts.ProductRepository)
		eventPublisher := c.MustGet("eventbus").(contracts.EventPublisher)

		return productService.NewProductService(
			productRepo,
			eventPublisher,
		)
	})
}

// registerOrderModule registers order module services
func registerOrderModule(c *container.Container) {
	// Order Repository
	c.RegisterSingleton("orderRepository", func() interface{} {
		db := c.MustGet("database").(*gorm.DB)
		return orderRepository.NewMySQLOrderRepository(db)
	})

	// Order Service
	c.RegisterSingleton("orderService", func() interface{} {
		orderRepo := c.MustGet("orderRepository").(contracts.OrderRepository)
		productSvc := c.MustGet("productService").(contracts.ProductService)
		userSvc := c.MustGet("userService").(contracts.UserService)
		eventPublisher := c.MustGet("eventbus").(contracts.EventPublisher)

		return orderService.NewOrderService(
			orderRepo,
			productSvc,
			userSvc,
			eventPublisher,
		)
	})
}

// registerModuleHandlers registers HTTP handlers for enabled modules
func registerModuleHandlers(c *container.Container, fw *framework.Framework) {
	if framework.IsEnabled("modules", "user") {
		c.RegisterSingleton("userHandler", func() interface{} {
			userService := c.MustGet("userService").(contracts.UserService)
			return userHandler.NewUserHandler(userService)
		})
	}

	if framework.IsEnabled("modules", "product") {
		c.RegisterSingleton("productHandler", func() interface{} {
			productSvc := c.MustGet("productService").(contracts.ProductService)
			return productHandler.NewProductHandler(productSvc)
		})
	}

	if framework.IsEnabled("modules", "order") {
		c.RegisterSingleton("orderHandler", func() interface{} {
			orderSvc := c.MustGet("orderService").(contracts.OrderService)
			return orderHandler.NewOrderHandler(orderSvc)
		})
	}
}

// registerGRPCServices registers gRPC services with the gRPC provider
func registerGRPCServices(c *container.Container, fw *framework.Framework) {
	// Get gRPC provider from framework
	grpcProv := fw.GetProvider("grpc")
	if grpcProv == nil {
		log.Printf("gRPC provider not found")
		return
	}

	grpcProvider, ok := grpcProv.(*grpcProvider.GRPCProvider)
	if !ok {
		log.Printf("Invalid gRPC provider type")
		return
	}

	// Register User gRPC Service
	if framework.IsEnabled("modules", "user") {
		userSvc := c.MustGet("userService").(contracts.UserService)
		userRepo := c.MustGet("userRepository").(contracts.UserRepository)
		userGRPCService := services.NewUserGRPCService(userSvc, userRepo)
		grpcProvider.RegisterService(userGRPCService)
		log.Printf("✅ User gRPC service registered")
	}

	// Register Product gRPC Service
	if framework.IsEnabled("modules", "product") {
		productSvc := c.MustGet("productService").(contracts.ProductService)
		productGRPCService := services.NewProductGRPCService(productSvc)
		grpcProvider.RegisterService(productGRPCService)
		log.Printf("✅ Product gRPC service registered")
	}

	// Register Order gRPC Service
	if framework.IsEnabled("modules", "order") {
		orderSvc := c.MustGet("orderService").(contracts.OrderService)
		orderGRPCService := services.NewOrderGRPCService(orderSvc)
		grpcProvider.RegisterService(orderGRPCService)
		log.Printf("✅ Order gRPC service registered")
	}

	// Agora inicializar o provider com todos os serviços registrados
	ctx := context.Background()
	deps := providers.Dependencies{
		Config: fw.GetConfig(),
	}
	if err := grpcProvider.Initialize(ctx, deps); err != nil {
		log.Printf("❌ Failed to initialize gRPC provider: %v", err)
		return
	}

	log.Printf("✅ gRPC server started successfully on port %s", grpcProvider.GetPort())
}
