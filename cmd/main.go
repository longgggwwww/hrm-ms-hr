package main

// import (
// 	"fmt"
// 	"log"
// 	"net"
// 	"os"

// 	"github.com/gin-gonic/gin"
// 	_ "github.com/lib/pq"
// 	"github.com/longgggwwww/hrm-ms-hr/ent"
// 	"github.com/longgggwwww/hrm-ms-hr/ent/proto/entpb"
// 	"github.com/longgggwwww/hrm-ms-hr/internal/grpc_clients"
// 	"github.com/longgggwwww/hrm-ms-hr/internal/handlers"
// 	"google.golang.org/grpc"
// )

// func startGRPCServer(cli *ent.Client) {
// 	perm := entpb.NewBranchService(cli)
// 	permGroup := entpb.NewCompanyService(cli)

// 	server := grpc.NewServer()
// 	fmt.Println("Starting gRPC server on port 5000...")

// 	entpb.RegisterBranchServiceServer(server, perm)
// 	entpb.RegisterCompanyServiceServer(server, permGroup)

// 	lis, err := net.Listen("tcp", ":5000")
// 	if err != nil {
// 		log.Fatalf("failed listening: %s", err)
// 	}

// 	if err := server.Serve(lis); err != nil {
// 		log.Fatalf("server ended: %s", err)
// 	}
// }

// func startHTTPServer(cli *ent.Client) {
// 	r := gin.Default()

// 	if err := r.Run(":8080"); err != nil {
// 		log.Fatalf("failed to start server: %v", err)
// 	}
// }

// func main() {
// 	DB_URL := os.Getenv("DB_URL")
// 	if DB_URL == "" {
// 		log.Fatal("DB_URL environment variable is not set")
// 	}

// 	cli, err := ent.Open("postgres", DB_URL)
// 	if err != nil {
// 		log.Fatalf("failed opening connection to postgres: %v", err)
// 	}
// 	defer cli.Close()

// 	// Initialize gRPC clients
// 	userClient, err := grpc_clients.NewUserClient()
// 	if err != nil {
// 		log.Fatalf("failed to initialize user client: %v", err)
// 	}

// 	roleHandler := handlers.NewRoleHandler(cli, userClient)

// 	// Pass roleHandler to HTTP server
// 	go startHTTPServer(cli, roleHandler)
// 	startGRPCServer(cli)
// }
