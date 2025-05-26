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

func registerGRPCServices(srv *grpc.Server, cli *ent.Client) {
	entpb.RegisterOrganizationServiceServer(srv, entpb.NewOrganizationService(cli))
	entpb.RegisterDepartmentServiceServer(srv, entpb.NewDepartmentService(cli))
	entpb.RegisterPositionServiceServer(srv, entpb.NewPositionService(cli))
	entpb.RegisterEmployeeServiceServer(srv, entpb.NewEmployeeService(cli))
	entpb.RegisterProjectServiceServer(srv, entpb.NewProjectService(cli))
	entpb.RegisterTaskServiceServer(srv, entpb.NewTaskService(cli))
	entpb.RegisterExtServiceServer(srv, entpb.NewExtService(cli))
}

func startGRPCServer(cli *ent.Client) {
	srv := grpc.NewServer()
	registerGRPCServices(srv, cli)

	log.Println("Starting gRPC server on port 5000...")
	lis, err := net.Listen("tcp", ":5000")
	if err != nil {
		log.Fatalf("failed listening: %s", err)
	}

	if err := srv.Serve(lis); err != nil {
		log.Fatalf("server ended: %s", err)
	}
}

func startHTTPServer(cli *ent.Client) {
	r := gin.Default()

	log.Println("Connecting to user service at:", os.Getenv("USER_SERVICE"))
	userServ := grpc_clients.NewUserClient(os.Getenv("USER_SERVICE"))

	handlersList := []struct {
		register func(*gin.Engine)
	}{
		{handlers.NewEmployeeHandler(cli, userServ).RegisterRoutes},
		{handlers.NewOrgHandler(cli, nil).RegisterRoutes},
		{handlers.NewDeptHandler(cli, nil).RegisterRoutes},
		{handlers.NewPositionHandler(cli, nil).RegisterRoutes},
		{handlers.NewProjectHandler(cli).RegisterRoutes},
	}
	for _, h := range handlersList {
		h.register(r)
	}

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

func main() {
	connStr := os.Getenv("DB_URL")
	if connStr == "" {
		log.Fatal("DB_URL environment variable is not set")
	}

	cli, err := ent.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}
	defer cli.Close()

	go startHTTPServer(cli)
	startGRPCServer(cli)
}
