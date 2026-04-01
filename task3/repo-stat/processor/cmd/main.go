package main

import (
	"flag"
	"log"
	"net"

	"google.golang.org/grpc"

	"repo-stat/api/config"
	"repo-stat/platform/logger"
	"repo-stat/processor/internal/adapter/collector"
	grpccontroller "repo-stat/processor/internal/controller/grpc"
	"repo-stat/processor/internal/usecase"
	processorpb "repo-stat/proto/processor"
)

func main() {
	lis, err := net.Listen("tcp", ":50001")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// config
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "server configuration file")
	flag.Parse()

	cfg := config.MustLoad(configPath)

	// logger
	log := logger.MustMakeLogger(cfg.Logger.LogLevel)

	log.Info("starting server...")
	log.Debug("debug messages are enabled")

	collClient, err := collector.NewClient("localhost:50002", log)
	if err != nil {
		log.Error("failed to init", "err", err)
	}
	defer collClient.Close()

	repoUC := usecase.NewGetRepositoryUseCase(collClient)

	handler := grpccontroller.NewHandler(repoUC)

	s := grpc.NewServer()

	processorpb.RegisterProcessorServer(s, handler)

	log.Info("Processor gRPC server is running on :50001...")
	if err := s.Serve(lis); err != nil {
		log.Error("failed to serve", "err", err)
	}
}
