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
	// registerProductRoutes(router, container) // Implementar conforme necessário
	// registerOrderRoutes(router, container)   // Implementar conforme necessário

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
	userService := container.MustGet("userService").(contracts.UserService)

	userGroup := router.Group("/api/v1/users")
	{
		userGroup.POST("/", createUserHandler(userService))
		userGroup.GET("/:id", getUserHandler(userService))
		userGroup.PUT("/:id", updateUserHandler(userService))
		userGroup.DELETE("/:id", deleteUserHandler(userService))
		userGroup.POST("/validate", validateUserHandler(userService))
	}
}

// Handlers HTTP para o módulo de usuário

func createUserHandler(userService contracts.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req contracts.CreateUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err := userService.CreateUser(c.Request.Context(), req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, user)
	}
}

func getUserHandler(userService contracts.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		user, err := userService.GetUserByID(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

func updateUserHandler(userService contracts.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var req contracts.UpdateUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err := userService.UpdateUser(c.Request.Context(), id, req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

func deleteUserHandler(userService contracts.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		err := userService.DeleteUser(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusNoContent, nil)
	}
}

func validateUserHandler(userService contracts.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Email    string `json:"email" binding:"required"`
			Password string `json:"password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err := userService.ValidateUser(c.Request.Context(), req.Email, req.Password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}
