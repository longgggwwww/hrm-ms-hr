package main

import (
	"fmt"
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

func startGRPCServer(cli *ent.Client) {
	branch := entpb.NewBranchService(cli)
	company := entpb.NewCompanyService(cli)
	employee := entpb.NewEmployeeService(cli)
	position := entpb.NewPositionService(cli)
	department := entpb.NewDepartmentService(cli)
	project := entpb.NewProjectService(cli)
	task := entpb.NewTaskService(cli)

	server := grpc.NewServer()
	fmt.Println("Starting gRPC server on port 5000...")

	entpb.RegisterBranchServiceServer(server, branch)
	entpb.RegisterCompanyServiceServer(server, company)
	entpb.RegisterEmployeeServiceServer(server, employee)
	entpb.RegisterPositionServiceServer(server, position)
	entpb.RegisterDepartmentServiceServer(server, department)
	entpb.RegisterProjectServiceServer(server, project)
	entpb.RegisterTaskServiceServer(server, task)
	entpb.RegisterExtServiceServer(server, entpb.NewExtService(cli))

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

	employee := handlers.NewEmployeeHandler(cli, nil)
	employee.RegisterRoutes(r)
	branch := handlers.NewBranchHandler(cli, nil)
	branch.RegisterRoutes(r)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

func main() {
	DB_URL := os.Getenv("DB_URL")
	if DB_URL == "" {
		log.Fatal("DB_URL environment variable is not set")
	}

	cli, err := ent.Open("postgres", DB_URL)
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}
	defer cli.Close()

	// Initialize gRPC clients
	userClient, err := grpc_clients.NewUserClient()
	if err != nil {
		log.Fatalf("failed to initialize user client: %v", err)
	}

	fmt.Println("User client initialized successfully", userClient)

	go startHTTPServer(cli)
	startGRPCServer(cli)
}
