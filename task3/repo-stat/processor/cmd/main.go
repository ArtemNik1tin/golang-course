package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

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

	collClient, err := collector.NewClient("localhost:50002", nil)
	if err != nil {
		log.Fatal(err)
	}
	defer collClient.Close()

	repoUC := usecase.NewGetRepositoryUseCase(collClient)

	handler := grpccontroller.NewHandler(repoUC)

	s := grpc.NewServer()

	processorpb.RegisterProcessorServer(s, handler)

	log.Println("Processor gRPC server is running on :50001...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
