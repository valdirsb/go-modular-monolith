package bootstrap

import (
	"context"
	"fmt"
	"log"
	"sync"

	"go-modular-monolith/pkg/contracts"
)

// SimpleLogger implementa uma versão simples da interface Logger
type SimpleLogger struct{}

func (l *SimpleLogger) Debug(msg string, fields ...contracts.Field) {
	log.Printf("[DEBUG] %s %v", msg, fields)
}

func (l *SimpleLogger) Info(msg string, fields ...contracts.Field) {
	log.Printf("[INFO] %s %v", msg, fields)
}

func (l *SimpleLogger) Warn(msg string, fields ...contracts.Field) {
	log.Printf("[WARN] %s %v", msg, fields)
}

func (l *SimpleLogger) Error(msg string, fields ...contracts.Field) {
	log.Printf("[ERROR] %s %v", msg, fields)
}

func (l *SimpleLogger) Fatal(msg string, fields ...contracts.Field) {
	log.Fatalf("[FATAL] %s %v", msg, fields)
}

func (l *SimpleLogger) With(fields ...contracts.Field) contracts.Logger {
	return l
}

// MockEmailService implementa uma versão mock do EmailService
type MockEmailService struct{}

func (m *MockEmailService) SendWelcomeEmail(ctx context.Context, userID, email string) error {
	log.Printf("[EMAIL] Welcome email sent to %s (User: %s)", email, userID)
	return nil
}

func (m *MockEmailService) SendPasswordResetEmail(ctx context.Context, userID, email, token string) error {
	log.Printf("[EMAIL] Password reset email sent to %s (User: %s, Token: %s)", email, userID, token)
	return nil
}

// MockTokenGenerator implementa uma versão mock do TokenGenerator
type MockTokenGenerator struct{}

func (m *MockTokenGenerator) GenerateAccessToken(userID string) (string, error) {
	return fmt.Sprintf("access_token_for_%s", userID), nil
}

func (m *MockTokenGenerator) GenerateRefreshToken(userID string) (string, error) {
	return fmt.Sprintf("refresh_token_for_%s", userID), nil
}

func (m *MockTokenGenerator) ValidateToken(token string) (string, error) {
	// Implementação simplificada - em produção usaria JWT
	return "user_id_from_token", nil
}

// InMemoryUserRepository implementa uma versão em memória do UserRepository
type InMemoryUserRepository struct {
	users map[string]*contracts.User
	mu    sync.RWMutex
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users: make(map[string]*contracts.User),
	}
}

func (r *InMemoryUserRepository) Create(ctx context.Context, user *contracts.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users[user.ID] = user
	return nil
}

func (r *InMemoryUserRepository) GetByID(ctx context.Context, id string) (*contracts.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	user, exists := r.users[id]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

func (r *InMemoryUserRepository) GetByEmail(ctx context.Context, email string) (*contracts.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

func (r *InMemoryUserRepository) Update(ctx context.Context, user *contracts.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.users[user.ID]; !exists {
		return fmt.Errorf("user not found")
	}
	r.users[user.ID] = user
	return nil
}

func (r *InMemoryUserRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.users[id]; !exists {
		return fmt.Errorf("user not found")
	}
	delete(r.users, id)
	return nil
}
