package grpc_clients

import (
	"log"

	pb "github.com/huynhthanhthao/hrm_user_service/generated"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewUserClient(grpcURL string) pb.UserServiceClient {
	conn, err := grpc.NewClient(grpcURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("Failed to create gRPC client connection:", err)
		return nil
	}
	return pb.NewUserServiceClient(conn)
}
