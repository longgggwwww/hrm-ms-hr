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

func registerGRPCServices(sv *grpc.Server, cli *ent.Client) {
	entpb.RegisterOrganizationServiceServer(sv, entpb.NewOrganizationService(cli))
	entpb.RegisterDepartmentServiceServer(sv, entpb.NewDepartmentService(cli))
	entpb.RegisterPositionServiceServer(sv, entpb.NewPositionService(cli))
	entpb.RegisterEmployeeServiceServer(sv, entpb.NewEmployeeService(cli))
	entpb.RegisterProjectServiceServer(sv, entpb.NewProjectService(cli))
	entpb.RegisterTaskServiceServer(sv, entpb.NewTaskService(cli))
	entpb.RegisterExtServiceServer(sv, entpb.NewExtService(cli))
}

func startGRPCServer(cli *ent.Client) {
	serv := grpc.NewServer()
	log.Println("Starting gRPC server on port 5000...")
	registerGRPCServices(serv, cli)

	lis, err := net.Listen("tcp", ":5000")
	if err != nil {
		log.Fatalf("failed listening: %s", err)
	}

	if err := serv.Serve(lis); err != nil {
		log.Fatalf("server ended: %s", err)
	}
}

func startHTTPServer(cli *ent.Client) {
	r := gin.Default()

	log.Println("Connecting to user service at:", os.Getenv("USER_SERVICE"))
	user, err := grpc_clients.NewUserClient(os.Getenv("USER_SERVICE"))
	if err != nil {
		log.Fatalf("failed to create user client: %v", err)
	}

	// Đăng ký các route cho HTTP server
	handlersList := []struct {
		register func(*gin.Engine)
	}{
		{handlers.NewEmployeeHandler(cli, user).RegisterRoutes},
		{handlers.NewOrgHandler(cli, nil).RegisterRoutes},
		{handlers.NewDeptHandler(cli, nil).RegisterRoutes},
		{handlers.NewPositionHandler(cli, nil).RegisterRoutes},
	}
	for _, h := range handlersList {
		h.register(r)
	}

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

func main() {
	// Initialize database client
	connStr := os.Getenv("DB_URL")
	if connStr == "" {
		log.Fatal("DB_URL environment variable is not set")
	}
	cli, err := ent.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}
	defer cli.Close()

	// Start HTTP and gRPC servers
	go startHTTPServer(cli)
	startGRPCServer(cli)
}
