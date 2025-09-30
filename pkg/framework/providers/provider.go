package providers

import (
	"context"
	"go-modular-monolith/pkg/framework"
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
	Config    *framework.FrameworkConfig
	Container interface{} // DI Container
}

// ProviderManager manages all framework providers
type ProviderManager struct {
	providers []Provider
	config    *framework.FrameworkConfig
}

// NewProviderManager creates a new provider manager
func NewProviderManager(config *framework.FrameworkConfig) *ProviderManager {
	return &ProviderManager{
		providers: []Provider{},
		config:    config,
	}
}

// Register registers a new provider
func (pm *ProviderManager) Register(provider Provider) {
	pm.providers = append(pm.providers, provider)
}

// InitializeAll initializes all enabled providers
func (pm *ProviderManager) InitializeAll(ctx context.Context, deps Dependencies) error {
	for _, provider := range pm.providers {
		if provider.IsEnabled() {
			if err := provider.Initialize(ctx, deps); err != nil {
				return err
			}
		}
	}
	return nil
}

// ShutdownAll gracefully shuts down all providers
func (pm *ProviderManager) ShutdownAll(ctx context.Context) error {
	for _, provider := range pm.providers {
		if provider.IsEnabled() {
			if err := provider.Shutdown(ctx); err != nil {
				return err
			}
		}
	}
	return nil
}

// HealthCheckAll performs health check on all providers
func (pm *ProviderManager) HealthCheckAll(ctx context.Context) map[string]error {
	results := make(map[string]error)
	for _, provider := range pm.providers {
		if provider.IsEnabled() {
			results[provider.Name()] = provider.HealthCheck(ctx)
		}
	}
	return results
}
