package grpc

import (
	"context"
	processorpb "repo-stat/proto/processor"
)

func (h *Handler) Ping(ctx context.Context, _ *processorpb.PingRequest) (*processorpb.PingResponse, error) {
	return &processorpb.PingResponse{
		Reply: "pong",
	}, nil
}
