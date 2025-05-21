package grpc_clients

import (
	pb "github.com/huynhthanhthao/hrm_user_service/generated"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/longgggwwww/hrm-ms-hr/internal/utils"
)

// Tạo một client gRPC cho user service
func NewUserClient(url string) (pb.UserServiceClient, error) {
	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, utils.WrapError("connect to user service", err)
	}

	return pb.NewUserServiceClient(conn), nil
}
