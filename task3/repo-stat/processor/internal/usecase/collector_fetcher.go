package usecase

import (
	"context"
	"fmt"
	"repo-stat/processor/internal/domain"
)

type CollectorFetcher interface {
	Fetch(ctx context.Context, owner, repo string) (*domain.Repository, error)
}

type GetRepositoryUseCase struct {
	fetcher CollectorFetcher
}

func NewGetRepositoryUseCase(fetcher CollectorFetcher) *GetRepositoryUseCase {
	return &GetRepositoryUseCase{fetcher: fetcher}
}

func (uc GetRepositoryUseCase) Execute(ctx context.Context, owner, repo string) (*domain.Repository, error) {
	if owner == "" || repo == "" {
		return nil, fmt.Errorf("invalid arguments: owner and repo must not be empty")
	}
	return uc.fetcher.Fetch(ctx, owner, repo)
}
