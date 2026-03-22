package grpc

import (
	"context"
	"fmt"

	pb "github.com/ArtemNik1tin/distributed-github/api/proto"
	"github.com/ArtemNik1tin/distributed-github/collector/internal/usecase"
)

type Handler struct {
	pb.UnimplementedCollectorServiceServer
	useCase *usecase.GetRepositoryUseCase
}

func NewHandler(useCase *usecase.GetRepositoryUseCase) Handler {
	return Handler{useCase: useCase}
}

func (handler Handler) GetRepository(ctx context.Context, request *pb.RepositoryRequest) (*pb.RepositoryResponse, error) {
	repo, err := handler.useCase.Execute(ctx, request.OwnerName, request.RepoName)
	if err != nil {
		return nil, fmt.Errorf("handler error: %w", err)
	}
	return &pb.RepositoryResponse{
		Name:        repo.Name,
		Description: repo.Description,
		Stars:       int32(repo.Stars),
		Forks:       int32(repo.Forks),
		CreatedAt:   repo.CreatedAt,
	}, nil
}
