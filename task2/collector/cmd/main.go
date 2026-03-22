package cmd

import (
	"fmt"
	"net"
	"os"

	pb "github.com/ArtemNik1tin/distributed-github/api/proto"
	"github.com/ArtemNik1tin/distributed-github/collector/internal/adapter/github"
	grpchandler "github.com/ArtemNik1tin/distributed-github/collector/internal/adapter/grpc"
	"github.com/ArtemNik1tin/distributed-github/collector/internal/usecase"
	"google.golang.org/grpc"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	client := github.GitHubClient{}
	getRepoUseCase := usecase.NewGetRepositoryUseCase(client)
	grpcHandler := grpchandler.NewHandler(&getRepoUseCase)
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		return fmt.Errorf("Listen error: %w", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterCollectorServiceServer(grpcServer, grpcHandler)
	return grpcServer.Serve(listener)
}
