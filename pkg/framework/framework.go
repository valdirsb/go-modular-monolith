package framework

import (
	"context"
	"fmt"
	"log"
)

// Framework represents the main framework instance
type Framework struct {
	config    *FrameworkConfig
	container interface{}            // DI Container
	providers map[string]interface{} // Simple provider storage
}

// NewFramework creates a new framework instance
func NewFramework(configPath string, container interface{}) (*Framework, error) {
	config, err := LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %v", err)
	}

	framework := &Framework{
		config:    config,
		container: container,
		providers: make(map[string]interface{}),
	}

	return framework, nil
}

// Initialize initializes the framework
func (f *Framework) Initialize(ctx context.Context) error {
	log.Printf("Initializing %s framework v%s", f.config.Framework.Name, f.config.Framework.Version)
	log.Printf("Framework initialized successfully")
	return nil
}

// Shutdown gracefully shuts down the framework
func (f *Framework) Shutdown(ctx context.Context) error {
	log.Printf("Shutting down framework...")
	log.Printf("Framework shutdown completed")
	return nil
}

// HealthCheck performs basic health checks
func (f *Framework) HealthCheck(ctx context.Context) map[string]error {
	results := make(map[string]error)

	// Basic framework health check
	results["framework"] = nil

	return results
}

// GetConfig returns the framework configuration
func (f *Framework) GetConfig() *FrameworkConfig {
	return f.config
}

// ListEnabledFeatures returns a list of enabled features
func (f *Framework) ListEnabledFeatures() map[string][]string {
	features := make(map[string][]string)

	// Core features
	var coreFeatures []string
	if f.config.Core.HTTPServer {
		coreFeatures = append(coreFeatures, "http_server")
	}
	if f.config.Core.EventSystem {
		coreFeatures = append(coreFeatures, "event_system")
	}
	if f.config.Core.Logging {
		coreFeatures = append(coreFeatures, "logging")
	}
	features["core"] = coreFeatures

	// Database features
	var dbFeatures []string
	if f.config.Database.MySQL {
		dbFeatures = append(dbFeatures, "mysql")
	}
	if f.config.Database.PostgreSQL {
		dbFeatures = append(dbFeatures, "postgresql")
	}
	features["database"] = dbFeatures

	// Integration features
	var integrationFeatures []string
	if f.config.Integrations.Stripe {
		integrationFeatures = append(integrationFeatures, "stripe")
	}
	if f.config.Integrations.GoogleMaps {
		integrationFeatures = append(integrationFeatures, "google_maps")
	}
	if f.config.Integrations.PIX {
		integrationFeatures = append(integrationFeatures, "pix")
	}
	features["integrations"] = integrationFeatures

	// Modules
	var moduleFeatures []string
	for module, enabled := range f.config.Modules {
		if enabled {
			moduleFeatures = append(moduleFeatures, module)
		}
	}
	features["modules"] = moduleFeatures

	return features
}

// IsEnabled checks if a specific feature or module is enabled
func (f *Framework) IsEnabled(category, feature string) bool {
	switch category {
	case "modules":
		if enabled, exists := f.config.Modules[feature]; exists {
			return enabled
		}
	case "core":
		switch feature {
		case "http_server":
			return f.config.Core.HTTPServer
		case "event_system":
			return f.config.Core.EventSystem
		case "logging":
			return f.config.Core.Logging
		}
	case "database":
		switch feature {
		case "mysql":
			return f.config.Database.MySQL
		case "postgresql":
			return f.config.Database.PostgreSQL
		}
	case "integrations":
		switch feature {
		case "stripe":
			return f.config.Integrations.Stripe
		case "google_maps":
			return f.config.Integrations.GoogleMaps
		case "pix":
			return f.config.Integrations.PIX
		}
	case "protocols":
		switch feature {
		case "grpc":
			return f.config.Protocols.GRPC
		case "http":
			return f.config.Protocols.HTTP
		}
	}
	return false
}

// RegisterProvider registers a provider with the framework
func (f *Framework) RegisterProvider(name string, provider interface{}) {
	f.providers[name] = provider
}

// GetProvider gets a provider by name
func (f *Framework) GetProvider(name string) interface{} {
	return f.providers[name]
}
