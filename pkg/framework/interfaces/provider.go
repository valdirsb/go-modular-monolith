package interfaces

import (
	"context"
)

// Provider interface that all framework providers must implement
type Provider interface {
	// Name returns the provider name
	Name() string

	// IsEnabled checks if this provider is enabled in configuration
	IsEnabled() bool

	// Initialize initializes the provider with dependencies
	Initialize(ctx context.Context, deps Dependencies) error

	// Shutdown gracefully shuts down the provider
	Shutdown(ctx context.Context) error

	// HealthCheck returns the health status of the provider
	HealthCheck(ctx context.Context) error
}

// Dependencies contains all the dependencies that providers might need
type Dependencies struct {
	Config    interface{} // Framework config
	Container interface{} // DI Container
}
