package ports

import (
	"context"
	"go-modular-monolith/pkg/contracts"
)

// UserService define as operações de negócio do módulo de usuário (Primary Port)
type UserService interface {
	contracts.UserService
}

// UserRepository define a interface para persistência (Secondary Port)
type UserRepository interface {
	contracts.UserRepository
}

// PasswordHasher define a interface para hashing de senhas
type PasswordHasher interface {
	Hash(password string) (string, error)
	Verify(password, hash string) bool
}

// EmailService define a interface para envio de emails
type EmailService interface {
	SendWelcomeEmail(ctx context.Context, userID, email string) error
	SendPasswordResetEmail(ctx context.Context, userID, email, token string) error
}

// TokenGenerator define a interface para geração de tokens
type TokenGenerator interface {
	GenerateAccessToken(userID string) (string, error)
	GenerateRefreshToken(userID string) (string, error)
	ValidateToken(token string) (string, error) // retorna userID
}
