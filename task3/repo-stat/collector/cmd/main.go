package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"

	"repo-stat/collector/config"
	"repo-stat/collector/internal/adapter/github"
	grpchandler "repo-stat/collector/internal/controller/grpc"
	"repo-stat/collector/internal/usecase"
	"repo-stat/platform/logger"
	collectorpb "repo-stat/proto/collector"

	"google.golang.org/grpc"
)

func main() {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	if err := run(ctx); err != nil {
		_, err = fmt.Fprintln(os.Stderr, err)
		if err != nil {
			fmt.Printf("launching server error: %s\n", err)
		}
		cancel()
		os.Exit(1)
	}
	cancel()
}

func run(ctx context.Context) error {
	// config
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "server configuration file")
	flag.Parse()

	cfg := config.MustLoad(configPath)

	// logger
	log := logger.MustMakeLogger(cfg.Logger.LogLevel)

	log.Info("starting server...")
	log.Debug("debug messages are enabled")

	// Hendlers setup
	client := github.GitHubClient{}
	getRepoUseCase := usecase.NewGitHubFetchUseCase(client)
	grpcHandler := grpchandler.NewHandler(getRepoUseCase)
	listener, err := net.Listen("tcp", cfg.GRPC.Address)
	if err != nil {
		return fmt.Errorf("listen error: %w", err)
	}

	// Server setup
	grpcServer := grpc.NewServer()
	collectorpb.RegisterCollectorServer(grpcServer, grpcHandler)

	go func() {
		log.Info("Collector gRPC server is running on " + cfg.GRPC.Address)
		_ = grpcServer.Serve(listener)
	}()
	<-ctx.Done()
	grpcServer.GracefulStop()
	return nil
}
