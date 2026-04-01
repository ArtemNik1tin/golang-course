package usecase

import (
	"context"
	"repo-stat/collector/internal/domain"
)

type GitHubFetcher interface {
	Fetch(ctx context.Context, ownerName, repoName string) (*domain.Repository, error)
}

type GitHubFetchUseCase struct {
	fetcher GitHubFetcher
}

func NewGitHubFetchUseCase(fetcher GitHubFetcher) *GitHubFetchUseCase {
	return &GitHubFetchUseCase{fetcher: fetcher}
}

func (uc *GitHubFetchUseCase) Execute(ctx context.Context, ownerName, repoName string) (*domain.Repository, error) {
	return uc.fetcher.Fetch(ctx, ownerName, repoName)
}
