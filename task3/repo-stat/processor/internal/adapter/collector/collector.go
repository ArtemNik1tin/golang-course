package collector

import (
	"context"
	"log/slog"
	"repo-stat/pkg/domain"

	collectorpb "repo-stat/proto/collector"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	log  *slog.Logger
	conn *grpc.ClientConn
	pb   collectorpb.CollectorClient
}

func NewClient(address string, log *slog.Logger) (*Client, error) {
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		log:  log,
		conn: conn,
		pb:   collectorpb.NewCollectorClient(conn),
	}, nil
}

func (c *Client) Fetch(ctx context.Context, owner string, repository string) (*domain.Repository, error) {
	getRepositoryInfoRequest := &collectorpb.GetRepositoryInfoRequest{
		Owner: owner,
		Repo:  repository,
	}

	response, err := c.pb.GetRepositoryInfo(ctx, getRepositoryInfoRequest)
	if err != nil {
		c.log.Error("repository GetRepositoryInfo failed", "error", err)
		return nil, err
	}

	return &domain.Repository{
		Name:        response.FullName,
		Description: response.Description,
		Stars:       response.Stars,
		Forks:       response.Forks,
		CreatedAt:   response.CreatedAt,
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}
