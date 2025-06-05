package grpc_clients

import (
	"context"
	"log"

	pb "github.com/huynhthanhthao/hrm_user_service/proto/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// UserServiceClient interface wraps the generated gRPC client
// This provides better abstraction and makes testing easier
type UserServiceClient interface {
	CreateUser(ctx context.Context, req *CreateUserRequest, opts ...grpc.CallOption) (*CreateUserResponse, error)
	UpdateUser(ctx context.Context, req *UpdateUserRequest, opts ...grpc.CallOption) (*UpdateUserResponse, error)
	GetUserById(ctx context.Context, req *GetUserByIdRequest, opts ...grpc.CallOption) (*GetUserByIdResponse, error)
	GetUsersByIDs(ctx context.Context, req *GetUsersByIDsRequest, opts ...grpc.CallOption) (*GetUsersByIDsResponse, error)
	DeleteUserByID(ctx context.Context, req *DeleteUserRequest, opts ...grpc.CallOption) (*DeleteUserResponse, error)
}

// userServiceClient implements UserServiceClient interface
type userServiceClient struct {
	client pb.UserServiceClient
}

// CreateUser creates a new user
func (c *userServiceClient) CreateUser(ctx context.Context, req *CreateUserRequest, opts ...grpc.CallOption) (*CreateUserResponse, error) {
	return c.client.CreateUser(ctx, req, opts...)
}

// UpdateUser updates an existing user
func (c *userServiceClient) UpdateUser(ctx context.Context, req *UpdateUserRequest, opts ...grpc.CallOption) (*UpdateUserResponse, error) {
	return c.client.UpdateUserByID(ctx, req, opts...)
}

// GetUserById retrieves a user by ID
func (c *userServiceClient) GetUserById(ctx context.Context, req *GetUserByIdRequest, opts ...grpc.CallOption) (*GetUserByIdResponse, error) {
	return c.client.GetUserById(ctx, req, opts...)
}

// GetUsersByIDs retrieves multiple users by their IDs
func (c *userServiceClient) GetUsersByIDs(ctx context.Context, req *GetUsersByIDsRequest, opts ...grpc.CallOption) (*GetUsersByIDsResponse, error) {
	return c.client.GetUsersByIDs(ctx, req, opts...)
}

// DeleteUserByID deletes a user by ID
func (c *userServiceClient) DeleteUserByID(ctx context.Context, req *DeleteUserRequest, opts ...grpc.CallOption) (*DeleteUserResponse, error) {
	return c.client.DeleteUserByID(ctx, req, opts...)
}

// NewUserClient creates a new user service client
func NewUserClient(grpcURL string) UserServiceClient {
	conn, err := grpc.NewClient(grpcURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("Failed to create gRPC client connection:", err)
		return nil
	}

	return &userServiceClient{
		client: pb.NewUserServiceClient(conn),
	}
}

// Re-export the protobuf types that other packages might need
// This allows other packages to use these types without importing the proto package directly
type (
	// User service request/response types
	CreateUserRequest     = pb.CreateUserRequest
	CreateUserResponse    = pb.CreateUserResponse
	UpdateUserRequest     = pb.UpdateUserRequest
	UpdateUserResponse    = pb.UpdateUserResponse
	GetUserByIdRequest    = pb.GetUserByIdRequest
	GetUserByIdResponse   = pb.GetUserByIdResponse
	GetUsersByIDsRequest  = pb.GetUsersByIDsRequest
	GetUsersByIDsResponse = pb.GetUsersByIDsResponse
	DeleteUserRequest     = pb.DeleteUserRequest
	DeleteUserResponse    = pb.DeleteUserResponse

	// User data types
	User    = pb.User
	Account = pb.Account
)
