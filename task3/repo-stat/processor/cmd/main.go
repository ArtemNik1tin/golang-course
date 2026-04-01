package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc"

	"repo-stat/api/config"
	"repo-stat/platform/logger"
	"repo-stat/processor/internal/adapter/collector"
	grpccontroller "repo-stat/processor/internal/controller/grpc"
	"repo-stat/processor/internal/usecase"
	processorpb "repo-stat/proto/processor"
)

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

	// Server setup
	listener, err := net.Listen("tcp", cfg.Services.ProcessorAddress)
	if err != nil {
		log.Error("failed to listen", "error", err)
		return err
	}

	collClient, err := collector.NewClient(cfg.Services.Collector, log)
	if err != nil {
		log.Error("failed to init collector client", "err", err)
		return err
	}
	defer func() { _ = collClient.Close() }()

	// Setup handlers
	repoUC := usecase.NewGetRepositoryUseCase(collClient)
	handler := grpccontroller.NewHandler(repoUC)

	grpcServer := grpc.NewServer()
	processorpb.RegisterProcessorServer(grpcServer, handler)

	go func() {
		log.Info("Processor gRPC server is running on " + cfg.Services.Processor)
		_ = grpcServer.Serve(listener)
	}()
	<-ctx.Done()
	grpcServer.GracefulStop()
	return nil
}

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
