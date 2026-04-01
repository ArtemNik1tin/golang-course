package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"repo-stat/pkg/domain"
)

type GithubRepository struct {
	FullName    string `json:"full_name"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Stars       int    `json:"stargazers_count"`
	Forks       int    `json:"forks_count"`
	CreatedAt   string `json:"created_at"`
}

type GitHubClient struct {
}

func (client GitHubClient) Fetch(ctx context.Context, ownerName string, repoName string) (*domain.Repository, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", ownerName, repoName)
	request, httpRequestErr := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if httpRequestErr != nil {
		return nil, fmt.Errorf("http request error: %w", httpRequestErr)
	}

	request.Header.Set("User-Agent", "repo-stat-collector")

	response, httpResponseErr := http.DefaultClient.Do(request)
	if httpResponseErr != nil {
		return nil, fmt.Errorf("http response error: %w", httpResponseErr)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github api error: %s", response.Status)
	}

	var repo GithubRepository
	if err := json.NewDecoder(response.Body).Decode(&repo); err != nil {
		return nil, err
	}

	return repo.toDomain(), nil
}

func (githubRepo *GithubRepository) toDomain() *domain.Repository {
	return &domain.Repository{
		Name:        githubRepo.FullName,
		Description: githubRepo.Description,
		Stars:       uint32(githubRepo.Stars),
		Forks:       uint32(githubRepo.Forks),
		CreatedAt:   githubRepo.CreatedAt,
	}
}
