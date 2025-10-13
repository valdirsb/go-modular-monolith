package services

import (
	"context"
	"fmt"
	"time"

	"go-modular-monolith/pkg/contracts"
	pb "go-modular-monolith/pkg/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UserGRPCService implements the gRPC UserService
type UserGRPCService struct {
	pb.UnimplementedUserServiceServer
	userService contracts.UserService
	userRepo    contracts.UserRepository
}

// NewUserGRPCService creates a new gRPC user service
func NewUserGRPCService(userService contracts.UserService, userRepo contracts.UserRepository) *UserGRPCService {
	return &UserGRPCService{
		userService: userService,
		userRepo:    userRepo,
	}
}

// RegisterWithServer registers the service with the gRPC server
func (s *UserGRPCService) RegisterWithServer(server *grpc.Server) {
	pb.RegisterUserServiceServer(server, s)
}

// CreateUser creates a new user
func (s *UserGRPCService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	// Validate input
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}
	if req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	// Create user request
	createReq := contracts.CreateUserRequest{
		Username: req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	// Call service
	user, err := s.userService.CreateUser(ctx, createReq)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to create user: %v", err))
	}

	// Convert to proto message
	protoUser := &pb.User{
		Id:        user.ID,
		Name:      user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}

	return &pb.CreateUserResponse{
		User:    protoUser,
		Message: "User created successfully",
	}, nil
}

// GetUser gets a user by ID
func (s *UserGRPCService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	user, err := s.userService.GetUserByID(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("user not found: %v", err))
	}

	protoUser := &pb.User{
		Id:        user.ID,
		Name:      user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}

	return &pb.GetUserResponse{
		User:    protoUser,
		Message: "User retrieved successfully",
	}, nil
}

// UpdateUser updates a user
func (s *UserGRPCService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	name := req.Name
	email := req.Email
	updateReq := contracts.UpdateUserRequest{
		Username: &name,
		Email:    &email,
	}

	user, err := s.userService.UpdateUser(ctx, req.Id, updateReq)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to update user: %v", err))
	}

	protoUser := &pb.User{
		Id:        user.ID,
		Name:      user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}

	return &pb.UpdateUserResponse{
		User:    protoUser,
		Message: "User updated successfully",
	}, nil
}

// DeleteUser deletes a user
func (s *UserGRPCService) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	err := s.userService.DeleteUser(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to delete user: %v", err))
	}

	return &pb.DeleteUserResponse{
		Message: "User deleted successfully",
	}, nil
}

// ListUsers lists users with pagination
func (s *UserGRPCService) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	// Para simplificar, vamos retornar apenas uma resposta vazia por enquanto
	// Em uma implementação real, você criaria um método no repository para paginação
	return &pb.ListUsersResponse{
		Users:    []*pb.User{},
		Total:    0,
		Page:     req.Page,
		PageSize: req.PageSize,
		Message:  "Users retrieved successfully",
	}, nil
}

// GetUserByEmail gets a user by email
func (s *UserGRPCService) GetUserByEmail(ctx context.Context, req *pb.GetUserByEmailRequest) (*pb.GetUserResponse, error) {
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("user not found: %v", err))
	}

	protoUser := &pb.User{
		Id:        user.ID,
		Name:      user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}

	return &pb.GetUserResponse{
		User:    protoUser,
		Message: "User retrieved successfully",
	}, nil
}
