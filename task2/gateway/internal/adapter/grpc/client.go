package grpc

import (
	"context"
	"fmt"

	pb "github.com/ArtemNik1tin/distributed-github/api/proto"
	"github.com/ArtemNik1tin/distributed-github/gateway/internal/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type CollectorClient struct {
	grpcClient pb.CollectorServiceClient
	conn       *grpc.ClientConn
}

func NewCollectorClient(address string) (*CollectorClient, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}
	return &CollectorClient{
		grpcClient: pb.NewCollectorServiceClient(conn),
	}, nil
}

func (client CollectorClient) Fetch(ctx context.Context, ownerName string, repoName string) (*domain.Repository, error) {
	grpcRequest := &pb.RepositoryRequest{OwnerName: ownerName, RepoName: repoName}
	response, err := client.grpcClient.GetRepository(ctx, grpcRequest)
	if err != nil {
		return nil, fmt.Errorf("GetRepository error: %w", err)
	}
	return &domain.Repository{
		Name:        response.Name,
		Description: response.Description,
		Stars:       int(response.Stars),
		Forks:       int(response.Forks),
		CreatedAt:   response.CreatedAt,
	}, nil
}

func (client *CollectorClient) Close() error {
	return client.conn.Close()
}
