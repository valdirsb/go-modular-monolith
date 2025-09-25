package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-modular-monolith/internal/bootstrap"
	"go-modular-monolith/pkg/container"
	"go-modular-monolith/pkg/contracts"

	"github.com/gin-gonic/gin"
)

func main() {
	// Inicializar container de dependências
	container, err := bootstrap.Bootstrap()
	if err != nil {
		log.Fatalf("Failed to bootstrap application: %v", err)
	}

	// Obter logger
	logger := container.MustGet("logger").(contracts.Logger)
	logger.Info("Starting Go Modular Monolith")

	// Configurar Gin
	router := gin.Default()

	// Middleware global
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"timestamp": time.Now(),
		})
	})

	// Registrar rotas dos módulos
	registerUserRoutes(router, container)
	registerProductRoutes(router, container)
	registerOrderRoutes(router, container)

	// Configurar servidor HTTP
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Iniciar servidor em goroutine
	go func() {
		logger.Info("Server starting on port 8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Aguardar sinal de interrupção
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	logger.Info("Server exited")
}

// registerUserRoutes registra as rotas do módulo de usuário
func registerUserRoutes(router *gin.Engine, container *container.Container) {
	userHandler := container.MustGet("userHandler").(contracts.UserHandler)

	userGroup := router.Group("/api/v1/users")
	{
		userGroup.POST("/", userHandler.CreateUser)
		userGroup.GET("/:id", userHandler.GetUser)
		userGroup.PUT("/:id", userHandler.UpdateUser)
		userGroup.DELETE("/:id", userHandler.DeleteUser)
		userGroup.POST("/validate", userHandler.ValidateUser)
	}
}

// registerProductRoutes registra as rotas do módulo de produto
func registerProductRoutes(router *gin.Engine, container *container.Container) {
	productHandler := container.MustGet("productHandler").(contracts.ProductHandler)

	productGroup := router.Group("/api/v1/products")
	{
		productGroup.POST("/", productHandler.CreateProduct)
		productGroup.GET("/", productHandler.GetProducts)
		productGroup.GET("/:id", productHandler.GetProduct)
		productGroup.PUT("/:id", productHandler.UpdateProduct)
		productGroup.DELETE("/:id", productHandler.DeleteProduct)
		productGroup.PUT("/:id/stock", productHandler.UpdateStock)
	}
}

// registerOrderRoutes registra as rotas do módulo de pedidos
func registerOrderRoutes(router *gin.Engine, container *container.Container) {
	orderHandler := container.MustGet("orderHandler").(contracts.OrderHandler)

	orderGroup := router.Group("/api/v1/orders")
	{
		orderGroup.POST("/", orderHandler.CreateOrder)
		orderGroup.GET("/:id", orderHandler.GetOrder)
		orderGroup.PUT("/:id/status", orderHandler.UpdateOrderStatus)
		orderGroup.POST("/:id/cancel", orderHandler.CancelOrder)
		orderGroup.GET("/user/:user_id", orderHandler.GetOrdersByUser)
	}
}
