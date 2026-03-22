package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ArtemNik1tin/distributed-github/collector/internal/domain"
)

type GitHubClient struct {
}

func (client GitHubClient) Fetch(ctx context.Context, ownerName string, repoName string) (*domain.Repository, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", ownerName, repoName)
	request, httpRequestErr := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if httpRequestErr != nil {
		return nil, fmt.Errorf("http request error: %w", httpRequestErr)
	}

	response, httpResponseErr := http.DefaultClient.Do(request)
	if httpResponseErr != nil {
		return nil, fmt.Errorf("http response error: %w", httpResponseErr)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github api error: %s", response.Status)
	}

	var repo domain.Repository
	if err := json.NewDecoder(response.Body).Decode(&repo); err != nil {
		return nil, err
	}

	return &repo, nil
}
