package usecase

import (
	"context"
	"repo-stat/api/internal/domain"
	repoDomain "repo-stat/pkg/domain"
)

type Pinger interface {
	Ping(ctx context.Context) domain.PingStatus
	Name() string
}

type RepositoryFetcher interface {
	Fetch(ctx context.Context, ownerName string, repoName string) (*repoDomain.Repository, error)
}
