package usecase

import (
	"context"

	"github.com/ArtemNik1tin/distributed-github/gateway/internal/domain"
)

type RepositoryFetcher interface {
	Fetch(ctx context.Context, ownerName string, repoName string) (*domain.Repository, error)
}

type GetRepositoryUseCase struct {
	fetcher RepositoryFetcher
}

func NewGetRepositoryUseCase(fetcher RepositoryFetcher) GetRepositoryUseCase {
	return GetRepositoryUseCase{fetcher: fetcher}
}

func (useCase *GetRepositoryUseCase) Execute(ctx context.Context, ownerName string, repoName string) (*domain.Repository, error) {
	return useCase.fetcher.Fetch(ctx, ownerName, repoName)
}
