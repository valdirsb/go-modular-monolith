package framework

import (
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

type FrameworkConfig struct {
	Framework    CoreFramework       `yaml:"framework"`
	Core         CoreFeatures        `yaml:"core"`
	Database     DatabaseProviders   `yaml:"database"`
	Cache        CacheProviders      `yaml:"cache"`
	Protocols    ProtocolFeatures    `yaml:"protocols"`
	Integrations IntegrationFeatures `yaml:"integrations"`
	Security     SecurityFeatures    `yaml:"security"`
	Development  DevelopmentFeatures `yaml:"development"`
	Modules      map[string]bool     `yaml:"modules"`
}

type CoreFramework struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}

type CoreFeatures struct {
	HTTPServer          bool `yaml:"http_server"`
	DependencyInjection bool `yaml:"dependency_injection"`
	EventSystem         bool `yaml:"event_system"`
	Logging             bool `yaml:"logging"`
	DatabaseMigrations  bool `yaml:"database_migrations"`
}

type DatabaseProviders struct {
	MySQL      bool `yaml:"mysql"`
	PostgreSQL bool `yaml:"postgresql"`
	MongoDB    bool `yaml:"mongodb"`
	SQLite     bool `yaml:"sqlite"`
}

type CacheProviders struct {
	Redis     bool `yaml:"redis"`
	Memcached bool `yaml:"memcached"`
	InMemory  bool `yaml:"in_memory"`
}

type ProtocolFeatures struct {
	HTTP       bool `yaml:"http"`
	GRPC       bool `yaml:"grpc"`
	WebSockets bool `yaml:"websockets"`
	GraphQL    bool `yaml:"graphql"`
}

type IntegrationFeatures struct {
	// Authentication & Authorization
	JWT          bool `yaml:"jwt"`
	OAuth2       bool `yaml:"oauth2"`
	FirebaseAuth bool `yaml:"firebase_auth"`
	Auth0        bool `yaml:"auth0"`

	// Cloud services
	AWSS3              bool `yaml:"aws_s3"`
	GoogleCloudStorage bool `yaml:"google_cloud_storage"`
	AzureBlob          bool `yaml:"azure_blob"`

	// Payment gateways
	Stripe      bool `yaml:"stripe"`
	PayPal      bool `yaml:"paypal"`
	MercadoPago bool `yaml:"mercado_pago"`
	PIX         bool `yaml:"pix"`

	// Maps & Location
	GoogleMaps bool `yaml:"google_maps"`
	Mapbox     bool `yaml:"mapbox"`

	// Messaging
	RabbitMQ bool `yaml:"rabbitmq"`
	Kafka    bool `yaml:"kafka"`
	SQS      bool `yaml:"sqs"`
	PubSub   bool `yaml:"pub_sub"`

	// Email services
	SendGrid bool `yaml:"sendgrid"`
	SES      bool `yaml:"ses"`
	SMTP     bool `yaml:"smtp"`

	// Push notifications
	FirebaseMessaging bool `yaml:"firebase_messaging"`
	APNS              bool `yaml:"apns"`

	// Monitoring & Observability
	Prometheus bool `yaml:"prometheus"`
	Jaeger     bool `yaml:"jaeger"`
	DataDog    bool `yaml:"datadog"`
	NewRelic   bool `yaml:"new_relic"`

	// External APIs
	SocialMedia bool `yaml:"social_media"`
	WeatherAPI  bool `yaml:"weather_api"`
	CurrencyAPI bool `yaml:"currency_api"`
}

type SecurityFeatures struct {
	RateLimiting   bool `yaml:"rate_limiting"`
	CORS           bool `yaml:"cors"`
	Helmet         bool `yaml:"helmet"`
	CSRFProtection bool `yaml:"csrf_protection"`
	Encryption     bool `yaml:"encryption"`
}

type DevelopmentFeatures struct {
	Swagger   bool `yaml:"swagger"`
	DebugMode bool `yaml:"debug_mode"`
	HotReload bool `yaml:"hot_reload"`
	Profiler  bool `yaml:"profiler"`
}

var globalConfig *FrameworkConfig

// LoadConfig loads the framework configuration from yaml file
func LoadConfig(configPath string) (*FrameworkConfig, error) {
	if configPath == "" {
		configPath = "framework.yaml"
	}

	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	var config FrameworkConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}

	globalConfig = &config
	return &config, nil
}

// GetConfig returns the global framework configuration
func GetConfig() *FrameworkConfig {
	if globalConfig == nil {
		log.Fatal("Framework config not loaded. Call LoadConfig() first.")
	}
	return globalConfig
}

// IsEnabled checks if a specific feature is enabled
func IsEnabled(category string, feature string) bool {
	config := GetConfig()

	switch category {
	case "core":
		return isCoreFatureEnabled(config.Core, feature)
	case "database":
		return isDatabaseEnabled(config.Database, feature)
	case "cache":
		return isCacheEnabled(config.Cache, feature)
	case "protocols":
		return isProtocolEnabled(config.Protocols, feature)
	case "integrations":
		return isIntegrationEnabled(config.Integrations, feature)
	case "security":
		return isSecurityEnabled(config.Security, feature)
	case "development":
		return isDevelopmentEnabled(config.Development, feature)
	case "modules":
		return config.Modules[feature]
	}

	return false
}

func isCoreFatureEnabled(core CoreFeatures, feature string) bool {
	switch feature {
	case "http_server":
		return core.HTTPServer
	case "dependency_injection":
		return core.DependencyInjection
	case "event_system":
		return core.EventSystem
	case "logging":
		return core.Logging
	case "database_migrations":
		return core.DatabaseMigrations
	}
	return false
}

func isDatabaseEnabled(db DatabaseProviders, feature string) bool {
	switch feature {
	case "mysql":
		return db.MySQL
	case "postgresql":
		return db.PostgreSQL
	case "mongodb":
		return db.MongoDB
	case "sqlite":
		return db.SQLite
	}
	return false
}

func isCacheEnabled(cache CacheProviders, feature string) bool {
	switch feature {
	case "redis":
		return cache.Redis
	case "memcached":
		return cache.Memcached
	case "in_memory":
		return cache.InMemory
	}
	return false
}

func isProtocolEnabled(protocols ProtocolFeatures, feature string) bool {
	switch feature {
	case "http":
		return protocols.HTTP
	case "grpc":
		return protocols.GRPC
	case "websockets":
		return protocols.WebSockets
	case "graphql":
		return protocols.GraphQL
	}
	return false
}

func isIntegrationEnabled(integrations IntegrationFeatures, feature string) bool {
	switch feature {
	case "jwt":
		return integrations.JWT
	case "oauth2":
		return integrations.OAuth2
	case "firebase_auth":
		return integrations.FirebaseAuth
	case "auth0":
		return integrations.Auth0
	case "aws_s3":
		return integrations.AWSS3
	case "google_cloud_storage":
		return integrations.GoogleCloudStorage
	case "azure_blob":
		return integrations.AzureBlob
	case "stripe":
		return integrations.Stripe
	case "paypal":
		return integrations.PayPal
	case "mercado_pago":
		return integrations.MercadoPago
	case "pix":
		return integrations.PIX
	case "google_maps":
		return integrations.GoogleMaps
	case "mapbox":
		return integrations.Mapbox
	case "rabbitmq":
		return integrations.RabbitMQ
	case "kafka":
		return integrations.Kafka
	case "sqs":
		return integrations.SQS
	case "pub_sub":
		return integrations.PubSub
	case "sendgrid":
		return integrations.SendGrid
	case "ses":
		return integrations.SES
	case "smtp":
		return integrations.SMTP
	case "firebase_messaging":
		return integrations.FirebaseMessaging
	case "apns":
		return integrations.APNS
	case "prometheus":
		return integrations.Prometheus
	case "jaeger":
		return integrations.Jaeger
	case "datadog":
		return integrations.DataDog
	case "new_relic":
		return integrations.NewRelic
	case "social_media":
		return integrations.SocialMedia
	case "weather_api":
		return integrations.WeatherAPI
	case "currency_api":
		return integrations.CurrencyAPI
	}
	return false
}

func isSecurityEnabled(security SecurityFeatures, feature string) bool {
	switch feature {
	case "rate_limiting":
		return security.RateLimiting
	case "cors":
		return security.CORS
	case "helmet":
		return security.Helmet
	case "csrf_protection":
		return security.CSRFProtection
	case "encryption":
		return security.Encryption
	}
	return false
}

func isDevelopmentEnabled(dev DevelopmentFeatures, feature string) bool {
	switch feature {
	case "swagger":
		return dev.Swagger
	case "debug_mode":
		return dev.DebugMode
	case "hot_reload":
		return dev.HotReload
	case "profiler":
		return dev.Profiler
	}
	return false
}
