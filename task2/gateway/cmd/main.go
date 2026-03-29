package main

import (
	"fmt"
	"net/http"
	"os"

	_ "github.com/ArtemNik1tin/distributed-github/gateway/docs"
	grpcclient "github.com/ArtemNik1tin/distributed-github/gateway/internal/adapter/grpc"
	"github.com/ArtemNik1tin/distributed-github/gateway/internal/adapter/rest"
	"github.com/ArtemNik1tin/distributed-github/gateway/internal/usecase"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title GitHub Repository API Gateway
// @version 1.0
// @description API Gateway for GitHub repository information
// @host localhost:8080
// @BasePath /
func main() {
	collectorAddr := os.Getenv("COLLECTOR_ADDR")
	if collectorAddr == "" {
		collectorAddr = "localhost:50051"
	}
	collectorClient, err := grpcclient.NewCollectorClient(collectorAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create collector client: %v\n", err)
		os.Exit(1)
	}
	defer collectorClient.Close()

	getRepoUseCase := usecase.NewGetRepositoryUseCase(collectorClient)

	restHandler := rest.NewHandler(&getRepoUseCase)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /repos/{owner}/{repo}", restHandler.GetRepository)
	mux.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)

	fmt.Println("Gateway starting on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}
