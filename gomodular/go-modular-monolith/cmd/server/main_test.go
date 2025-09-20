package main_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-modular-monolith/internal/bootstrap"
	// "go-modular-monolith/pkg/container"
	"go-modular-monolith/pkg/contracts"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserAPI(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	container, err := bootstrap.Bootstrap()
	require.NoError(t, err)

	router := gin.New()
	// registerUserRoutes(router, container)

	// Test Create User
	t.Run("Create User", func(t *testing.T) {
		user := contracts.CreateUserRequest{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "password123",
		}

		body, _ := json.Marshal(user)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/users/", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response contracts.User
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "testuser", response.Username)
		assert.Equal(t, "test@example.com", response.Email)
		assert.NotEmpty(t, response.ID)
	})

	// Test Get User
	t.Run("Get User", func(t *testing.T) {
		// Primeiro, criar um usuário
		userService := container.MustGet("userService").(contracts.UserService)

		createReq := contracts.CreateUserRequest{
			Username: "getuser",
			Email:    "get@example.com",
			Password: "password123",
		}

		createdUser, err := userService.CreateUser(context.Background(), createReq)
		require.NoError(t, err)

		// Agora, buscar o usuário
		req := httptest.NewRequest(http.MethodGet, "/api/v1/users/"+createdUser.ID, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response contracts.User
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, createdUser.ID, response.ID)
		assert.Equal(t, "getuser", response.Username)
	})

	// Test Update User
	t.Run("Update User", func(t *testing.T) {
		// Criar usuário
		userService := container.MustGet("userService").(contracts.UserService)

		createReq := contracts.CreateUserRequest{
			Username: "updateuser",
			Email:    "update@example.com",
			Password: "password123",
		}

		createdUser, err := userService.CreateUser(context.Background(), createReq)
		require.NoError(t, err)

		// Atualizar usuário
		newUsername := "updateduser"
		updateReq := contracts.UpdateUserRequest{
			Username: &newUsername,
		}

		body, _ := json.Marshal(updateReq)
		req := httptest.NewRequest(http.MethodPut, "/api/v1/users/"+createdUser.ID, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response contracts.User
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "updateduser", response.Username)
	})

	// Test Validate User
	t.Run("Validate User", func(t *testing.T) {
		// Criar usuário
		userService := container.MustGet("userService").(contracts.UserService)

		createReq := contracts.CreateUserRequest{
			Username: "validateuser",
			Email:    "validate@example.com",
			Password: "password123",
		}

		_, err := userService.CreateUser(context.Background(), createReq)
		require.NoError(t, err)

		// Validar credenciais
		validateReq := map[string]string{
			"email":    "validate@example.com",
			"password": "password123",
		}

		body, _ := json.Marshal(validateReq)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/users/validate", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response contracts.User
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "validateuser", response.Username)
	})
}

// Função helper para registrar rotas (copiada do main.go)
// func registerUserRoutes(router *gin.Engine, container *container.Container) {
// 	userService := container.MustGet("userService").(contracts.UserService)

// 	userGroup := router.Group("/api/v1/users")
// 	{
// 		userGroup.POST("/", createUserHandler(userService))
// 		userGroup.GET("/:id", getUserHandler(userService))
// 		userGroup.PUT("/:id", updateUserHandler(userService))
// 		userGroup.DELETE("/:id", deleteUserHandler(userService))
// 		userGroup.POST("/validate", validateUserHandler(userService))
// 	}
// }
