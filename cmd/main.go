package main

import (
	"log"
	"net"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"

	"github.com/longgggwwww/hrm-ms-hr/ent"
	"github.com/longgggwwww/hrm-ms-hr/ent/proto/entpb"
	"github.com/longgggwwww/hrm-ms-hr/internal/grpc_clients"
	"github.com/longgggwwww/hrm-ms-hr/internal/handlers"
	"github.com/longgggwwww/hrm-ms-hr/internal/utils"
)

func main() {
	// Initialize custom validators
	utils.InitValidator()

	cli := initDatabase()
	defer cli.Close()

	log.Println("Starting HR microservice...")
	go startHTTPServer(cli)
	startGRPCServer(cli)
}

func initDatabase() *ent.Client {
	connStr := os.Getenv("DB_URL")
	if connStr == "" {
		log.Fatal("DB_URL environment variable is not set")
	}

	cli, err := ent.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}

	log.Println("Database connection established successfully")
	return cli
}

func startHTTPServer(cli *ent.Client) {
	r := gin.Default()
	userServ := setupExternalServices()
	registerHTTPRoutes(r, cli, userServ)

	log.Println("Starting HTTP server on port 8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to start HTTP server: %v", err)
	}
}

func setupExternalServices() grpc_clients.UserServiceClient {
	userServiceAddr := os.Getenv("USER_SERVICE")
	if userServiceAddr == "" {
		log.Println("Warning: USER_SERVICE environment variable is not set")
		return nil
	}

	log.Println("Connecting to user service at:", userServiceAddr)
	return grpc_clients.NewUserClient(userServiceAddr)
}

func registerHTTPRoutes(r *gin.Engine, cli *ent.Client, userServ grpc_clients.UserServiceClient) {
	handlersList := []struct {
		name     string
		register func(*gin.Engine)
	}{
		{"Organization", handlers.NewOrgHandler(cli, nil).RegisterRoutes},
		{"Department", handlers.NewDeptHandler(cli, nil).RegisterRoutes},
		{"Position", handlers.NewPositionHandler(cli, nil).RegisterRoutes},
		{"Employee", handlers.NewEmployeeHandler(cli, userServ).RegisterRoutes},
		{"Project", handlers.NewProjectHandler(cli, userServ).RegisterRoutes},
		{"Task", handlers.NewTaskHandler(cli).RegisterRoutes},
		{"TaskReport", handlers.NewTaskReportHandler(cli).RegisterRoutes},
		{"Label", handlers.NewLabelHandler(cli).RegisterRoutes},
		{"LeaveRequest", handlers.NewLeaveRequestHandler(cli).RegisterRoutes},
	}

	for _, h := range handlersList {
		log.Printf("Registering %s routes...", h.name)
		h.register(r)
	}
}

func startGRPCServer(cli *ent.Client) {
	srv := grpc.NewServer()
	registerGRPCServices(srv, cli)

	log.Println("Starting gRPC server on port 5000...")
	lis, err := net.Listen("tcp", ":5000")
	if err != nil {
		log.Fatalf("failed to listen on port 5000: %v", err)
	}

	if err := srv.Serve(lis); err != nil {
		log.Fatalf("gRPC server ended with error: %v", err)
	}
}

func registerGRPCServices(srv *grpc.Server, cli *ent.Client) {
	log.Println("Registering gRPC services...")

	entpb.RegisterOrganizationServiceServer(srv, entpb.NewOrganizationService(cli))
	entpb.RegisterDepartmentServiceServer(srv, entpb.NewDepartmentService(cli))
	entpb.RegisterPositionServiceServer(srv, entpb.NewPositionService(cli))
	entpb.RegisterEmployeeServiceServer(srv, entpb.NewEmployeeService(cli))
	entpb.RegisterProjectServiceServer(srv, entpb.NewProjectService(cli))
	entpb.RegisterTaskServiceServer(srv, entpb.NewTaskService(cli))
	entpb.RegisterTaskReportServiceServer(srv, entpb.NewTaskReportService(cli))
	entpb.RegisterExtServiceServer(srv, entpb.NewExtService(cli))

	log.Println("All gRPC services registered successfully")
}
