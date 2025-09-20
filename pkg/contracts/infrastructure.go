package contracts

import (
	"context"
	"time"
)

// Config define a interface para configurações globais
type Config interface {
	GetDatabaseURL() string
	GetServerAddress() string
	GetJWTSecret() string
	GetEnvironment() string
	IsProduction() bool
	GetLogLevel() string
}

// Logger define a interface para logging
type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)
	With(fields ...Field) Logger
}

type Field struct {
	Key   string
	Value interface{}
}

// Validator define a interface para validação
type Validator interface {
	Validate(s interface{}) error
}

// Cache define a interface para cache
type Cache interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (interface{}, error)
	Delete(ctx context.Context, key string) error
}

// Database define a interface para transações de banco
type Database interface {
	BeginTx(ctx context.Context) (Transaction, error)
	Health() error
}

type Transaction interface {
	Commit() error
	Rollback() error
	UserRepository() UserRepository
	ProductRepository() ProductRepository
	OrderRepository() OrderRepository
}
