package processor

import (
	"context"
	"log/slog"
	apidomain "repo-stat/api/internal/domain"
	pkgdomain "repo-stat/pkg/domain"

	processorpb "repo-stat/proto/processor"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	log  *slog.Logger
	conn *grpc.ClientConn
	pb   processorpb.ProcessorClient
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
		pb:   processorpb.NewProcessorClient(conn),
	}, nil
}

func (c *Client) Fetch(ctx context.Context, owner string, repository string) (*pkgdomain.Repository, error) {
	getRepositoryInfoRequest := &processorpb.GetRepositoryInfoRequest{
		Owner: owner,
		Repo:  repository,
	}

	response, err := c.pb.GetRepositoryInfo(ctx, getRepositoryInfoRequest)
	if err != nil {
		c.log.Error("repository GetRepositoryInfo failed", "error", err)
		return nil, err
	}

	return &pkgdomain.Repository{
		Name:        response.FullName,
		Description: response.Description,
		Stars:       response.Stars,
		Forks:       response.Forks,
		CreatedAt:   response.CreatedAt,
	}, nil
}

func (c *Client) Ping(ctx context.Context) apidomain.PingStatus {
	_, err := c.pb.Ping(ctx, &processorpb.PingRequest{})
	if err != nil {
		c.log.Error("processor ping failed", "error", err)
		return apidomain.PingStatusDown
	}

	return apidomain.PingStatusUp
}

func (c *Client) Name() string {
	return "processor"
}

func (c *Client) Close() error {
	return c.conn.Close()
}
