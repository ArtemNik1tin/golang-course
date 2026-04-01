package main

import (
	"fmt"
	"net"
	"os"

	"repo-stat/collector/internal/adapter/github"
	grpchandler "repo-stat/collector/internal/controller/grpc"
	"repo-stat/collector/internal/usecase"
	collectorpb "repo-stat/proto/collector"

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
	getRepoUseCase := usecase.NewGitHubFetchUseCase(client)
	grpcHandler := grpchandler.NewHandler(getRepoUseCase)
	listener, err := net.Listen("tcp", ":50002")
	if err != nil {
		return fmt.Errorf("Listen error: %w", err)
	}
	grpcServer := grpc.NewServer()
	collectorpb.RegisterCollectorServer(grpcServer, grpcHandler)
	return grpcServer.Serve(listener)
}
