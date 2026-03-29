package usecase

import (
	"context"
	"fmt"

	"github.com/ArtemNik1tin/distributed-github/collector/internal/domain"
)

type RepositoryFetcher interface {
	Fetch(ctx context.Context, ownerName string, repoName string) (*domain.Repository, error)
}

type GetRepositoryUseCase struct {
	fetcher RepositoryFetcher
}

func NewGetRepositoryUseCase(repositoryFetcher RepositoryFetcher) GetRepositoryUseCase {
	return GetRepositoryUseCase{fetcher: repositoryFetcher}
}

func (useCase *GetRepositoryUseCase) Execute(ctx context.Context, ownerName string, repoName string) (*domain.Repository, error) {
	if ownerName == "" || repoName == "" {
		return nil, fmt.Errorf("UseCase execute error: repoName and ownerName should be non empty.")
	}
	return useCase.fetcher.Fetch(ctx, ownerName, repoName)
}
