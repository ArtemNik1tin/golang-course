package grpc

import (
	"context"
	"repo-stat/processor/internal/usecase"
	processorpb "repo-stat/proto/processor"
)

type Handler struct {
	processorpb.UnimplementedProcessorServer
	useCase *usecase.GetRepositoryUseCase
}

func NewHandler(uc *usecase.GetRepositoryUseCase) *Handler {
	return &Handler{
		useCase: uc,
	}
}

func (h *Handler) GetRepositoryInfo(
	ctx context.Context,
	req *processorpb.GetRepositoryInfoRequest,
) (*processorpb.GetRepositoryInfoResponse, error) {

	owner := req.GetOwner()
	repoName := req.GetRepo()

	repo, err := h.useCase.Execute(ctx, owner, repoName)
	if err != nil {
		return nil, err
	}

	return &processorpb.GetRepositoryInfoResponse{
		FullName:    repo.Name,
		Description: repo.Description,
		Stars:       repo.Stars,
		Forks:       repo.Forks,
		CreatedAt:   repo.CreatedAt,
	}, nil
}
