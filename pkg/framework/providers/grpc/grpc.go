package grpc

import (
	"context"
	"fmt"
	"net"

	"go-modular-monolith/pkg/framework"
	"go-modular-monolith/pkg/framework/providers"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

// GRPCProvider provides gRPC server functionality
type GRPCProvider struct {
	server   *grpc.Server
	listener net.Listener
	config   *framework.FrameworkConfig
	port     string
	services []ServiceRegistrar
}

// ServiceRegistrar defines how services register with gRPC
type ServiceRegistrar interface {
	RegisterWithServer(server *grpc.Server)
}

// NewGRPCProvider creates a new gRPC provider
func NewGRPCProvider(port string) *GRPCProvider {
	return &GRPCProvider{
		port:     port,
		services: []ServiceRegistrar{},
	}
}

// Name returns the provider name
func (g *GRPCProvider) Name() string {
	return "grpc"
}

// IsEnabled checks if gRPC is enabled in configuration
func (g *GRPCProvider) IsEnabled() bool {
	return framework.IsEnabled("protocols", "grpc")
}

// Initialize initializes the gRPC server
func (g *GRPCProvider) Initialize(ctx context.Context, deps providers.Dependencies) error {
	g.config = deps.Config

	// Create listener
	listener, err := net.Listen("tcp", ":"+g.port)
	if err != nil {
		return fmt.Errorf("failed to listen on port %s: %v", g.port, err)
	}
	g.listener = listener

	// Create gRPC server with options
	opts := []grpc.ServerOption{
		grpc.MaxRecvMsgSize(4 * 1024 * 1024), // 4MB
		grpc.MaxSendMsgSize(4 * 1024 * 1024), // 4MB
	}

	g.server = grpc.NewServer(opts...)

	// Register health service
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(g.server, healthServer)
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

	// Enable reflection for development
	if framework.IsEnabled("development", "debug_mode") {
		reflection.Register(g.server)
	}

	// Register all services
	for _, service := range g.services {
		service.RegisterWithServer(g.server)
	}

	// Start server in goroutine
	go func() {
		if err := g.server.Serve(g.listener); err != nil {
			fmt.Printf("gRPC server failed to serve: %v\n", err)
		}
	}()

	return nil
}

// RegisterService registers a gRPC service
func (g *GRPCProvider) RegisterService(service ServiceRegistrar) {
	g.services = append(g.services, service)
}

// Shutdown gracefully stops the gRPC server
func (g *GRPCProvider) Shutdown(ctx context.Context) error {
	if g.server != nil {
		// Graceful shutdown with timeout
		done := make(chan struct{})
		go func() {
			g.server.GracefulStop()
			close(done)
		}()

		select {
		case <-done:
			return nil
		case <-ctx.Done():
			g.server.Stop()
			return ctx.Err()
		}
	}
	return nil
}

// HealthCheck checks if the gRPC server is running
func (g *GRPCProvider) HealthCheck(ctx context.Context) error {
	if g.server == nil {
		return fmt.Errorf("gRPC server not initialized")
	}
	if g.listener == nil {
		return fmt.Errorf("gRPC listener not available")
	}
	return nil
}

// GetServer returns the gRPC server instance
func (g *GRPCProvider) GetServer() *grpc.Server {
	return g.server
}

// GetPort returns the port the server is listening on
func (g *GRPCProvider) GetPort() string {
	return g.port
}
