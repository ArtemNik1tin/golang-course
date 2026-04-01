package controller

import (
	"context"
	"repo-stat/collector/internal/usecase"
	collectorpb "repo-stat/proto/collector"
)

type Handler struct {
	collectorpb.UnimplementedCollectorServer
	useCase *usecase.GitHubFetchUseCase
}

func NewHandler(uc *usecase.GitHubFetchUseCase) *Handler {
	return &Handler{useCase: uc}
}

func (h *Handler) GetRepositoryInfo(
	ctx context.Context,
	request *collectorpb.GetRepositoryInfoRequest,
) (*collectorpb.GetRepositoryInfoResponse, error) {
	repoInfo, err := h.useCase.Execute(ctx, request.Owner, request.Repo)
	if err != nil {
		return nil, err
	}

	return &collectorpb.GetRepositoryInfoResponse{
		FullName:    repoInfo.Name,
		Description: repoInfo.Description,
		Stars:       repoInfo.Stars,
		Forks:       repoInfo.Forks,
		CreatedAt:   repoInfo.CreatedAt,
	}, nil
}
