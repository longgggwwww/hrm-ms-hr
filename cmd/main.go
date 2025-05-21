package main

import (
	"log"
	"net"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/longgggwwww/hrm-ms-hr/ent"
	"github.com/longgggwwww/hrm-ms-hr/ent/proto/entpb"
	"github.com/longgggwwww/hrm-ms-hr/internal/grpc_clients"
	"github.com/longgggwwww/hrm-ms-hr/internal/handlers"
	"google.golang.org/grpc"
)

func registerGRPCServices(server *grpc.Server, cli *ent.Client) {
	entpb.RegisterOrganizationServiceServer(server, entpb.NewOrganizationService(cli))
	entpb.RegisterDepartmentServiceServer(server, entpb.NewDepartmentService(cli))
	entpb.RegisterPositionServiceServer(server, entpb.NewPositionService(cli))
	entpb.RegisterEmployeeServiceServer(server, entpb.NewEmployeeService(cli))
	entpb.RegisterProjectServiceServer(server, entpb.NewProjectService(cli))
	entpb.RegisterTaskServiceServer(server, entpb.NewTaskService(cli))
	entpb.RegisterHRExtServiceServer(server, entpb.NewHRExtService(cli))
}

func startGRPCServer(cli *ent.Client) {
	server := grpc.NewServer()
	log.Println("Starting gRPC server on port 5000...")
	registerGRPCServices(server, cli)

	lis, err := net.Listen("tcp", ":5000")
	if err != nil {
		log.Fatalf("failed listening: %s", err)
	}

	if err := server.Serve(lis); err != nil {
		log.Fatalf("server ended: %s", err)
	}
}

func startHTTPServer(cli *ent.Client) {
	r := gin.Default()

	// Register all handlers and routes
	handlersList := []struct {
		register func(*gin.Engine)
	}{
		{handlers.NewEmployeeHandler(cli, nil).RegisterRoutes},
		{handlers.NewOrgHandler(cli, nil).RegisterRoutes},
		{handlers.NewDepartmentHandler(cli, nil).RegisterRoutes},
		{handlers.NewPositionHandler(cli, nil).RegisterRoutes},
	}
	for _, h := range handlersList {
		h.register(r)
	}

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

// getDBClient loads the DB_URL from environment and returns a connected ent.Client.
func getDBClient() *ent.Client {
	connStr := os.Getenv("DB_URL")
	if connStr == "" {
		log.Fatal("DB_URL environment variable is not set")
	}
	client, err := ent.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}
	return client
}

// initUserClient initializes the gRPC user client.
func initUserClient() (interface{}, error) {
	return grpc_clients.NewUserClient()
}

func main() {
	// Initialize database client
	cli := getDBClient()
	defer cli.Close()

	// Initialize gRPC clients
	userClient, err := initUserClient()
	if err != nil {
		log.Fatalf("failed to initialize user client: %v", err)
	}
	log.Println("User client initialized successfully", userClient)

	// Start HTTP and gRPC servers
	go startHTTPServer(cli)
	startGRPCServer(cli)
}
