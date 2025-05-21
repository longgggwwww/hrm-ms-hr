package grpc_clients

import (
	"fmt"
	"os"

	pb "github.com/huynhthanhthao/hrm_user_service/generated"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/longgggwwww/hrm-ms-hr/internal/utils"
)

// Tạo một client gRPC cho user service
func NewUserClient() (pb.UserServiceClient, error) {
	url := os.Getenv("USER_SERVICE")
	if url == "" {
		return nil, fmt.Errorf("USER_SERVICE environment variable is not set")
	}

	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, utils.WrapError("connect to user service", err)
	}

	return pb.NewUserServiceClient(conn), nil
}
