package grpc_clients

import (
	"fmt"
	"os"

	pb "github.com/huynhthanhthao/hrm_user_service/generated"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// NewUserClient creates and returns a new UserServiceClient connected to the USER_SERVICE address.
func NewUserClient() (pb.UserServiceClient, error) {
	url := os.Getenv("USER_SERVICE")
	if url == "" {
		return nil, fmt.Errorf("USER_SERVICE environment variable is not set")
	}

	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to user service: %w", err)
	}

	return pb.NewUserServiceClient(conn), nil
}
